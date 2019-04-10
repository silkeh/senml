package senml

import "testing"

func TestEncodeXMLExamples(t *testing.T) {
	for n, example := range testVectors {
		b, err := EncodeXML(example.Result)
		if err != nil {
			t.Errorf("Error encoding %s: %s", n, err)
			return
		}
		t.Logf("Result for %s: %s", n, b)
	}
}

func TestDecodeXMLExamples(t *testing.T) {
	for n, example := range testVectors {
		if example.XML == "" {
			continue
		}

		res, err := DecodeXML([]byte(example.XML))
		if err != nil {
			t.Errorf("Error decoding %s: %s", n, err)
			continue
		}

		if !equal(res, example.Result) {
			t.Errorf("Decode for example %s incorrect, got:\n%s\nexpected:\n%s", n, toString(res), toString(example.Result))
		}
	}
}

