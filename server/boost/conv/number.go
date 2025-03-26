package conv

import (
	"math"
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

func Float64ToString(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}

func UIntToString(n uint) string { return strconv.FormatUint(uint64(n), 10) }

func UInt64ToString(n uint64) string { return strconv.FormatUint(n, 10) }

func Int64ToString(n int64) string { return strconv.FormatInt(n, 10) }

func StringToFloat64[T string | *string](s T) *float64 {
	var str string
	switch v := (any)(s).(type) {
	case string:
		str = v
	case *string:
		if v != nil {
			str = *v
		}
	}
	if len(str) == 0 {
		return nil
	}
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil
	}
	return &v
}

func StringToBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BooleanToInt[T constraints.Integer, K bool | *bool](source K) T {
	var boolean bool
	switch v := (any)(source).(type) {
	case bool:
		boolean = v
	case *bool:
		if v != nil {
			boolean = *v
		}
	}
	if boolean {
		return T(1)
	}
	return T(0)
}

func IntToBoolean[T constraints.Integer](t T) bool { return t != 0 }

func NumberToString[T constraints.Integer | constraints.Float](n T) string {
	switch any(n).(type) {
	case float32, float64:
		return Float64ToString(float64(n))
	default:
		return Int64ToString(int64(n))
	}
}

func CeilToHalf[T constraints.Float](v T) float64 {
	if v <= 0 {
		return exp.Zero[float64]()
	}

	return math.Ceil(float64(v)*2) / 2
}

func StringToNumber[R constraints.Integer | constraints.Float](source string) *R {
	if len(source) == 0 {
		return nil
	}

	var zero, result R
	switch any(zero).(type) {
	case float32, float64:
		v, err := strconv.ParseFloat(source, 64)
		if err != nil {
			return nil
		}
		result = R(v)
	default:
		v, err := strconv.ParseInt(source, 10, 64)
		if err != nil {
			return nil
		}
		result = R(v)
	}
	return &result
}
