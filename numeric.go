package senml

import (
	"encoding/json"
	"fmt"
	"time"
)

// Numeric represents a numeric value. This can be any integer or floating point
// type (including json.Number), or a Decimal fraction.
// Numeric is equal to the empty interface, but using it for anything other than
// those types will result in a panic.
type Numeric interface{}

// sumNumeric adds two Numeric types and returns the sum.
func sumNumeric(a, b Numeric) Numeric {
	switch v := a.(type) {
	case nil:
		return b
	case int, int8, int16, int32, int64:
		return addToNumericInt(a, b)
	case uint, uint8, uint16, uint32, uint64:
		return addToNumericUint(a, b)
	case float32, float64, Decimal, *Decimal:
		return numericToFloat64(a) + numericToFloat64(b)
	case json.Number:
		return sumNumeric(jsonToInt64OrFloat64(v), b)
	default:
		panic(fmt.Sprintf("invalid value type: %T", a))
	}
}

// compareNumeric returns true if a and b are equal.
func compareNumeric(a, b Numeric) bool {
	switch x := a.(type) {
	case nil:
		return b == nil

	case int, int8, int16, int32, int64:
		switch b.(type) {
		case int, int8, int16, int32, int64:
			return numericToInt64(a) == numericToInt64(b)
		case uint, uint8, uint16, uint32, uint64:
			if numericToInt64(b) < 0 {
				return false
			}
			return numericToInt64(a) == numericToInt64(b)
		default:
			return false
		}

	case uint, uint8, uint16, uint32, uint64:
		switch b.(type) {
		case int, int8, int16, int32, int64:
			if numericToInt64(b) < 0 {
				return false
			}
			return numericToUint64(a) == numericToUint64(b)
		case uint, uint8, uint16, uint32, uint64:
			return numericToUint64(a) == numericToUint64(b)
		default:
			return false
		}

	case float32, float64:
		return numericToFloat64(a) == numericToFloat64(b)

	case Decimal:
		switch y := b.(type) {
		case Decimal:
			return x.Equal(y)
		case *Decimal:
			return x.Equal(*y)
		default:
			return numericToFloat64(a) == numericToFloat64(b)
		}

	case json.Number:
		switch y := b.(type) {
		case json.Number:
			return x == y
		default:
			return numericToFloat64(a) == numericToFloat64(b)
		}
	default:
		panic(fmt.Sprintf("invalid value type: %T", a))
	}
}

// addToNumericInt sums two numeric values, the first of which is an integer type.
// An attempt is made to preserve integer types, but this is not guaranteed.
func addToNumericInt(a, b Numeric) Numeric {
	switch b.(type) {
	case nil:
		return a
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return numericToInt64(a) + numericToInt64(b)
	case float32, float64, Decimal, *Decimal:
		return numericToFloat64(a) + numericToFloat64(b)
	default:
		panic(fmt.Sprintf("invalid value type: %T", b))
	}
}

// addToNumericInt sums two numeric values, the first of which is an unsigned integer type.
// An attempt is made to preserve unsigned integer types, but this is not guaranteed.
func addToNumericUint(a, b Numeric) Numeric {
	switch b.(type) {
	case nil:
		return a
	case int, int8, int16, int32, int64:
		return numericToInt64(a) + numericToInt64(b)
	case uint, uint8, uint16, uint32, uint64:
		return numericToUint64(a) + numericToUint64(b)
	case float32, float64, Decimal, *Decimal:
		return numericToFloat64(a) + numericToFloat64(b)
	default:
		panic(fmt.Sprintf("invalid value type: %T", b))
	}
}

// timeToNumeric converts a time.Time to a Numeric value.
func timeToNumeric(t time.Time) Numeric {
	if t.IsZero() {
		return nil
	}
	if t.Nanosecond() == 0 {
		return t.Unix()
	}
	return float64(t.UnixNano()) / 1e9
}

