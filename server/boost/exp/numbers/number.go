package numbers

import (
	"code.chenji.com/pkg/boost/ccmp"
	"code.chenji.com/pkg/boost/collection/collext"
	"golang.org/x/exp/constraints"
)

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

type Numbers[T Number] []T

func (r Numbers[T]) Equals(t Numbers[T]) bool { return ccmp.ArraysEqual(r, t) }

func (r Numbers[T]) Copy() Numbers[T] { return collext.Copy(r) }

func (r Numbers[T]) Base() []T { return r }
