package senml_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/silkeh/senml"
)

func ExampleDecimal_Compare() {
	numbers := []senml.Decimal{
		senml.NewDecimal(0, 0),
		senml.NewDecimal(-1, 2),
		senml.NewDecimal(1, 3),
	}

	for _, a := range numbers {
		for _, b := range numbers {
			switch a.Compare(b) {
			case -1:
				fmt.Printf("%e < %e\n", a, b)
			case 0:
				fmt.Printf("%e = %e\n", a, b)
			case 1:
				fmt.Printf("%e > %e\n", a, b)
			}

		}
	}

	// Output:
	// 0e+0 = 0e+0
	// 0e+0 < 2e-1
	// 0e+0 < 3e+1
	// 2e-1 > 0e+0
	// 2e-1 = 2e-1
	// 2e-1 < 3e+1
	// 3e+1 > 0e+0
	// 3e+1 > 2e-1
	// 3e+1 = 3e+1
}

func ExampleDecimal_Equal() {
	numbers := []senml.Decimal{
		senml.NewDecimal(1, 1),
		senml.NewDecimal(-1, 2),
		senml.NewDecimal(0, 10),
	}

	for _, a := range numbers {
		for _, b := range numbers {
			fmt.Printf("%e == %e ? %t\n", a, b, a.Equal(b))
		}
	}

	// Output:
	// 1e+1 == 1e+1 ? true
	// 1e+1 == 2e-1 ? false
	// 1e+1 == 10e+0 ? true
	// 2e-1 == 1e+1 ? false
	// 2e-1 == 2e-1 ? true
	// 2e-1 == 10e+0 ? false
	// 10e+0 == 1e+1 ? true
	// 10e+0 == 2e-1 ? false
	// 10e+0 == 10e+0 ? true
}

func TestDecimal_Compare(t *testing.T) {
	tests := map[string]struct {
		F, G senml.Decimal
		Exp  int
	}{
		"equal": {
			F:   senml.NewDecimal(17, 1),
			G:   senml.NewDecimal(17, 1),
			Exp: 0,
		},
		"same scale larger": {
			F:   senml.NewDecimal(17, 2),
			G:   senml.NewDecimal(17, 1),
			Exp: 1,
		},
		"same scale smaller": {
			F:   senml.NewDecimal(17, 3),
			G:   senml.NewDecimal(17, 4),
			Exp: -1,
		},
		"min/max int": {
			F:   senml.NewDecimal(math.MinInt64, 1),
			G:   senml.NewDecimal(math.MaxInt64, 1),
			Exp: -1,
		},
		"max/min int": {
			F:   senml.NewDecimal(math.MaxInt64, 1),
			G:   senml.NewDecimal(math.MinInt32, 1),
			Exp: 1,
		},
		"scaled equal": {
			F:   senml.NewDecimal(10, 10),
			G:   senml.NewDecimal(11, 1),
			Exp: 0,
		},
		"scaled not equal": {
			F:   senml.NewDecimal(10, 11),
			G:   senml.NewDecimal(11, 1),
			Exp: 1,
		},
		"max int64 scale": {
			F:   senml.NewDecimal(18, 9),
			G:   senml.NewDecimal(0, 9e18),
			Exp: 0,
		},
		"max int64 diff": {
			F:   senml.NewDecimal(18, 9),
			G:   senml.NewDecimal(0, math.MaxInt64),
			Exp: -1,
		},
		"exceed max int64 diff": {
			F:   senml.NewDecimal(19, 1),
			G:   senml.NewDecimal(0, math.MaxInt64),
			Exp: 1,
		},
		"actual max": {
			F:   senml.NewDecimal(36, 1),
			G:   senml.NewDecimal(0, math.MaxInt64),
			Exp: 1,
		},
		"exceed max range": {
			F:   senml.NewDecimal(37, 1),
			G:   senml.NewDecimal(0, math.MaxInt64),
			Exp: 1,
		},
	}

	for n, test := range tests {
		t.Run(n, func(t *testing.T) {
			if got := test.F.Compare(test.G); got != test.Exp {
				t.Errorf("want %d, got %d", test.Exp, got)
			}
		})
	}
}

