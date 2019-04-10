package senml

import (
	"fmt"
	"log"
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
	for i := range a {
		if !a[i].Equal(b[i]) {
			log.Printf("\n%#v\n%#v", a[i], b[i])
			return false
		}
	}
	return true
}
