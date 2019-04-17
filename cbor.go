package senml

import "github.com/ugorji/go/codec"

var cbor codec.CborHandle

// EncodeCBOR encodes a list of measurements into CBOR.
func EncodeCBOR(list []Measurement) (b []byte, err error) {
	err = codec.NewEncoderBytes(&b, &cbor).Encode(Encode(list))
	return
}

// DecodeCBOR decodes a list of measurements from CBOR.
func DecodeCBOR(c []byte) ([]Measurement, error) {
	obj := make([]Record, 0)
	err := codec.NewDecoderBytes(c, &cbor).Decode(&obj)
	if err != nil {
		return nil, err
	}
	return Decode(obj)
}
