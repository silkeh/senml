package senml

import "testing"

func TestEncodeCBORExamples(t *testing.T) {
	for n, example := range testVectors {
		b, err := EncodeCBOR(example.Result)
		if err != nil {
			t.Errorf("Error encoding %s: %s", n, err)
			return
		}

		t.Logf("Result for %s: %x", n, b)
	}
}

func TestDecodeCBORExamples(t *testing.T) {
	for n, example := range testVectors {
		if example.CBOR == nil {
			continue
		}

		res, err := DecodeCBOR(example.CBOR)
		if err != nil {
			t.Errorf("Error decoding %s: %s", n, err)
			continue
		}

		if !equal(res, example.Result) {
			t.Errorf("Decode for example %s incorrect, got:\n%s\nexpected:\n%s", n, toString(res), toString(example.Result))
		}
	}
}
