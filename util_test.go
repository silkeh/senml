package senml

import (
	"fmt"
	"strings"
)

func toString(ml []Measurement) string {
	strs := make([]string, len(ml))
	for i, m := range ml {
		strs[i] = fmt.Sprintf("%#v", m)
	}

	return strings.Join(strs, ",\n")
}

func equal(a, b []Measurement) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}
