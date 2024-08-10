package collext

import "code.chenji.com/pkg/boost/exp"

func Map[T any, K comparable](slice []T, kFn func(T) K) map[K]T {
	var mapping = make(map[K]T)
	for _, v := range slice {
		mapping[kFn(v)] = v
	}
	return mapping
}

func MapKV[T any, K comparable, R any](slice []T, kFn func(T) K, vFn func(T) R) map[K]R {
	var mapping = make(map[K]R)
	for _, v := range slice {
		mapping[kFn(v)] = vFn(v)
	}
	return mapping
}

func Group[T any, K comparable](slice []T, kFn func(T) K) map[K][]T {
	var mapping = make(map[K][]T)
	for _, t := range slice {
		key := kFn(t)
		var values []T
		if v, ok := mapping[key]; ok {
			values = v
		}
		values = append(values, t)
		mapping[key] = values
	}
	return mapping
}

func GroupMap[T any, K comparable, R any](slice []T, kFn func(T) K, vFn func(T) R) map[K][]R {
	var mapping = make(map[K][]R)
	for _, t := range slice {
		key := kFn(t)
		var values []R
		if v, ok := mapping[key]; ok {
			values = v
		}
		values = append(values, vFn(t))
		mapping[key] = values
	}
	return mapping
}

func SelectMap[T any, K comparable](source []T, selector func(T) (K, bool)) map[K]T {
	mapping := make(map[K]T)
	for _, v := range source {
		if k, ok := selector(v); ok {
			mapping[k] = v
		}
	}
	return mapping
}

func MapTakeOne[K comparable, V any, T ~map[K]V](source T) (K, V) {
	for k, v := range source {
		return k, v
	}
	return exp.Zero[K](), exp.Zero[V]()
}

func KeySets[T any, K comparable](source []T, selector func(T) K) map[K]struct{} {
	sets := make(map[K]struct{})
	for _, v := range source {
		sets[selector(v)] = struct{}{}
	}
	return sets
}

func Sets[T comparable](source []T) map[T]struct{} {
	sets := make(map[T]struct{})
	for _, v := range source {
		sets[v] = struct{}{}
	}
	return sets
}
