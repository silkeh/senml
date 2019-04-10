package senml

import "testing"

func TestEncodeJSONExamples(t *testing.T) {
	for n, example := range testVectors {
		b, err := EncodeJSON(example.Result)
		if err != nil {
			t.Errorf("Error encoding %s: %s", n, err)
			return
		}

		t.Logf("Result for %s: %s", n, b)
	}
}


func TestExamplesDecodeJSON(t *testing.T) {
	for n, example := range testVectors {
		res, err := DecodeJSON([]byte(example.JSON))
		if err != nil {
			t.Errorf("Decode error in example %s: %s", n, err)
			continue
		}

		if !equal(res, example.Result) {
			t.Errorf("Decode for example %s incorrect, got:\n%s\nexpected:\n%s", n, toString(res), toString(example.Result))
		}
	}
}

func BenchmarkEncodeJSON(b *testing.B) {
	v := "Multiple Measurements"
	ms := testVectors[v].Result
	for n := 0; n < b.N; n++ {
		_, err := EncodeJSON(ms)
		if err != nil {
			b.Fatalf("Error encoding %s: %s", v, err)
		}
	}
}

func BenchmarkDecodeJSON(b *testing.B) {
	v := "Multiple Data Points 2"
	t := []byte(testVectors[v].JSON)
	for n := 0; n < b.N; n++ {
		_, err := DecodeJSON(t)
		if err != nil {
			b.Fatalf("Error encoding %s: %s", v, err)
		}
	}
}
