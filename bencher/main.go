package main

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var logger *zap.Logger

func main() {
	logger, _ = zap.NewProduction()
	rand.Seed(int64(time.Now().Nanosecond()))

	client := &http.Client{}
	url := fmt.Sprintf("https://bencher:%s@reflektive-herokunewrelic.herokuapp.com", os.Getenv("APP_SECRET"))

	concurrency := 75

	wg := &sync.WaitGroup{}
	wg.Add(1)

	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				body, frameCount := bodyWithFrames()

				logger.Debug(fmt.Sprintf("URL: %s", url))
				logger.Debug(fmt.Sprintf("Body Length/Frames: %d/%d", len(body), frameCount))

				req, _ := http.NewRequest("POST", url, strings.NewReader(body))
				req.Header.Add("Logplex-Msg-Count", strconv.Itoa(frameCount))
				req.Header.Add("Logplex-Drain-Token", "bencher-drain-token")

				benchmark("req", func() {
					resp, err := client.Do(req)
					if err != nil {
						logger.Debug(fmt.Sprintf("Error: %v", err))
					} else {
						logger.Debug(fmt.Sprintf("Status code: %d", resp.StatusCode))
					}
				})
			}
		}()
	}

	wg.Wait()
}

func bodyWithFrames() (string, int) {
	var frames = [][]byte{
		[]byte("195 <40>1 2012-11-30T06:45:29+00:00 host heroku web.4 - source=web.4 dyno=heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b sample#load_avg_1m=0.27 sample#load_avg_5m=0.22 sample#load_avg_15m=0.19"),
		[]byte("332 <40>2 2012-11-30T06:45:29+00:00 host heroku web.4 - source=web.4 dyno=heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b sample#memory_total=7020.55MB sample#memory_rss=7013.48MB sample#memory_cache=7.08MB sample#memory_swap=0.00MB sample#memory_pgpgin=1949515pages sample#memory_pgpgout=779250pages sample#memory_quota=14336.00MB"),
		[]byte("551 <40>3 2012-11-30T06:45:29+00:00 host app heroku-postgres - source=HEROKU_POSTGRESQL_ONYX sample#current_transaction=153674642 sample#db_size=53666224664bytes sample#tables=105 sample#active-connections=20 sample#waiting-connections=0 sample#index-cache-hit-rate=0.99832 sample#table-cache-hit-rate=0.9657 sample#load-avg-1m=0.01 sample#load-avg-5m=0.045 sample#load-avg-15m=0.025 sample#read-iops=0.58264 sample#write-iops=3.8988 sample#memory-total=8173656kB sample#memory-free=1124196kB sample#memory-cached=5935596kB sample#memory-postgres=111996kB"),
	}

	randMsgCount := rand.Intn(100)
	logger.Debug(fmt.Sprintf("Rand %d", randMsgCount))
	randMsgCount = rand.Intn(100)
	logger.Debug(fmt.Sprintf("Rand %d", randMsgCount))

	for i := 0; i < randMsgCount+50; i++ {
		frames = append(frames, []byte(randomLogFrame()))
	}

	buffer := bytes.Buffer{}

	for _, p := range frames {
		buffer.Write(p)
	}

	return buffer.String(), len(frames)
}

func benchmark(name string, f func()) {
	start := time.Now()

	f()

	end := time.Now()
	delta := end.Sub(start) / time.Millisecond

	logger.Info(fmt.Sprintf("Finished %s in %dms", name, delta))
}

func randomLogFrame() string {
	time := time.Now().UTC()
	str := fmt.Sprintf("<40>3 %s host app someproc.1 - %s",
		time.Format("2006-01-02T15:04:05Z07:00"),
		`at=info method=GET path="/" host=some-app.herokuapp.com request_id=78b767f8-fa4d-42a9-8b1d-b71bc22a8934 fwd="127.0.0.1" dyno=web.4 connect=1ms service=28ms status=200 bytes=1084 protocol=https`,
	)

	return fmt.Sprintf("%d %s", len(str), str)
}
