package senml

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/fxamacker/cbor/v2"
)

const (
	maxIntDecimals  = 18
	maxIntScale     = 1e18
	maxUintDecimals = 20
)

// Decimal represents a number in the form of a mantissa and a base-10 scaling factor (m*10^e).
// It implements the decimal fraction from RFC8949 section 3.4.4.
type Decimal struct {
	e, m int64
}

// NewDecimal returns a decimal fraction that represents m*10^e.
func NewDecimal(e, m int64) Decimal {
	return Decimal{m: m, e: e}
}

// Compare returns
//
//	-1 if f is less than g,
//	 0 if f equals g,
//	+1 if f is greater than g.
func (f Decimal) Compare(g Decimal) int {
	fman, frem, gman, grem := f.scaleEqual(g)

	switch {
	case fman < gman:
		return -1
	case fman > gman:
		return 1
	case frem > grem:
		return 1
	case frem < grem:
		return -1
	default:
		return 0
	}
}

func (f Decimal) scaleEqual(g Decimal) (fman, frem, gman, grem int64) {
	if f.e == g.e {
		return f.m, 0, g.m, 0
	}

	if f.e > g.e {
		div, mul, inBounds := f.scaleFactors(f.e, g.e)
		if !inBounds {
			return 1, 0, 0, 0
		}

		return f.m * mul, 0, g.m / div, g.m % div
	}

	div, mul, inBounds := f.scaleFactors(g.e, f.e)
	if !inBounds {
		return 0, 0, 1, 0
	}

	return f.m / div, f.m % div, g.m * mul, 0
}

func (f Decimal) scaleFactors(fe, ge int64) (div, mul int64, inBounds bool) {
	diff := fe - ge

	switch {
	case diff < 0 || diff > 2*maxIntDecimals:
		return 0, 0, false
	case diff > maxIntDecimals:
		return maxIntScale, f.pow10(diff - maxIntDecimals), true
	default:
		return f.pow10(diff), 1, true
	}
}

// Equal returns true if the value in f is the same as g.
func (f Decimal) Equal(g Decimal) bool {
	fman, frem, gman, grem := f.scaleEqual(g)

	return fman == gman && frem == grem
}

// Float returns the [big.Float] equivalent of the decimal fraction.
// The result is infinity when the value exceeds the maximum value of a [big.Float],
// and zero when it is smaller than can be represented.
func (f Decimal) Float() *big.Float {
	m := new(big.Float)
	m.SetInt64(f.m)

	// convert the base 10 exponent into a base 2 fraction and split it.
	exp, rem := math.Modf(float64(f.e) * math.Ln10 / math.Ln2)

	// scale the mantissa using the remainder of the exponent.
	m.Mul(m, big.NewFloat(math.Exp2(rem)))

	// scale the mantissa to the new exponent.
	m.SetMantExp(m, int(exp))

	return m
}

// Float64 returns the float64 equivalent of the decimal fraction.
// The result is infinity when the value exceeds the maximum value of a float64,
// and zero when it is smaller than can be represented.
func (f Decimal) Float64() float64 {
	return float64(f.m) * math.Pow10(int(f.e))
}

// Format converts the decimal fraction to a string in the same format as floating-point values.
// It supports the following formatting directives:
//
//	%e, %g, %s  scientific notation, eg: 1234e+10
//	%E, %G      scientific notation with a capital, eg: 1234E+10
//
// Other formats, such as %b and %x are not implemented.
// It implements [fmt.Formatter].
func (f Decimal) Format(s fmt.State, verb rune) {
	slen, b := f.formatBytes(verb, s.Flag('+'))
	if len(b) == 0 {
		_, _ = fmt.Fprintf(s, "%%!%c(%T=%s)", verb, f, f.String())
	}

	sign := b[:slen]
	b = b[slen:]

	pre, post, zeroes := f.padding(s, len(b)+len(sign))

	writeChar(s, ' ', pre)
	_, _ = s.Write(sign)
	writeChar(s, '0', zeroes)
	_, _ = s.Write(b)
	writeChar(s, ' ', post)
}

func (f Decimal) padding(s fmt.State, length int) (pre, post, zeroes int) {
	if width, ok := s.Width(); ok && width > length {
		switch pad := width - length; {
		case s.Flag('-'):
			post = pad
		case s.Flag('0'):
			zeroes = pad
		default:
			pre = pad
		}
	}
	return
}

func (f Decimal) formatBytes(verb rune, signFlag bool) (slen int, b []byte) {
	sep := f.formatSeparator(verb)
	if sep == "" {
		return 0, nil
	}

	sign, slen := f.appendSign(nil, f.m, signFlag)

	b = make([]byte, 0, 39+len(sign)+len(sep))
	b = append(b, sign...)
	b = strconv.AppendInt(b, f.m, 10)
	b = append(b, sep...)
	b, _ = f.appendSign(b, f.e, true)
	b = strconv.AppendInt(b, f.e, 10)

	return slen, b
}

func (Decimal) formatSeparator(verb rune) (sep string) {
	switch verb {
	case 'e', 'g', 's', 'v':
		return "e"
	case 'E', 'G':
		return "E"
	default:
		return ""
	}
}

func (Decimal) appendSign(b []byte, v int64, flag bool) ([]byte, int) {
	if flag && v >= 0 {
		return append(b, '+'), 1
	}

	if v < 0 {
		return b, 1
	}

	return b, 0
}