func ExampleDecimal_Float() {
	numbers := []senml.Decimal{
		senml.NewDecimal(0, 1),
		senml.NewDecimal(-1, 11),
		// Largest valid positive number:
		senml.NewDecimal(646456974, 8808063447118036156),
		// Smallest valid positive non-zero number:
		senml.NewDecimal(-646457012, 2838308228594544048),
		// Largest valid negative number:
		senml.NewDecimal(-646457012, -2838308228594544048),
		// Smallest valid negative number:
		senml.NewDecimal(646456974, -8808063447118036156),
		// Too large
		senml.NewDecimal(math.MaxInt32, 1),
		// Too small
		senml.NewDecimal(-math.MaxInt32, 1),
	}

	for _, n := range numbers {
		dm, de := n.MantExp()
		fm := n.Float()
		fe := fm.MantExp(fm)

		fmt.Printf("%d * 10^%d =~ %g * 2^%d\n", dm, de, fm, fe)
	}

	// Output:
	// 1 * 10^0 =~ 0.5 * 2^1
	// 11 * 10^-1 =~ 0.5499999999999998779 * 2^1
	// 8808063447118036156 * 10^646456974 =~ 0.9999999999999999999 * 2^2147483647
	// 2838308228594544048 * 10^-646457012 =~ 0.5000000000000000001 * 2^-2147483648
	// -2838308228594544048 * 10^-646457012 =~ -0.5000000000000000001 * 2^-2147483648
	// -8808063447118036156 * 10^646456974 =~ -0.9999999999999999999 * 2^2147483647
	// 1 * 10^2147483647 =~ +Inf * 2^0
	// 1 * 10^-2147483647 =~ 0 * 2^0
}

func ExampleDecimal_Float64() {
	numbers := []senml.Decimal{
		senml.NewDecimal(0, 1),
		senml.NewDecimal(-1, 11),
		// Largest valid positive number:
		senml.NewDecimal(292, 17976931348623157),
		// Smallest valid positive non-zero number:
		senml.NewDecimal(-323, 1),
		// Largest valid negative number:
		senml.NewDecimal(-323, -1),
		// Smallest valid negative number:
		senml.NewDecimal(292, -17976931348623157),
		// Too large
		senml.NewDecimal(400, 1),
		// Too small
		senml.NewDecimal(-400, 1),
	}

	for _, n := range numbers {
		fmt.Printf("%g =~ %g\n", n, n.Float64())
	}

	// Output:
	// 1e+0 =~ 1
	// 11e-1 =~ 1.1
	// 17976931348623157e+292 =~ 1.7976931348623157e+308
	// 1e-323 =~ 1e-323
	// -1e-323 =~ -1e-323
	// -17976931348623157e+292 =~ -1.7976931348623157e+308
	// 1e+400 =~ +Inf
	// 1e-400 =~ 0
}

func ExampleDecimal_Format() {
	formatters := []string{
		"%e", "%E", "%g", "%G",
		"%+e", "%+E", "%+g", "%+G",
		"%9e", "%09e", "%-9e|",
		"%o", "%s",
	}

	for _, f := range formatters {
		fmt.Printf("%-5s  "+f+"\n", f, senml.NewDecimal(20, 1234))
	}

	// Output:
	// %e     1234e+20
	// %E     1234E+20
	// %g     1234e+20
	// %G     1234E+20
	// %+e    +1234e+20
	// %+E    +1234E+20
	// %+g    +1234e+20
	// %+G    +1234E+20
	// %9e     1234e+20
	// %09e   01234e+20
	// %-9e|  1234e+20 |
	// %o     %!o(senml.Decimal=1234e+20)
	// %s     1234e+20
}

func ExampleDecimal_Int() {
	numbers := []senml.Decimal{
		senml.NewDecimal(0, 1),
		senml.NewDecimal(-1, 11),
		senml.NewDecimal(24, 1),
		// Too small to be represented:
		senml.NewDecimal(-1, 1),
		senml.NewDecimal(-1, -1),
	}

	for _, n := range numbers {
		fmt.Printf("%s =~ %d\n", n, n.Int())
	}

	// Output:
	// 1e+0 =~ 1
	// 11e-1 =~ 1
	// 1e+24 =~ 1000000000000000000000000
	// 1e-1 =~ 0
	// -1e-1 =~ -1
}

