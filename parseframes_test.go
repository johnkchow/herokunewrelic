package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFrameReturnsCorrectMessages(t *testing.T) {
	msgs, err := parseFrames(bytes.NewReader([]byte("3 foo2 ba")))

	assert.Nil(t, err)
	assert.Len(t, msgs, 2)
	assert.Equal(t, [][]byte{[]byte("foo"), []byte("ba")}, msgs)
}

func TestParseFrameBodyExceedingBufferSize(t *testing.T) {
	body := bytes.NewReader([]byte("10 123456789 5 12 3 "))
	buffer := make([]byte, 5)
	msgs, err := parseFramesWithBuffer(body, buffer)

	assert.Nil(t, err)
	assert.Len(t, msgs, 2)
	assert.Equal(t, [][]byte{[]byte("123456789 "), []byte("12 3 ")}, msgs)
}

func TestParseFrameLengthEndOfBuffer(t *testing.T) {
	body := bytes.NewReader([]byte("2 1210 1234567890"))
	buffer := make([]byte, 5)
	msgs, err := parseFramesWithBuffer(body, buffer)

	assert.Nil(t, err)
	assert.Len(t, msgs, 2)
	assert.Equal(t, [][]byte{[]byte("12"), []byte("1234567890")}, msgs)
}

func TestParseFrameLengthExceedActualStringSize(t *testing.T) {
	msg := []byte("9 1")

	body := bytes.NewReader(msg)
	_, err := parseFrames(body)

	assert.NotNil(t, err)
}

func TestParseFrameLogplex(t *testing.T) {
	msg := []byte("387 <134>1 2018-06-15T21:39:20+00:00 host app heroku-redis - source=HEROKU_REDIS_WHITE sample#active-connections=1 sample#load-avg-1m=0.04 sample#load-avg-5m=0.09 sample#load-avg-15m=0.085 sample#read-iops=0 sample#write-iops=0 sample#memory-total=15664360kB sample#memory-free=9891092kB sample#memory-cached=3415688kB sample#memory-redis=1757408bytes sample#hit-rate=1 sample#evicted-keys=0")

	body := bytes.NewReader(msg)
	frames, err := parseFrames(body)

	assert.Nil(t, err)

	assert.Equal(t, string(frames[0]), "<134>1 2018-06-15T21:39:20+00:00 host app heroku-redis - source=HEROKU_REDIS_WHITE sample#active-connections=1 sample#load-avg-1m=0.04 sample#load-avg-5m=0.09 sample#load-avg-15m=0.085 sample#read-iops=0 sample#write-iops=0 sample#memory-total=15664360kB sample#memory-free=9891092kB sample#memory-cached=3415688kB sample#memory-redis=1757408bytes sample#hit-rate=1 sample#evicted-keys=0")
}
