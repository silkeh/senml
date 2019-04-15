package senml

import (
	"math"
	"time"
)

// GNU Free Documentation License 1.2

// lcp finds the longest common prefix of the input strings.
// It compares by bytes instead of runes (Unicode code points).
// It's up to the caller to do Unicode normalization if desired
// (e.g. see golang.org/x/text/unicode/norm).
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

func floatToTime(t float64) time.Time {
	s, n := math.Modf(t)
	return time.Unix(int64(s), int64(n*1e9))
}

func timeToFloat(t time.Time) float64 {
	if t.IsZero() {
		return 0
	}
	return float64(t.UnixNano()) / 1e9
}

func floatToDuration(d float64) time.Duration {
	return time.Duration(d * float64(time.Second))
}

func parseTime(base, val float64, now time.Time) (t time.Time) {
	// Convert base time to Time
	if base == 0 {
		t = now
	} else if base >= (1 << 28) {
		t = floatToTime(base)
	} else {
		t = now.Add(floatToDuration(base))
	}

	// Convert value to Time
	if val == 0 {
		return
	} else if t.IsZero() {
		t = floatToTime(val)
	} else {
		t = t.Add(floatToDuration(val))
	}

	return
}
