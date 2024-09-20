package collext

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"onij/boost/exp"
)

func SelectKeys[T any, K comparable](source []T, selector func(T) (K, bool)) []K {
	set := make(map[K]struct{})
	for _, v := range source {
		if k, ok := selector(v); ok {
			set[k] = struct{}{}
		}
	}
	return maps.Keys(set)
}

func Select[T any, K any](source []T, selector func(T) (K, bool)) []K {
	if selector == nil {
		return nil
	}

	ks := make([]K, 0, len(source))
	for _, v := range source {
		if k, ok := selector(v); ok {
			ks = append(ks, k)
		}
	}
	return ks
}

func SelectOne[T any, K any](source []T, selector func(T) (K, bool)) K {
	if selector == nil {
		return exp.Zero[K]()
	}

	for _, v := range source {
		if k, ok := selector(v); ok {
			return k
		}
	}
	return exp.Zero[K]()
}

func Pick[T any, K any](source []T, selector func(T) K) []K {
	if selector == nil || source == nil {
		return nil
	}

	ks := make([]K, 0, len(source))
	for _, v := range source {
		ks = append(ks, selector(v))
	}
	return ks
}

func PickOne[T any, K any](source []T, selector func(T) K) K {
	if selector == nil {
		return exp.Zero[K]()
	}

	for _, v := range source {
		return selector(v)
	}

	return exp.Zero[K]()
}

func PickWithIndex[T any, K any](source []T, selector func(int, T) K) []K {
	if selector == nil || source == nil {
		return nil
	}

	ks := make([]K, 0, len(source))
	for i, v := range source {
		ks = append(ks, selector(i, v))
	}
	return ks
}

func PickCombine[T, K any](source []T, f func(T) []K) []K {
	if f == nil || source == nil {
		return nil
	}

	items := make([]K, 0)
	for _, v := range source {
		items = append(items, f(v)...)
	}
	return items
}

func Combine[T any](sources ...[]T) []T {
	if len(sources) == 0 {
		return nil
	}

	var count int
	for _, v := range sources {
		count += len(v)
	}

	items := make([]T, 0, count)
	for _, v := range sources {
		items = append(items, v...)
	}
	return items
}

func Distinct[T comparable](source []T, exclude ...T) []T {
	sets := make(map[T]struct{})
	exps := Sets(exclude)
	for _, v := range source {
		if _, ok := exps[v]; ok {
			continue
		}
		sets[v] = struct{}{}
	}
	return maps.Keys(sets)
}

func Paging[T any](source []T, offset, limit int) []T {
	total := len(source)
	if offset > total {
		return nil
	}
	if offset+limit > total {
		limit = total - offset
	}
	return source[offset : offset+limit]
}

func CombineDistinct[T any, K comparable](selector func(T) K, sources [][]T, excludes ...T) []T {
	if len(sources) == 0 {
		return nil
	}

	var count int
	for _, v := range sources {
		count += len(v)
	}

	exps := KeySets(excludes, selector)
	sets := make(map[K]struct{})
	items := make([]T, 0, count)
	for _, source := range sources {
		for _, v := range source {
			key := selector(v)
			if _, ok := exps[key]; ok {
				continue
			}
			if _, ok := sets[key]; ok {
				continue
			}
			sets[key] = struct{}{}
			items = append(items, v)
		}
	}
	return items
}

func Index[T any](source []T, i int) (T, bool) {
	if i < 0 || i >= len(source) {
		return exp.Zero[T](), false
	}
	return source[i], true
}

func Sum[T constraints.Integer | constraints.Float](source []T) T {
	var sum T
	for _, v := range source {
		sum += v
	}
	return sum
}

func PickSum[T any, K constraints.Integer | constraints.Float](source []T, selector func(T) K) K {
	var sum K
	for _, v := range source {
		sum += selector(v)
	}
	return sum
}
