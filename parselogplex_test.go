package main

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func testRegex(regexStr, body string) {
	regex, err := regexp.Compile(regexStr)

	if err != nil {
		panic(err)
	}

	res := regex.FindAllStringSubmatch(body, -1)
	logger.Infof("Found? %v", regex.FindStringIndex(body))

	logger.Infof("regex match? %v", len(res) != 0)
	logger.Infof("  regex: %s", regexStr)
	logger.Infof("  str:   %s", body)

	if len(res) != 0 {
		names := regex.SubexpNames()

		for i, str := range res[0] {
			logger.Infof("  `%s`: `%s`", names[i], str)
		}
	}
}

func TestSuccessfulLogplexParsing(t *testing.T) {
	body := []byte(`<40>1 2012-11-30T06:45:29+00:00 host app heroku-postgres - source=web.4 dyno=heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b sample#load_avg_1m=0.27 sample#load_avg_5m=0.22 sample#load_avg_15m=0.19`)
	msg, err := parseLogplex(body)

	assert.Nil(t, err)
	assert.Equal(t, "source=web.4 dyno=heroku.18186867.69773a68-35f0-4cc3-aa48-e956a04c2c8b sample#load_avg_1m=0.27 sample#load_avg_5m=0.22 sample#load_avg_15m=0.19", msg.Msg)
}
