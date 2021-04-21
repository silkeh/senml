package senml

import (
	"fmt"
	"math"
	"reflect"
)

// Decimal represents a CBOR decimal type.
// See RFC8949 section 3.4.4.
type Decimal [2]int

// NewDecimal creates a new Decimal value.
func NewDecimal(exponent, value int) Decimal {
	return Decimal{exponent, value}
}

// slice returns the slice for encoding.
func (n Decimal) slice() []int {
	return n[:]
}

// Float returns the floating point representation of the Decimal value.
// Some precision may be lost.
func (n Decimal) Float() float64 {
	return float64(n[1]) * math.Pow10(n[0])
}

// Int returns the integer representation of the Decimal value.
// Any fractional part will be lost.
func (n Decimal) Int() int {
	if n[0] < 0 {
		return n[1] / pow10(-n[0])
	}
	return n[1] * pow10(n[0])
}

// decimalType is used for (de)serialization of the Decimal type.
type decimalType struct{}

// ConvertExt converts a Decimal value to a slice.
func (t *decimalType) ConvertExt(v interface{}) interface{} {
	d, ok := v.(*Decimal)
	if !ok {
		panic(fmt.Sprintf("unsupported type %T (%#v)", v, v))
	}
	return d.slice()
}

// UpdateExt updates the destination with the value from the given array.
func (t *decimalType) UpdateExt(dst interface{}, src interface{}) {
	a, ok := src.([]int)
	if !ok {
		panic(fmt.Sprintf("unsupported type %T (%#v)", src, src))
	}

	if len(a) != 2 {
		panic(fmt.Sprintf("invalid decimal size: %v", len(a)))
	}

	copy(dst.(*Decimal)[:], a)
}

// init registers the interface extension for the decimal type.
func init() {
	err := cbor.SetInterfaceExt(reflect.TypeOf(Decimal{}), 4, new(decimalType))
	if err != nil {
		panic(err)
	}
}