func ExampleDecimal_Int64() {
	numbers := []senml.Decimal{
		senml.NewDecimal(0, 1),
		senml.NewDecimal(-1, 11),
		senml.NewDecimal(-18, math.MaxInt64),
		senml.NewDecimal(-18, math.MinInt64),
		senml.NewDecimal(100, 0),
		// Too small to be represented:
		senml.NewDecimal(-1, 1),
		senml.NewDecimal(-1, -1),
		senml.NewDecimal(-19, 1),
		senml.NewDecimal(-19, -1),
		// Too large to be represented:
		senml.NewDecimal(1, 922337203685477581),
		senml.NewDecimal(1, -922337203685477581),
		senml.NewDecimal(19, 1),
		senml.NewDecimal(19, -1),
	}

	for _, n := range numbers {
		fmt.Printf("%e =~ %d\n", n, n.Int64())
	}

	// Output:
	// 1e+0 =~ 1
	// 11e-1 =~ 1
	// 9223372036854775807e-18 =~ 9
	// -9223372036854775808e-18 =~ -9
	// 0e+100 =~ 0
	// 1e-1 =~ 0
	// -1e-1 =~ -1
	// 1e-19 =~ 0
	// -1e-19 =~ -1
	// 922337203685477581e+1 =~ 9223372036854775807
	// -922337203685477581e+1 =~ -9223372036854775808
	// 1e+19 =~ 9223372036854775807
	// -1e+19 =~ -9223372036854775808
}

func ExampleDecimal_Uint64() {
	numbers := []senml.Decimal{
		senml.NewDecimal(0, 1),
		senml.NewDecimal(-1, 11),
		senml.NewDecimal(-18, math.MaxInt64),
		senml.NewDecimal(-18, math.MinInt64),
		senml.NewDecimal(100, 0),
		// Too small to be represented:
		senml.NewDecimal(-1, 1),
		senml.NewDecimal(-19, 1),
		// Too large to be represented:
		senml.NewDecimal(1, 922337203685477581),
		senml.NewDecimal(21, 1),
	}

	for _, n := range numbers {
		fmt.Printf("%e =~ %d\n", n, n.Uint64())
	}

	// Output:
	// 1e+0 =~ 1
	// 11e-1 =~ 1
	// 9223372036854775807e-18 =~ 9
	// -9223372036854775808e-18 =~ 0
	// 0e+100 =~ 0
	// 1e-1 =~ 0
	// 1e-19 =~ 0
	// 922337203685477581e+1 =~ 9223372036854775810
	// 1e+21 =~ 18446744073709551615
}

func ExampleDecimal_IsInt() {
	numbers := []senml.Decimal{
		senml.NewDecimal(0, 1),
		senml.NewDecimal(-1, 11),
		senml.NewDecimal(10, 11),
	}

	for _, n := range numbers {
		fmt.Printf("%e is int: %t\n", n, n.IsInt())
	}

	// Output:
	// 1e+0 is int: true
	// 11e-1 is int: false
	// 11e+10 is int: true
}

func ExampleDecimal_IsInt64() {
	numbers := []senml.Decimal{
		senml.NewDecimal(0, 1),
		senml.NewDecimal(-1, 11),
		senml.NewDecimal(18, 9),
		senml.NewDecimal(17, 93),
		senml.NewDecimal(0, math.MaxInt64),
	}

	for _, n := range numbers {
		fmt.Printf("%e is int64: %t\n", n, n.IsInt64())
	}

	// Output:
	// 1e+0 is int64: true
	// 11e-1 is int64: false
	// 9e+18 is int64: true
	// 93e+17 is int64: false
	// 9223372036854775807e+0 is int64: true
}

func ExampleDecimal_MarshalCBOR() {
	f := senml.NewDecimal(-2, 27315)

	data, _ := cbor.Marshal(f)
	fmt.Printf("% 02x\n", data)
	// Output:
	// c4 82 21 19 6a b3
}

func BenchmarkDecimal_MarshalCBOR(b *testing.B) {
	v := senml.NewDecimal(-2, 27315)

	for b.Loop() {
		_, _ = cbor.Marshal(v)
	}
}

func ExampleDecimal_MarshalJSON() {
	f := senml.NewDecimal(-2, 27315)

	data, _ := json.Marshal(f)
	fmt.Println(string(data))
	// Output:
	// 27315e-2
}

