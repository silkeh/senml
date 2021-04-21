package senml

import (
	"math"
	"time"
)

// lcp finds the longest common prefix of the input strings.
// It compares by bytes instead of runes (Unicode code points).
// It's up to the caller to do Unicode normalization if desired
// (e.g. see golang.org/x/text/unicode/norm).
// This function was sourced from Rosetta Code, and licensed under the
// GNU Free Documentation License 1.2.
// See: https://rosettacode.org/wiki/Longest_common_prefix#Go
func lcp(l []string) string {
	// Special cases first
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0]
	}
	// LCP of min and max (lexigraphically)
	// is the LCP of the whole set.
	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return min
}

// maxUnit returns the unit with the greatest value from a map.
func maxUnit(units map[Unit]int) (unit Unit) {
	maxV := 1
	for u, c := range units {
		if c > maxV {
			unit = u
			maxV = c
		}
	}

	return
}

// intToTime converts a 64-bit integer Unix timestamp to time.Time.
func intToTime(t int64) time.Time {
	return time.Unix(t, 0)
}

// floatToTime converts a 64-bit floating point value to time.Time.
func floatToTime(t float64) time.Time {
	s, n := math.Modf(t)
	return time.Unix(int64(s), int64(n*1e9))
}

// floatToDuration converts a 64-bit floating point value to a time.Duration.
func floatToDuration(d float64) time.Duration {
	return time.Duration(d * float64(time.Second))
}

// parseTime converts a Numeric time value and base value to an actual timestamp.
func parseTime(base, val Numeric, now time.Time) (t time.Time) {
	// Convert base time to Time
	baseFloat := numericToFloat64(base)
	if base == nil || baseFloat == 0 {
		t = now
	} else if baseFloat >= (1 << 28) {
		t = numericToTime(base)
	} else {
		t = now.Add(numericToDuration(base))
	}

	// Convert value to Time
	if val == nil {
		return
	}

	if t.IsZero() {
		return numericToTime(val)
	}

	return t.Add(numericToDuration(val))
}

// pow10 returns 10^n with n >=0.
func pow10(n int) (v int) {
	if n < 0 {
		panic("n must be positive")
	}
	v = 1
	for i := 0; i < n; i++ {
		v *= 10
	}
	return
}