// Int returns the [big.Int] representation of the decimal fraction.
// The result is zero for very small positive values and minus one for very small negative values.
// This function is *very* slow for large decimal fractions, which also result in very large integers.
func (f Decimal) Int() *big.Int {
	m := big.NewInt(f.m)
	e := new(big.Int).Abs(big.NewInt(f.e))

	emul := new(big.Int).Exp(big.NewInt(10), e, nil)
	if f.e < 0 {
		return new(big.Int).Div(m, emul)
	}

	return new(big.Int).Mul(m, emul)
}

// Int64 returns the int64 representation of the decimal fraction.
// The result is:
//   - truncated (rounded down) for all normal values.
//   - zero for very small positive values.
//   - minus one for very small negative values.
//   - [math.MaxInt64] for very large positive values.
//   - [math.MinInt64] for very large negative values.
func (f Decimal) Int64() int64 {
	switch {
	case f.e > maxIntDecimals:
		switch {
		case f.m == 0:
			return 0
		case f.m > 0:
			return math.MaxInt64
		default:
			return math.MinInt64
		}
	case f.e < -maxIntDecimals:
		switch {
		case f.m >= 0:
			return 0
		default:
			return -1
		}
	case f.e < 0:
		v := f.m / f.pow10(-f.e)
		if f.m < 0 && v >= 0 {
			return v - 1
		}
		return v
	default:
		v := f.m * f.pow10(f.e)
		switch {
		case f.m > 0 && v < 0:
			return math.MaxInt64
		case f.m < 0 && v > 0:
			return math.MinInt64
		default:
			return v
		}
	}
}

// Uint64 returns the int64 representation of the decimal fraction.
// The result is:
//   - truncated (rounded down) for all normal values.
//   - zero for very small positive values.
//   - [math.MaxUint64] for very large positive values.
func (f Decimal) Uint64() uint64 {
	switch {
	case f.m <= 0:
		return 0
	case f.e > maxUintDecimals:
		return math.MaxUint64
	case f.e < 0:
		return uint64(f.m) / f.pow10u(-f.e)
	default:
		v := uint64(f.m) * f.pow10u(f.e)
		if v == 0 {
			return math.MaxUint64
		}
		return v
	}
}

// IsInt returns true if the decimal fraction can be represented by a [big.Int].
func (f Decimal) IsInt() bool {
	return f.e >= 0
}

// IsInt64 returns true if the decimal fraction can be represented by an int64.
func (f Decimal) IsInt64() bool {
	if f.e < 0 || f.e > maxIntDecimals {
		return false
	}

	v := int64(math.MaxInt64)

	for i := range int64(maxIntDecimals) {
		if f.e == i {
			return f.m <= v
		}

		v /= 10
	}

	return f.m <= v
}

// MantExp returns the mantissa and the base-10 exponent of the decimal fraction.
func (f Decimal) MantExp() (m int64, e int64) {
	return f.m, f.e
}

// String returns the string representation fo the decimal fraction.
func (f Decimal) String() string {
	_, b := f.formatBytes('e', false)

	return string(b)
}

// MarshalCBOR implements [cbor.Marshaler].
func (f Decimal) MarshalCBOR() ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte(0xc4)

	if err := cbor.NewEncoder(buffer).Encode([2]int64{f.e, f.m}); err != nil {
		return nil, fmt.Errorf("senml decimal: %w", err)
	}

	return buffer.Bytes(), nil
}

// MarshalJSON implements [encoding/json.Marshaler].
func (f Decimal) MarshalJSON() ([]byte, error) {
	return f.MarshalText()
}

// MarshalText implements [encoding.TextMarshaler].
func (f Decimal) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "%de%d", f.m, f.e), nil
}

// UnmarshalCBOR implements [Unmarshaler].
func (f *Decimal) UnmarshalCBOR(data []byte) error {
	if len(data) >= 1 && data[0] == 0xc4 {
		data = data[1:]
	}

	parts := [2]int64{}

	if err := cbor.Unmarshal(data, &parts); err != nil {
		return fmt.Errorf("senml decimal: %w", err)
	}

	*f = Decimal{e: parts[0], m: parts[1]}

	return nil
}

// UnmarshalJSON implements [encoding/json.Unmarshaler].
func (f *Decimal) UnmarshalJSON(data []byte) error {
	return f.UnmarshalText(data)
}

// UnmarshalText implements [encoding.TextUnmarshaler].
func (f *Decimal) UnmarshalText(data []byte) error {
	_, err := fmt.Sscanf(string(bytes.ToLower(data)), "%de%d", &f.m, &f.e)
	if err != nil {
		return fmt.Errorf("unmarshal decimal fraction from JSON: %w", err)
	}

	return nil
}

// matchScale matches the exponent of two numbers.
func (Decimal) matchScale(a, b Decimal) (x, y Decimal) {
	switch {
	case a.e < b.e:
		return
	case a.e > b.e:
		return
	default:
		return a, b
	}
}

// pow10 returns 10^n with n >=0.
// It is undefined for large powers.
func (Decimal) pow10(n int64) (v int64) {
	if n < 0 {
		panic("n must be positive")
	}

	if n > maxIntDecimals {
		return -1
	}

	v = 1

	for i := int64(0); i < n; i++ {
		v *= 10
	}

	return
}

func (Decimal) pow10u(n int64) (v uint64) {
	if n < 0 {
		panic("n must be positive")
	}

	if n > maxUintDecimals {
		return 0
	}

	v = 1

	for i := int64(0); i < n; i++ {
		v *= 10
	}

	return
}

func writeChar(s fmt.State, char byte, count int) {
	for i := 0; i < count; i++ {
		_, _ = s.Write([]byte{char})
	}
}
