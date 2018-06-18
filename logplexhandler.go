package main

import (
	nr "github.com/newrelic/go-agent"
	"net/http"
	"strconv"
)

func newLogplexHandler(app nr.Application) http.HandlerFunc {
	return func(rr http.ResponseWriter, rq *http.Request) {
		rawMsgs, err := parseFrames(rq.Body)
		appName, _, _ := rq.BasicAuth()
		logDrainToken := rq.Header.Get("Logplex-Drain-Token")

		if appName == "" {
			// Unfortunately, the Logplex message does **NOT** contain the Heroku
			// app's name whatsoever, so we must set it as the username when adding a
			// HTTPs log drain
			logger.Errorf("Basic auth's username, used for Heroku app name, was not set for drain token `%s` thus we cannot record metric", logDrainToken)
			rr.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err != nil {
			logger.Errorf("Error parsing frames: %s", err.Error())
			rr.WriteHeader(http.StatusBadRequest)
			return
		}

		msgCount, _ := strconv.Atoi(rq.Header.Get("Logplex-Msg-Count"))

		if msgCount != len(rawMsgs) {
			// TODO record error to Bugsnag
			logger.Error("Frame count does not match Logplex-Msg-Count header")
			rr.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, rawMsg := range rawMsgs {
			logger.Debugf("Parsing raw msg: `%s`", rawMsg)

			msg, err := parseLogplex(rawMsg)

			logger.Debugf("Logplex: %+v", msg)

			if err != nil {
				// TODO: Record error to Bugsnag
				logger.Warnf("Malformed Logplex format")
				rr.WriteHeader(http.StatusBadRequest)
				return
			}

			var metrics map[string]interface{}

			metrics, err = parseMetrics(msg.Msg)

			if err != nil {
				logger.Debugf("Error parsing metrics: %s", err.Error())
				continue
			}

			// NewRelic only receives either seconds or milliseconds.
			// See
			// https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/insert-custom-events-insights-api#timestamps
			metrics["timestamp"] = msg.Timestamp.UnixNano() / 1000000

			// We cannot use `appName` since it's reserved to the current app's name
			metrics["sourceAppName"] = appName

			app.RecordCustomEvent("DynoMetric", metrics)
		}

		rr.WriteHeader(http.StatusOK)
	}
}
