package conv

import (
	"strconv"
	"unsafe"

	"golang.org/x/exp/constraints"
)

func Float64ToString(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}

func UIntToString(n uint) string { return strconv.FormatUint(uint64(n), 10) }

func UInt64ToString(n uint64) string { return strconv.FormatUint(n, 10) }

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
