package exp

func Ptr[T any](source T) *T { return &source }

func Zero[T any]() T { var v T; return v }

func ValueOrZero[T any](t *T) T {
	if t == nil {
		return Zero[T]()
	}
	return *t
}

func PtrCopy[T any](v *T) *T {
	if v == nil {
		return nil
	}
	return Ptr(ValueOrZero(v))
}

func Values[T any](t T) []T { return []T{t} }

func PtrValues[T any](t *T, zero bool) []T {
	if t == nil {
		if zero {
			return []T{Zero[T]()}
		}
		return nil
	}
	return []T{*t}
}
