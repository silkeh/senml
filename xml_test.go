package senml

import "testing"

func TestEncodeXMLExamples(t *testing.T) {
	for n, example := range testVectors {
		b, err := EncodeXML(example.Result)
		if err != nil {
			t.Errorf("Error encoding %s: %s", n, err)
			return
		}
		t.Logf("Result for %s (%v bytes): %s", n, len(b), b)
	}
}

func TestDecodeXMLExamples(t *testing.T) {
	AutoTime = false
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

func BenchmarkEncodeXML(b *testing.B) {
	v := "Multiple Measurements"
	ms := testVectors[v].Result
	for n := 0; n < b.N; n++ {
		_, err := EncodeXML(ms)
		if err != nil {
			b.Fatalf("Error encoding %s: %s", v, err)
		}
	}
}

func BenchmarkDecodeXML(b *testing.B) {
	v := "Multiple Data Points 2"
	t := []byte(testVectors[v].XML)
	for n := 0; n < b.N; n++ {
		_, err := DecodeXML(t)
		if err != nil {
			b.Fatalf("Error encoding %s: %s", v, err)
		}
	}
}
