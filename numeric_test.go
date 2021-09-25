package senml

import (
	"encoding/xml"
	"fmt"
	"math"
	"testing"
)

type numericTest struct {
	N Numeric
	F float64
	I int64
	U uint64
	S string
}

var numericTests = []numericTest{
	{N: uint(5), F: 5, I: 5, U: 5, S: "5"},
	{N: uint8(44), F: 44, I: 44, U: 44, S: "44"},
	{N: uint16(11290), F: 11290, I: 11290, U: 11290, S: "11290"},
	{N: uint32(739910425), F: 739910425, I: 739910425, U: 739910425, S: "739910425"},
	{N: uint64(3177891079421588183), F: 3177891079421588183, I: 3177891079421588183, U: 3177891079421588183, S: "3177891079421588183"},
	{N: int(-6), F: -6, I: -6, U: math.MaxUint64 - 5, S: "-6"},
	{N: int8(-45), F: -45, I: -45, U: math.MaxUint64 - 44, S: "-45"},
	{N: int16(-11291), F: -11291, I: -11291, U: math.MaxUint64 - 11290, S: "-11291"},
	{N: int32(-739910426), F: -739910426, I: -739910426, U: math.MaxUint64 - 739910425, S: "-739910426"},
	{N: int64(-3177891079421588184), F: -3177891079421588184, I: -3177891079421588184, U: math.MaxUint64 - 3177891079421588183, S: "-3177891079421588184"},
	{N: float32(8.1171875), F: 8.1171875, I: 8, U: 8, S: "8.1171875"},
	{N: float32(-2.9390623569488525), F: -2.9390623569488525, I: -2, U: math.MaxUint64 - 1, S: "-2.9390624"},
	{N: float64(1.7), F: 1.7, I: 1, U: 1, S: "1.7"},
	{N: NewDecimal(-2, 27115), F: 271.15, I: 271, U: 271, S: "27115e-2"},
	{N: NewDecimal(2, -271), F: -27100, I: -27100, U: math.MaxUint64 - 27099, S: "-271e2"},
	{N: &Decimal{-8, 12345678}, F: 0.12345678, I: 0, U: 0, S: "12345678e-8"},
	{N: &Decimal{8, -123}, F: -123e8, I: -123e8, U: math.MaxUint64 - 123e8 + 1, S: "-123e8"},
}

func TestNumericConversion(t *testing.T) {
	for _, n := range numericTests {
		f := numericToFloat64(n.N)
		i := numericToInt64(n.N)
		u := numericToUint64(n.N)
		if f != n.F {
			t.Errorf("Error converting %#v to float64: got %v", n.N, f)
		}
		if i != n.I {
			t.Errorf("Error converting %#v to int64: got %v", n.N, i)
		}

		if u != n.U {
			t.Errorf("Error converting %#v to uint64: expected %v, got %v", n.N, n.U, u)
		}
	}
}

type testXML struct {
	XMLName xml.Name `json:"-" xml:"s" codec:"-"`
	Attr    Numeric  `json:"a,omitempty" xml:"a,attr,omitempty" codec:"a"`
	Element Numeric  `json:"e,omitempty" xml:"e,omitempty" codec:"e"`
}

func TestNumericEncoding(t *testing.T) {
	for _, n := range numericTests {
		e := fmt.Sprintf(`<s a="%s"><e>%s</e></s>`, n.S, n.S)
		s := testXML{Attr: n.N, Element: n.N}
		x, err := xml.Marshal(s)
		if err != nil {
			t.Errorf("Error encoding %#v to XML: %s", n.N, err)
			continue
		}
		if string(x) != e {
			t.Errorf("Error converting %#v to XML: expected %s, got %s", n.N, e, x)
		}
	}
}

type numericSumTest struct {
	A, B, C Numeric
}

var numericSumTests = []numericSumTest{
	// nil
	{A: nil, B: nil, C: nil},
	{A: nil, B: int(1), C: int(1)},
	{A: int(1), B: nil, C: int(1)},
	{A: nil, B: uint(1), C: uint(1)},
	{A: uint(1), B: nil, C: uint(1)},
	{A: nil, B: float64(1), C: float64(1)},
	{A: float64(1), B: nil, C: float64(1)},
	// int + int
	{A: int(1), B: int(1), C: int64(2)},
	{A: int8(1), B: int8(1), C: int64(2)},
	{A: int16(1), B: int16(1), C: int64(2)},
	{A: int32(1), B: int32(1), C: int64(2)},
	{A: int64(1), B: int64(1), C: int64(2)},
	// int + uint
	{A: int(1), B: uint(1), C: int64(2)},
	{A: int8(1), B: uint8(1), C: int64(2)},
	{A: int16(1), B: uint16(1), C: int64(2)},
	{A: int32(1), B: uint32(1), C: int64(2)},
	{A: int64(1), B: uint64(1), C: int64(2)},
	// int + float
	{A: int(1), B: float32(1), C: float64(2)},
	{A: int(1), B: float64(1), C: float64(2)},
	// uint + uint
	{A: uint(1), B: uint(1), C: uint64(2)},
	{A: uint8(1), B: uint8(1), C: uint64(2)},
	{A: uint16(1), B: uint16(1), C: uint64(2)},
	{A: uint32(1), B: uint32(1), C: uint64(2)},
	{A: uint64(1), B: uint64(1), C: uint64(2)},
	// uint + int
	{A: uint(1), B: int(1), C: int64(2)},
	{A: uint8(1), B: int8(1), C: int64(2)},
	{A: uint16(1), B: int16(1), C: int64(2)},
	{A: uint32(1), B: int32(1), C: int64(2)},
	{A: uint64(1), B: int64(1), C: int64(2)},
	// uint + float
	{A: uint(1), B: float32(1), C: float64(2)},
	{A: uint(1), B: float64(1), C: float64(2)},
	// float + float
	{A: float32(1), B: float32(1), C: float64(2)},
	{A: float64(1), B: float64(1), C: float64(2)},
	// decimal + decimal
	{A: Decimal{0, 1}, B: &Decimal{0, 1}, C: float64(2)},
	{A: &Decimal{0, 1}, B: Decimal{0, 1}, C: float64(2)},
}

func TestNumericSum(t *testing.T) {
	for _, n := range numericSumTests {
		c := sumNumeric(n.A, n.B)
		if c != n.C {
			t.Errorf("Error adding %v (%T) to %v (%T): expected %v (%T) got %v (%T)",
				n.A, n.A, n.B, n.B, n.C, n.C, c, c)
		}
	}
}
