package senml

import (
	"encoding/json"
	"testing"
)

// TestEncode tests if encoding the expected result gives the same result as decoding the JSON into Objects.
func TestEncode(t *testing.T) {
	for n, test := range testVectors {
		if test.SkipEncode {
			continue
		}

		var exp []Record
		err := json.Unmarshal([]byte(test.JSON), &exp)
		if err != nil {
			t.Errorf("JSON error in test %s: %s", n, err)
			continue
		}

		res := Encode(test.Result)
		if len(exp) != len(res) {
			t.Errorf("Encode for test %s incorrect, got:\n%#v\nexpected:\n%#v", n, res, exp)
		}
		for i := range exp {
			if !exp[i].Equal(&res[i]) {
				t.Errorf("Encode for test %s incorrect, got:\n%#v\nexpected:\n%#v", n, res, exp)
				break
			}
		}
	}
}

func TestEncodeCompare(t *testing.T) {
	for n, example := range testVectors {
		c, _ := EncodeCBOR(example.Result)
		j, _ := EncodeJSON(example.Result)
		x, _ := EncodeXML(example.Result)

		t.Logf("Comparison for %s CBOR/JSON/XML (bytes):  %03d/%03d/%03d", n, len(c), len(j), len(x))
	}
}
