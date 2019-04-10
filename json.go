package senml

import "encoding/json"

// EncodeJSON encodes a list of measurements into JSON.
func EncodeJSON(list []Measurement) ([]byte, error) {
	return json.Marshal(Encode(list))
}

// DecodeJSON decodes a list of measurements from JSON.
func DecodeJSON(j []byte) ([]Measurement, error) {
	obj := make([]Object, 0)
	err := json.Unmarshal(j, &obj)
	if err != nil {
		return nil, err
	}
	return Decode(obj)
}
