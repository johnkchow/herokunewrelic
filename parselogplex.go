package main

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"time"
)

// LogplexMsg Represents Heroku's HTTPS log drain version of the syslog format:
// > “application/logplex-1” does not conform to RFC5424. It leaves out
// > STRUCTURED-DATA but does not replace it with a NILVALUE.
// See https://devcenter.heroku.com/articles/log-drains#https-drain-caveats for
// more details
type LogplexMsg struct {
	Priority  uint8
	Version   uint16
	Timestamp time.Time
	Hostname  string
	Appname   string
	ProcID    string
	MsgID     string
	Msg       string
}

var regex = regexp.MustCompile(`^\s*<(?P<priority>\d+)>(?P<version>\d+) (?P<timestamp>[!-~]+) (?P<hostname>[!-~]+) (?P<appname>[!-~]+) (?P<procID>[!-~]+) (?P<msgID>[!-~]+) (?P<msg>(.|\s)*)$`)

func parseLogplex(body []byte) (*LogplexMsg, error) {
	res := regex.FindAllSubmatch(body, -1)

	if len(res) == 0 {
		return nil, errors.New("Body does not match expected Logplex format")
	}

	matches := res[0]

	logplexMsg := new(LogplexMsg)

	for i, n := range matches {
		switch i {
		case 1:
			val, err := strconv.Atoi(string(n))

			if err != nil {
				return nil, errors.Wrap(err, "Priority parsing error")
			}

			logplexMsg.Priority = uint8(val)
		case 2:
			val, err := strconv.Atoi(string(n))

			if err != nil {
				return nil, errors.Wrap(err, "Priority parsing error")
			}

			logplexMsg.Version = uint16(val)
		case 3:
			ts, err := time.Parse(time.RFC3339Nano, string(n))

			if err != nil {
				return nil, errors.Wrap(err, "Timestamp parsing error")
			}

			logplexMsg.Timestamp = ts
		case 4:
			logplexMsg.Hostname = string(n)
		case 5:
			logplexMsg.Appname = string(n)
		case 6:
			logplexMsg.ProcID = string(n)
		case 7:
			logplexMsg.MsgID = string(n)
		case 8:
			logplexMsg.Msg = string(n)
		}
	}

	return logplexMsg, nil
}
