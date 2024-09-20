package numbers

import (
	"onij/boost/conv"
	"strconv"
	"strings"
)

type CommaNumber[T Number] string

func NewCommaNumber[T Number](s string) CommaNumber[T] { return CommaNumber[T](s) }

func JoinComma[T Number](source []T, wrap bool) CommaNumber[T] {
	strs := make([]string, len(source))
	for i, v := range source {
		s, f := Which(v)
		if s {
			strs[i] = strconv.Itoa(int(v))
			continue
		}
		if f {
			strs[i] = conv.Float64ToString(float64(v))
			continue
		}
	}
	join := strings.Join(strs, comma)
	if len(join) != 0 && wrap {
		join = comma + join + comma
	}
	return NewCommaNumber[T](join)
}

func (r CommaNumber[T]) String() string { return string(r) }

func (r CommaNumber[T]) Parse() []T {
	if len(r) == 0 {
		return nil
	}

	parts := strings.Split(string(r), comma)
	numbers := make([]T, 0, len(parts))
	for _, part := range parts {
		v := strings.TrimSpace(part)
		if len(v) == 0 {
			continue
		}

		var err error
		var number T
		switch any(number).(type) {
		case int, int8, int16, int32, int64:
			var n int64
			n, err = strconv.ParseInt(v, 10, 64)
			number = T(n)
		case float32, float64:
			var n float64
			n, err = strconv.ParseFloat(v, 64)
			number = T(n)
		default:
			continue
		}
		if err != nil {
			return nil
		}

		numbers = append(numbers, number)
	}

	return numbers
}