// numericToTime converts a Numeric value to a time.Time.
func numericToTime(v Numeric) time.Time {
	switch n := v.(type) {
	case nil:
		return time.Time{}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return intToTime(numericToInt64(v))
	case float32, float64, Decimal, *Decimal:
		return floatToTime(numericToFloat64(v))
	case json.Number:
		return numericToTime(jsonToInt64OrFloat64(n))
	default:
		panic("unsupported type")
	}
}

// numericToDuration converts a Numeric value to a duration.
func numericToDuration(v Numeric) time.Duration {
	switch n := v.(type) {
	case nil:
		return 0
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return time.Duration(numericToInt64(v)) * time.Second
	case float32, float64, Decimal, *Decimal:
		return floatToDuration(numericToFloat64(v))
	case json.Number:
		return numericToDuration(jsonToInt64OrFloat64(n))
	default:
		panic("unsupported type")
	}
}

// numericToInt64 converts a numeric value to a 64-bit integer.
// Fractional parts will be lost in the conversion.
func numericToInt64(v Numeric) int64 {
	switch i := v.(type) {
	case nil:
		return 0
	case int:
		return int64(i)
	case int8:
		return int64(i)
	case int16:
		return int64(i)
	case int32:
		return int64(i)
	case int64:
		return i
	case uint:
		return int64(i)
	case uint8:
		return int64(i)
	case uint16:
		return int64(i)
	case uint32:
		return int64(i)
	case uint64:
		return int64(i)
	case float32:
		return int64(i)
	case float64:
		return int64(i)
	case Decimal:
		return int64(i.Int())
	case *Decimal:
		return int64(i.Int())
	case json.Number:
		v, _ := i.Int64()
		return v
	default:
		panic(fmt.Sprintf("invalid value type: %T", i))
	}
}

// numericToUint64 converts a Numeric value to a 64-bit unsigned integer.
// Fractional parts and signs will be lost in the conversion.
func numericToUint64(v Numeric) uint64 {
	switch i := v.(type) {
	case nil:
		return 0
	case int:
		return uint64(i)
	case int8:
		return uint64(i)
	case int16:
		return uint64(i)
	case int32:
		return uint64(i)
	case int64:
		return uint64(i)
	case uint:
		return uint64(i)
	case uint8:
		return uint64(i)
	case uint16:
		return uint64(i)
	case uint32:
		return uint64(i)
	case uint64:
		return i
	case float32:
		return uint64(i)
	case float64:
		return uint64(i)
	case Decimal:
		return uint64(i.Int())
	case *Decimal:
		return uint64(i.Int())
	case json.Number:
		v, _ := i.Int64()
		return uint64(v)
	default:
		panic(fmt.Sprintf("invalid value type: %T", i))
	}
}

// numericToFloat64 converts a Numeric value to a 64-bit floating point value.
// Precision may be lost in the conversion.
func numericToFloat64(v Numeric) float64 {
	switch f := v.(type) {
	case nil:
		return 0
	case int:
		return float64(f)
	case int8:
		return float64(f)
	case int16:
		return float64(f)
	case int32:
		return float64(f)
	case int64:
		return float64(f)
	case uint:
		return float64(f)
	case uint8:
		return float64(f)
	case uint16:
		return float64(f)
	case uint32:
		return float64(f)
	case uint64:
		return float64(f)
	case float32:
		return float64(f)
	case float64:
		return f
	case Decimal:
		return f.Float()
	case *Decimal:
		return f.Float()
	case json.Number:
		v, _ := f.Float64()
		return v
	default:
		panic(fmt.Sprintf("invalid value type: %T", f))
	}
}

// jsonToInt64OrFloat64 returns a json.Number as either an int64 or float64.
func jsonToInt64OrFloat64(n json.Number) Numeric {
	if i, err := n.Int64(); err == nil {
		return i
	}
	f, _ := n.Float64()
	return f
}
