package main

import (
	"github.com/pkg/errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

var kvpRegex = regexp.MustCompile(`\s*(?P<key>[a-zA-Z0-9.#\-_]+)=(?P<value>"(\w|\s)*"|[!-~]*)\s*`)

// Regex to find/replace with nothing
var normalizeRegex = regexp.MustCompile(`^sample#`)

var numWithUnitRegex = regexp.MustCompile(`^(\d+(\.\d+)?)([a-zA-Z]+)$`)

// parseKvp This function returns a map of all the key-value pairs within a log
// message. This returns an `error` object if the string isn't a valid kvp
// format.
func parseKvp(msg string) (map[string]interface{}, error) {
	matches, err := extractKvp(msg)

	if err != nil {
		return nil, err
	}

	payload := make(map[string]interface{})

	for _, m := range matches {
		k := m[0]
		v := m[1]

		if v == "" {
			payload[k] = nil
			continue
		} else if v[0] == '"' && v[len(v)-1] == '"' {
			v = v[1 : len(v)-1]
		}

		samples := parseMetricValue(v)

		for _, s := range samples {
			name := normalizeMetricName(k, s.Unit)
			payload[name] = s.Value
		}
	}

	return payload, nil
}

// extractKvp Returns a slice of tuple strings where t[0] is the key and
// t[1] is the value. The slice will be nil if there's a parsing error
func extractKvp(msg string) ([][]string, error) {
	tuples := [][]string{}
	logger.Debugf("matching `%s`", msg)

	for i := 0; i < len(msg); {
		match := kvpRegex.FindStringSubmatchIndex(msg[i:])

		if match == nil {
			return nil, errors.New("Parsing error: invalid kvp format")
		}

		if match[0] != 0 {
			return nil, errors.New("Parsing error: invalid kvp format")
		}

		logger.Debugf("Match? %v, `%s`", match, msg[i+match[0]:i+match[1]])

		tuples = append(tuples, []string{msg[i+match[2] : i+match[3]], msg[i+match[4] : i+match[5]]})

		logger.Debugf("%+v", tuples)

		i += match[1]
	}

	return tuples, nil
}

func normalizeMetricName(rawMetric string, unit *string) string {
	nn := normalizeRegex.ReplaceAllString(rawMetric, "")
	// NOTE: For whatever reason, some of the keys for the Postgres's metrics use
	// hyphens instead of underscores; this simply converts it...
	nn = strings.Replace(nn, "-", "_", -1)

	if nn == "memory_pgpgin" || nn == "memory_pgpgout" || unit == nil {
		// Don't suffix, since the only unit is pages
		return nn
	}

	return nn + "_" + *unit
}

type metricSample struct {
	Value interface{}
	Unit  *string
}

// parseMetricValue Returns a parsed value. If `unit` is nil, that means that
// the parsed value is simply a string.
func parseMetricValue(value string) []*metricSample {
	samples := make([]*metricSample, 0)

	// Checks for unit-less numerical value
	if n, err := parseNumber(value); err == nil {
		f := big.NewFloat(n)

		if f.IsInt() {
			i, _ := f.Int64()
			samples = append(samples, &metricSample{Value: int(i)})
		} else {
			samples = append(samples, &metricSample{Value: n})
		}

		return samples
	}

	matches := numWithUnitRegex.FindAllStringSubmatch(value, -1)

	// Checks if value doesn't match a numerical value w/ unit
	if len(matches) == 0 {
		samples = append(samples, &metricSample{Value: value})

		return samples
	}

	pv, _ := parseNumber(matches[0][1])
	pUnit := matches[0][3]

	samples = append(samples, &metricSample{Value: pv, Unit: &pUnit})

	nv, nUnit := normalizeStorageSize(pv, pUnit)

	if nUnit != pUnit {
		samples = append(samples, &metricSample{Value: nv, Unit: &nUnit})
	}

	return samples
}

func normalizeStorageSize(num float64, unit string) (float64, string) {
	switch strings.ToLower(unit) {
	case "b", "bytes":
		return num / 1024 / 1024, "MB"
	case "kb", "kilobytes":
		return num / 1024, "MB"
	case "mb", "megabytes":
		return num, "MB"
	case "gb", "gigabytes":
		return num * 1024, "MB"
	default:
		// Unsure what the unit of measure is, just return it as is
		return num, unit
	}
}

func parseNumber(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
