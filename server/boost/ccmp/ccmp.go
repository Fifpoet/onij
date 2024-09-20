package ccmp

import (
	"cmp"
	"hash/crc32"
	"onij/boost/collection/collext"
	"onij/boost/conv"
	"onij/boost/exp"
	"slices"
)

func ArraysEqual[T cmp.Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	slices.Sort(a)
	slices.Sort(b)
	return slices.Equal(a, b)
}

func ArraysEqualFunc[T any, K cmp.Ordered](a, b []T, orderFn func(T) K) bool {
	if len(a) != len(b) {
		return false
	}

	return ArraysEqual(collext.Pick(a, orderFn), collext.Pick(b, orderFn))
}

func StringEqual(o, a string) bool {
	if len(o) != len(a) {
		return false
	}

	if len(o) <= 1024*2 {
		return o == a
	}

	crc32O := crc32.ChecksumIEEE(conv.StringToBytes(o))
	crc32A := crc32.ChecksumIEEE(conv.StringToBytes(a))
	return crc32O == crc32A
}

func PtrEqual[T comparable](o, a *T) bool {
	if o == nil && a == nil {
		return true
	}
	if o == nil || a == nil {
		return false
	}
	switch v := (any)(o).(type) {
	case *string:
		return StringEqual(*v, *(any)(a).(*string))
	default:
		return *o == *a
	}
}

func PtrValueOrZeroEqual[T comparable](o, a *T) bool {
	switch vo := (any)(o).(type) {
	case *string:
		return StringEqual(exp.ValueOrZero(vo), exp.ValueOrZero((any)(a).(*string)))
	default:
		return exp.ValueOrZero(o) == exp.ValueOrZero(a)
	}
}
