package main

import (
	"bytes"
	"github.com/johnkchow/herokunewrelic/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var frames = [][]byte{
	[]byte("195 <40>1 2012-11-30T06:45:29+00:00 host heroku web.4 - source=web.4 dyno=heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b sample#load_avg_1m=0.27 sample#load_avg_5m=0.22 sample#load_avg_15m=0.19"),
	[]byte("77 <40>3 2012-11-30T06:45:29+00:00 host app web.4 - Some application log message"),
	[]byte(`245 <40>3 2012-11-30T06:45:29+00:00 host heroku router - at=info method=GET path="/" host=some-app.herokuapp.com request_id=78b767f8-fa4d-42a9-8b1d-b71bc22a8934 fwd="127.0.0.1" dyno=web.4 connect=1ms service=28ms status=200 bytes=1084 protocol=https`),
	[]byte("332 <40>2 2012-11-30T06:45:29+00:00 host heroku web.4 - source=web.4 dyno=heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b sample#memory_total=7020.55MB sample#memory_rss=7013.48MB sample#memory_cache=7.08MB sample#memory_swap=0.00MB sample#memory_pgpgin=1949515pages sample#memory_pgpgout=779250pages sample#memory_quota=14336.00MB"),
	[]byte("551 <40>3 2012-11-30T06:45:29+00:00 host app heroku-postgres - source=HEROKU_POSTGRESQL_ONYX sample#current_transaction=153674642 sample#db_size=53666224664bytes sample#tables=105 sample#active-connections=20 sample#waiting-connections=0 sample#index-cache-hit-rate=0.99832 sample#table-cache-hit-rate=0.9657 sample#load-avg-1m=0.01 sample#load-avg-5m=0.045 sample#load-avg-15m=0.025 sample#read-iops=0.58264 sample#write-iops=3.8988 sample#memory-total=8173656kB sample#memory-free=1124196kB sample#memory-cached=5935596kB sample#memory-postgres=111996kB"),
	[]byte("387 <134>1 2018-06-15T21:39:20+00:00 host app heroku-redis - source=HEROKU_REDIS_WHITE sample#active-connections=1 sample#load-avg-1m=0.04 sample#load-avg-5m=0.09 sample#load-avg-15m=0.085 sample#read-iops=0 sample#write-iops=0 sample#memory-total=15664360kB sample#memory-free=9891092kB sample#memory-cached=3415688kB sample#memory-redis=1757408bytes sample#hit-rate=1 sample#evicted-keys=0"),
	[]byte("293 <40>3 2012-11-30T06:45:29+00:00 host app web.4 - event_name=PgBouncerMetric source=PGBOUNCER sample#avg_req=1 sample#avg_recv=2 sample#avg_sent=3 sample#avg_query=4 sample#cl_active=5 sample#cl_waiting=6 sample#sv_active=7 sample#sv_idle=8 sample#sv_used=9 sample#sv_login=10 sample#maxwait=11"),
	[]byte("88 <40>3 2012-11-30T06:45:29+00:00 host app web.4 - event_name=AppDefinedCustomMetric foo=1"),
	[]byte(`197 <40>3 2012-11-30T06:45:29+00:00 host heroku router - at=error code=H11 desc="Backlog too deep" method=GET path="/" host=myapp.herokuapp.com fwd=17.17.17.17 dyno= connect= service= status=503 bytes=`),
}

// NOTE: This variable is a concat of the `frames` above, which is done in init()
var logplexBody []byte

func TestLogplexHandlerNewRelicIsCalledCorrectly(t *testing.T) {
	req := buildRequest()

	rr := httptest.NewRecorder()

	app := new(mocks.Application)

	app.On("RecordCustomEvent", "DynoMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "PostgresMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "RedisMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "PgBouncerMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "AppDefinedCustomMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "HerokuError", mock.Anything).Return(nil)

	newLogplexHandler(app)(rr, req)
	app.AssertNumberOfCalls(t, "RecordCustomEvent", 7)

	assert.Equal(t, map[string]interface{}{
		"source":        "web.4",
		"timestamp":     int64(1354257929000),
		"dyno":          "heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b",
		"load_avg_1m":   0.27,
		"load_avg_5m":   0.22,
		"load_avg_15m":  0.19,
		"sourceAppName": "my-heroku-app",
	}, app.Calls[0].Arguments.Get(1))

	assert.Equal(t, "DynoMetric", app.Calls[0].Arguments.Get(0))
	assert.Equal(t, "DynoMetric", app.Calls[1].Arguments.Get(0))
	assert.Equal(t, "PostgresMetric", app.Calls[2].Arguments.Get(0))
	assert.Equal(t, "RedisMetric", app.Calls[3].Arguments.Get(0))
	assert.Equal(t, "PgBouncerMetric", app.Calls[4].Arguments.Get(0))
	assert.Equal(t, "AppDefinedCustomMetric", app.Calls[5].Arguments.Get(0))
	assert.Equal(t, "HerokuError", app.Calls[6].Arguments.Get(0))
}

func TestLogplexHandlerEmptySuccessResponseReturned(t *testing.T) {
	req := buildRequest()

	rr := httptest.NewRecorder()

	app := new(mocks.Application)

	app.On("RecordCustomEvent", "DynoMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "PostgresMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "RedisMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "PgBouncerMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "AppDefinedCustomMetric", mock.Anything).Return(nil)
	app.On("RecordCustomEvent", "HerokuError", mock.Anything).Return(nil)

	newLogplexHandler(app)(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 0, rr.Body.Len())
}

func init() {
	buffer := bytes.Buffer{}

	for _, p := range frames {
		buffer.Write(p)
	}

	logplexBody = buffer.Bytes()
}

func buildRequest() *http.Request {
	req, _ := http.NewRequest("GET", "/", bytes.NewReader(logplexBody))

	req.Header.Add("Logplex-Msg-Count", "9")
	req.Header.Add("Logplex-Msg-Id", "someid")
	req.Header.Add("Logplex-Drain-Token", "draintoken")
	req.Header.Add("User-Agent", "Logplex/v72")
	req.Header.Add("Content-Type", "application/logplex-1")

	req.SetBasicAuth("my-heroku-app", "doesnt-matter-for-this-test")

	return req
}
