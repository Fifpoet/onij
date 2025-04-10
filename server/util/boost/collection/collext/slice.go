package collext

import (
	"math"
	"onij/util/boost/collection"
	"onij/util/boost/exp"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

// Deprecated: use DistinctSelect instead
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
		ks := f(v)
		if len(ks) == 0 {
			continue
		}
		items = append(items, ks...)
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
		if len(v) == 0 {
			continue
		}
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
	return collection.New(source).Paging(offset, limit)
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
	return collection.New(source).Index(i)
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

func Batch[T any](items []T, size int) [][]T {
	if len(items) == 0 {
		return nil
	}
	if size <= 0 {
		return [][]T{items}
	}
	return BatchPick(items, size, func(t T) T { return t })
}

func BatchPick[T, K any](items []T, size int, selector func(T) K) [][]K {
	if len(items) == 0 {
		return nil
	}
	if size <= 0 {
		return [][]K{Pick(items, selector)}
	}

	total := len(items)
	batches := make([][]K, 0, int(math.Ceil(float64(total)/float64(size))))
	for start := 0; start < total; start += size {
		end := min(start+size, total)
		batches = append(batches, Pick(items[start:end], selector))
	}
	return batches
}

func DistinctSelect[T any, K comparable](source []T, selector func(T) (K, bool)) []K {
	if selector == nil {
		return nil
	}

	ks := make([]K, 0, len(source))
	sets := make(map[K]struct{}, len(source))
	for _, v := range source {
		k, ok := selector(v)
		if !ok {
			continue
		}
		if _, ok = sets[k]; ok {
			continue
		}
		sets[k] = struct{}{}
		ks = append(ks, k)
	}
	return ks
}

func DistinctPick[T any, K comparable](source []T, selector func(T) K) []K {
	if selector == nil || source == nil {
		return nil
	}

	ks := make([]K, 0, len(source))
	sets := make(map[K]struct{}, len(source))
	for _, v := range source {
		k := selector(v)
		if _, ok := sets[k]; ok {
			continue
		}
		sets[k] = struct{}{}
		ks = append(ks, k)
	}
	return ks
}

func Copy[T any](source []T) []T { return collection.New(source).Copy() }

func Max[T constraints.Integer | constraints.Float](source []T) T {
	var max T
	for _, v := range source {
		if v > max {
			max = v
		}
	}
	return max
}

func Min[T constraints.Integer | constraints.Float](source []T) T {
	var min T
	for _, v := range source {
		if v < min {
			min = v
		}
	}
	return min
}

func PickMax[T any, K constraints.Integer | constraints.Float](source []T, selector func(T) K) K {
	var max K
	for _, s := range source {
		if v := selector(s); v > max {
			max = v
		}
	}
	return max
}

func PickMin[T any, K constraints.Integer | constraints.Float](source []T, selector func(T) K) K {
	var min K
	for _, s := range source {
		if v := selector(s); v < min {
			min = v
		}
	}
	return min
}
