package numbers

import "golang.org/x/exp/constraints"

const (
	comma = ","
	dot   = "."
)

type Number interface {
	constraints.Float | constraints.Signed
}

func Which[T Number](number T) (signed, float bool) {
	switch any(number).(type) {
	case int, int8, int16, int32, int64:
		return true, false
	case float32, float64:
		return false, true
	}
	return false, false
}
