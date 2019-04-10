package senml

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestEncode tests if encoding the expected result gives the same result as decoding the JSON into Objects.
func TestEncode(t *testing.T) {
	for n, test := range testVectors {
		if test.SkipEncode {
			continue
		}

		var exp []Object
		err := json.Unmarshal([]byte(test.JSON), &exp)
		if err != nil {
			t.Errorf("JSON error in test %s: %s", n, err)
			continue
		}

		res := Encode(test.Result)
		if !reflect.DeepEqual(exp, res) {
			t.Errorf("Encode for test %s incorrect, got:\n%#v\nexpected:\n%#v", n, res, exp)
		}
	}
}
