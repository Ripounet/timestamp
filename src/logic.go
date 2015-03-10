package timestamp

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type timeUnit string 

// "interesting" dates are between years 1800 and 2200
const (
	SECONDS_PER_DAY  = 24 * 60 * 60
	SECONDS_PER_YEAR = (365 * SECONDS_PER_DAY) + (SECONDS_PER_DAY / 4)

	TS_SECONDS_MIN = (1800 - 1970) * SECONDS_PER_YEAR
	TS_SECONDS_MAX = (2200 - 1970) * SECONDS_PER_YEAR

	TS_MILLISECONDS_MIN = 1000 * TS_SECONDS_MIN
	TS_MILLISECONDS_MAX = 1000 * TS_SECONDS_MAX

	TS_MICROSECONDS_MIN = 1000 * TS_MILLISECONDS_MIN
	TS_MICROSECONDS_MAX = 1000 * TS_MILLISECONDS_MAX

	TS_NANOSECONDS_MIN = 1000 * TS_MICROSECONDS_MIN
	TS_NANOSECONDS_MAX = 1000 * TS_MICROSECONDS_MAX
)

// parseUnknown tries to guess the input format, and returns the
// parsed date, detected unit, cleaned input string.
func parseUnknown(s string) (time.Time, timeUnit, string, error) {
	// Is it a number?
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		switch {
		// Date-like
		case 18000101 < n && n < 22000101:
			t, err := time.Parse("20060102", s)
			return t, "day", s, err
		case 180001010000 < n && n <= 220001012359:
			t, err := time.Parse("200601021504", s)
			return t, "day hours minutes", s, err
		case 18000101000000 < n && n <= 22000101235959:
			t, err :=  time.Parse("20060102150405", s)
			return t, "day hours minutes seconds", s, err

		// Timestamp-like
		case TS_SECONDS_MIN < n && n < TS_SECONDS_MAX:
			return time.Unix(n, 0), "seconds", s, nil
		case TS_MILLISECONDS_MIN < n && n < TS_MILLISECONDS_MAX:
			return time.Unix(n/1000, (n%1000)*1000000), "milliseconds", s, nil
		case TS_MICROSECONDS_MIN < n && n < TS_MICROSECONDS_MAX:
			return time.Unix(n/1000000, (n%1000000)*1000), "microseconds", s, nil
		case TS_NANOSECONDS_MIN < n && n < TS_NANOSECONDS_MAX:
			return time.Unix(n/1000000000, n%1000000000), "nanoseconds", s, nil
		}
	}
	// Try removing useless chars
	stripped := strip(s)
	if stripped != s {
		tStripped, unit, clean, errStripped := parseUnknown(stripped)
		if errStripped == nil {
			return tStripped, unit, clean, nil
		}
	}
	return time.Time{}, "", s, fmt.Errorf("Could not determine format of input[%v]", s)
}

// Remove all non-digit chars
func strip(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}
