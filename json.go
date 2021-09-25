package senml

import (
	"bytes"
	"encoding/json"
)

// EncodeJSON encodes a list of measurements into JSON.
func EncodeJSON(list []Measurement) ([]byte, error) {
	return json.Marshal(Encode(list))
}

// DecodeJSON decodes a list of measurements from JSON.
func DecodeJSON(j []byte) ([]Measurement, error) {
	dec := json.NewDecoder(bytes.NewReader(j))
	dec.UseNumber()

	obj := make([]Record, 0)
	err := dec.Decode(&obj)
	if err != nil {
		return nil, err
	}
	return Decode(obj)
}
