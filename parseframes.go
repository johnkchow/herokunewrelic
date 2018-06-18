package main

import (
	"bytes"
	bp "github.com/oxtoacart/bpool"
	"github.com/pkg/errors"
	"io"
	"os"
	"strconv"
)

func minInt(x, y int) int {
	if x < y {
		return x
	}

	return y
}

// By default, we will pre-allocate ~500kB of byte buffers
var bPool = bp.NewBytePool(
	atoi(getEnvFallback("REQUEST_BUFFER_POOL_SIZE", "500")),
	atoi(getEnvFallback("REQUEST_BUFFER_SIZE", "1024")),
)

func getEnvFallback(key string, fallback string) string {
	res := os.Getenv(key)
	if res == "" {
		return fallback
	}
	return res
}

func atoi(s string) int {
	val, err := strconv.Atoi(s)

	if err != nil {
		panic(err)
	}

	return val
}

// parseFrames Parses a request's body and returns an array framed messsages
//
// Heroku's request body are formatted using the octect counting framing method.
//
// See more:
// 	https://devcenter.heroku.com/articles/log-drains
// 	https://tools.ietf.org/html/rfc6587#section-3.4.1
func parseFrames(body io.Reader) ([][]byte, error) {
	buffer := bPool.Get()
	defer bPool.Put(buffer)

	return parseFramesWithBuffer(body, buffer)
}

func parseFramesWithBuffer(body io.Reader, buffer []byte) ([][]byte, error) {
	// State machine to parse body.
	// 1 - "length" - Reads until finds space
	// 2 - "body" - Reads until length is reached

	state := 1
	lastState := 0
	msgLen := 0
	bufIdx := -1
	bufLen := 0
	totalRead := 0
	eof := false

	var err error
	var lenBuffer bytes.Buffer
	var msgBuffer bytes.Buffer
	messages := [][]byte{}

	for !eof || bufIdx < bufLen {
		logger.Debugf("totalRead: %v, msgLen: %v, bufIdx: %v, bufLen: %v, eof: %v", totalRead, msgLen, bufIdx, bufLen, eof)

		if bufLen == 0 || bufIdx == bufLen {

			bufLen, err = body.Read(buffer)
			totalRead += bufLen
			bufIdx = 0

			logger.Debugf("Read %v bytes: `%s`", bufLen, buffer)

			// NOTE: This can happen if in the last read, we had read the end of the
			// body and it happened to match up with the buffer's capacity
			if bufLen == 0 && err == io.EOF {
				break
			}

			if err != nil && err != io.EOF {
				return nil, errors.Wrap(err, "Reading body failed!")
			}

			if err == io.EOF {
				eof = true
			}
		}

		if state == 1 {
			logger.Debugf("State 1")

			lastState = 1
			mi := bufIdx

			for mi < bufLen && buffer[mi] != ' ' {
				mi++
			}

			logger.Debugf("Writing to lenBuffer: '%s'", buffer[bufIdx:mi])

			lenBuffer.Write(buffer[bufIdx:mi])
			bufIdx = mi

			if mi < bufLen {
				var err error

				msgLen, err = strconv.Atoi(lenBuffer.String())

				if err != nil {
					return nil, errors.Wrap(err, "Error converting string to int")
				}

				logger.Debugw("Transitioning state 2",
					"msgLen", msgLen,
					"bufIdx", bufIdx,
				)

				bufIdx++
				state = 2
				lenBuffer.Reset()
			}
		}

		if state == 2 {
			logger.Debugf("State 2")
			lastState = 2

			bound := minInt(bufIdx+msgLen, bufLen)

			logger.Debugf("Writing '%v' to msgBuffer", string(buffer[bufIdx:bound]))
			logger.Debugf("msgLen: %v, bound: %v, bufIdx: %v, readLen: %v, bufLen: %v", msgLen, bound, bufIdx, len(buffer[bufIdx:bound]), bufLen)
			msgBuffer.Write(buffer[bufIdx:bound])

			msgLen = msgLen - len(buffer[bufIdx:bound])
			bufIdx = bound

			logger.Debugw("Finished state 2",
				"msgLen", msgLen,
				"bufIdx", bufIdx,
			)

			if msgLen == 0 {
				state = 1
				// NOTE: We "copy" from the msgBuffer so that the slice's underlying
				// array is a complete copy. Otherwise, it's possible that subsequent
				// msgBuffer.Write will modify the array that's backing multiple
				// slices.
				messages = append(messages, append([]byte{}, msgBuffer.Bytes()...))
				msgBuffer.Reset()
			}
		}
	}

	for _, m := range messages {
		logger.Debugf("Message '%s'", string(m))
	}
	logger.Debugf("Final state: %v, lastState: %v", state, lastState)

	if state == 1 && lastState == 2 {
		return messages, nil
	}

	return messages, errors.New("Parsing error")
}
