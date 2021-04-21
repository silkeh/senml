package senml

import "testing"

func TestEncodeCBORExamples(t *testing.T) {
	for n, example := range testVectors {
		b, err := EncodeCBOR(example.Result)
		if err != nil {
			t.Errorf("Error encoding %s: %s", n, err)
			return
		}

		t.Logf("Result for %s (%v bytes): %x", n, len(b), b)
	}
}

func TestDecodeCBORExamples(t *testing.T) {
	AutoTime = false
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

func BenchmarkEncodeCBOR(b *testing.B) {
	v := "Multiple Measurements"
	ms := testVectors[v].Result
	for n := 0; n < b.N; n++ {
		_, err := EncodeCBOR(ms)
		if err != nil {
			b.Fatalf("Error encoding %s: %s", v, err)
		}
	}
}

func BenchmarkDecodeCBOR(b *testing.B) {
	v := "Multiple Data Points 2"
	t := testVectors[v].CBOR
	for n := 0; n < b.N; n++ {
		_, err := DecodeCBOR(t)
		if err != nil {
			b.Fatalf("Error encoding %s: %s", v, err)
		}
	}
}
