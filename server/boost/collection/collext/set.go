package collext

import "onij/boost/collection"

func Intersect[K comparable](sources ...[]K) collection.Collection[K] {
	if len(sources) == 0 {
		return nil
	}
	if len(sources) == 1 {
		return Copy(sources[0])
	}

	var set collection.Set[K]
	for i, v := range sources {
		if i == 0 {
			set = collection.NewSet(v...)
			continue
		}
		if set = set.Intersect(collection.NewSet(v...)); set.Empty() {
			return nil
		}
	}
	return set.Elements()
}

func Union[K comparable](sources ...[]K) collection.Collection[K] {
	if len(sources) == 0 {
		return nil
	}
	if len(sources) == 1 {
		return Copy(sources[0])
	}

	var set collection.Set[K]
	for i, v := range sources {
		if i == 0 {
			set = collection.NewSet(v...)
			continue
		}
		set = set.Union(collection.NewSet(v...))
	}
	return set.Elements()
}

func Subtract[K comparable](sources ...[]K) collection.Collection[K] {
	if len(sources) == 0 || len(sources) == 1 {
		return nil
	}

	sets := Pick(sources, func(source []K) collection.Set[K] {
		return collection.NewSet(source...)
	})
	union := collection.NewSet(Union(sources...)...)
	return union.Subtract(sets...).Elements()
}

func IntersectFunc[T any, K comparable](fn func(T) K, sources ...[]T) collection.Collection[T] {
	if len(sources) == 0 {
		return nil
	}
	if len(sources) == 1 {
		return Copy(sources[0])
	}

	var set collection.Set[T]
	for i, v := range sources {
		if i == 0 {
			set = collection.NewAnySetWithKeyFunc(fn, v...)
			continue
		}
		if set = set.Intersect(collection.NewAnySetWithKeyFunc(fn, v...)); set.Empty() {
			return nil
		}
	}
	return set.Elements()

}

func UnionFunc[T any, K comparable](fn func(T) K, sources ...[]T) collection.Collection[T] {
	if len(sources) == 0 {
		return nil
	}
	if len(sources) == 1 {
		return Copy(sources[0])
	}

	var set collection.Set[T]
	for i, v := range sources {
		if i == 0 {
			set = collection.NewAnySetWithKeyFunc(fn, v...)
			continue
		}
		set = set.Union(collection.NewAnySetWithKeyFunc(fn, v...))
	}
	return set.Elements()

}

func SubtractFunc[T any, K comparable](fn func(T) K, sources ...[]T) collection.Collection[T] {
	if len(sources) == 0 || len(sources) == 1 {
		return nil
	}

	sets := Pick(sources, func(source []T) collection.Set[T] {
		return collection.NewAnySetWithKeyFunc(fn, source...)
	})
	union := collection.NewAnySetWithKeyFunc(fn, UnionFunc(fn, sources...)...)
	return union.Subtract(sets...).Elements()
}