func ExampleDecimal_MarshalText() {
	f := senml.NewDecimal(-2, 27315)

	data, _ := f.MarshalText()
	fmt.Println(string(data))
	// Output:
	// 27315e-2
}

func ExampleDecimal_UnmarshalCBOR() {
	var f senml.Decimal

	cbor.Unmarshal([]byte{0xC4, 0x82, 0x21, 0x19, 0x6a, 0xb3}, &f)

	fmt.Println(f)
	// Output:
	// 27315e-2
}

func BenchmarkDecimal_UnmarshalCBOR(b *testing.B) {
	for b.Loop() {
		v := senml.Decimal{}
		_ = cbor.Unmarshal([]byte{0xC4, 0x82, 0x21, 0x19, 0x6a, 0xb3}, &v)
	}
}

func TestDecimal_UnmarshalCBOR(t *testing.T) {
	var f senml.Decimal

	err := cbor.Unmarshal([]byte{0x82, 0x21, 0x19, 0x6a, 0xb3}, &f)
	if err != nil {
		t.Fatal(err)
	}

	if f.Compare(senml.NewDecimal(-2, 27315)) != 0 {
		t.Errorf("expected 27315e-1, got %e", f)
	}
}

func ExampleDecimal_UnmarshalJSON() {
	var f senml.Decimal

	json.Unmarshal([]byte("27315e-2"), &f)

	fmt.Println(f)
	// Output:
	// 27315e-2
}

func ExampleDecimal_UnmarshalText() {
	var f senml.Decimal

	f.UnmarshalText([]byte("27315e-2"))

	fmt.Println(f)
	// Output:
	// 27315e-2
}

func TestDecimal_cbor(t *testing.T) {
	tests := map[string]struct {
		Value senml.Decimal
		CBOR  []byte
	}{
		"zero": {
			Value: senml.Decimal{},
			CBOR:  []byte{0xc4, 0x82, 0x00, 0x00},
		},
		"rfc8949": {
			Value: senml.NewDecimal(-2, 27315),
			CBOR:  []byte{0xc4, 0x82, 0x21, 0x19, 0x6a, 0xb3},
		},
		"1.1": {
			Value: senml.NewDecimal(-1, 11),
			CBOR:  []byte{0xc4, 0x82, 0x20, 0x0b},
		},
		"1": {
			Value: senml.NewDecimal(0, 1),
			CBOR:  []byte{0xc4, 0x82, 0x00, 0x01},
		},
		"int8": {
			Value: senml.NewDecimal(math.MinInt8, math.MaxInt8),
			CBOR:  []byte{0xc4, 0x82, 0x38, 0x7f, 0x18, 0x7f},
		},
		"int16": {
			Value: senml.NewDecimal(math.MinInt16, math.MaxInt16),
			CBOR:  []byte{0xc4, 0x82, 0x39, 0x7f, 0xff, 0x19, 0x7f, 0xff},
		},
		"int32": {
			Value: senml.NewDecimal(math.MinInt64, math.MaxInt32),
			CBOR: []byte{
				0xc4, 0x82, 0x3b, 0x7f, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0x1a, 0x7f, 0xff, 0xff, 0xff,
			},
		},
		"int64": {
			Value: senml.NewDecimal(math.MinInt64, math.MaxInt64),
			CBOR: []byte{
				0xc4, 0x82, 0x3b, 0x7f, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0x1b, 0x7f, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff},
		},
	}

	for n, test := range tests {
		t.Run(n, func(t *testing.T) {
			gotCBOR, err := cbor.Marshal(test.Value)
			if err != nil {
				t.Errorf("marshal: %s", err)
			}

			if !bytes.Equal(gotCBOR, test.CBOR) {
				t.Errorf("marshal: want % 02x, got % 02x", test.CBOR, gotCBOR)
			}

			gotValue := senml.Decimal{}

			err = cbor.Unmarshal(test.CBOR, &gotValue)
			if err != nil {
				t.Errorf("unmarshal: %s", err)
			}

			if test.Value.Compare(gotValue) != 0 {
				t.Errorf("unmarshal: want %e, got %e", test.Value, gotValue)
			}
		})
	}
}
