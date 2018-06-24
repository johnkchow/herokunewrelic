package main

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

var kvpRegex = regexp.MustCompile(`^\s*([!-~]+=[!-~]+\s*)+$`)

// Regex to find/replace with nothing
var normalizeRegex = regexp.MustCompile(`^sample#`)

var numWithUnitRegex = regexp.MustCompile(`^(\d+(\.\d+)?)([a-zA-Z]+)$`)

// parseKvp This function returns a map of all the key-value pairs within a log
// message. This returns an `error` object if the string isn't a valid kvp
// format.
func parseKvp(msg string) (map[string]interface{}, error) {

	if !kvpRegex.MatchString(msg) {
		return nil, errors.New("Parsing error: invalid kvp format")
	}

	payload := make(map[string]interface{})

	for _, str := range strings.Split(msg, " ") {
		parts := strings.Split(str, "=")

		samples := parseMetricValue(parts[1])

		for _, s := range samples {
			name := normalizeMetricName(parts[0], s.Unit)
			payload[name] = s.Value
		}
	}

	return payload, nil
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
		samples = append(samples, &metricSample{Value: n})

		return samples
	}

	matches := numWithUnitRegex.FindAllStringSubmatch(value, -1)

	// Checks if value doesn't match a numerical value w/ unit
	if len(matches) == 0 {
		samples = append(samples, &metricSample{Value: value})

		return samples
	}

	num, _ := parseNumber(matches[0][1])

	pv := num
	unit := matches[0][3]

	samples = append(samples, &metricSample{Value: num, Unit: &unit})

	nv, nUnit := normalizeStorageSize(pv, unit)

	if nUnit != unit {
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
