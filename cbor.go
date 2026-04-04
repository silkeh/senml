package senml

import (
	"fmt"
	"reflect"

	"github.com/fxamacker/cbor/v2"
)

var cborTags cbor.TagSet

func init() {
	cborTags = cbor.NewTagSet()
	cborTags.Add(
		cbor.TagOptions{EncTag: cbor.EncTagRequired, DecTag: cbor.DecTagRequired},
		reflect.TypeOf(Decimal{}),
		4,
	)
}

// EncodeCBOR encodes a list of measurements into CBOR.
func EncodeCBOR(list []Measurement) (b []byte, err error) {
	encode, _ := cbor.EncOptions{}.EncModeWithTags(cborTags)

	b, err = encode.Marshal(Encode(list))
	if err != nil {
		return nil, fmt.Errorf("senml cbor: %w", err)
	}

	return
}

// DecodeCBOR decodes a list of measurements from CBOR.
func DecodeCBOR(c []byte) ([]Measurement, error) {
	var records []Record

	decode, _ := cbor.DecOptions{}.DecModeWithTags(cborTags)

	err := decode.Unmarshal(c, &records)
	if err != nil {
		return nil, fmt.Errorf("senml cbor: %w", err)
	}

	return Decode(records)
}
