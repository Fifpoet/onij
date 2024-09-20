package collection

import (
	"onij/boost/collection/collext"
	"onij/boost/exp"
	"slices"
)

type Collection[T any] []T

type VisitContext[T any] struct {
	Root         T
	Parent       T
	Current      T
	CurrentIndex int
	Depth        int
}

type DualVisitContext[T any] struct {
	Current   T
	CurrentIn bool
	Target    T
	TargetIn  bool
	Index     int
}

type CollectionVisit[T any] func(*VisitContext[T]) (stop bool, skip bool, err error)
type DualCollectionVisit[T any] func(*DualVisitContext[T]) (stop bool, skip bool, err error)
type SubCollection[T any] func(T) []T

func New[T any](source []T) Collection[T] {
	return Collection[T](source)
}

func (r Collection[T]) Find(f func(T) bool) T {
	for _, v := range r {
		if f(v) {
			return v
		}
	}
	var z T
	return z
}

func (r Collection[T]) Any(f func(T) bool) bool {
	for _, v := range r {
		if f(v) {
			return true
		}
	}
	return false
}

func (r Collection[T]) All(f func(T) bool) bool {
	for _, v := range r {
		if !f(v) {
			return false
		}
	}
	return true
}

func (r Collection[T]) Where(f func(T) bool) Collection[T] {
	selected := make(Collection[T], 0, len(r))
	for _, v := range r {
		if f(v) {
			selected = append(selected, v)
		}
	}
	return selected
}

func (r Collection[T]) Paging(offset, limit int) Collection[T] {
	return collext.Paging(r, offset, limit)
}

func (r Collection[T]) Each(f func(T)) {
	for _, v := range r {
		f(v)
	}
}

func (r Collection[T]) EachWithIndex(f func(int, T)) {
	for i, v := range r {
		f(i, v)
	}
}

func (r Collection[T]) Flatten(f func(T) []T) []T {
	if f == nil {
		return r
	}

	items := make([]T, 0)
	for _, v := range r {
		items = append(items, v)
		items = append(items, Collection[T](f(v)).Flatten(f)...)
	}
	return items
}

func (r Collection[T]) First() (T, bool) {
	if len(r) == 0 {
		return exp.Zero[T](), false
	}
	return r[0], true
}

func (r Collection[T]) Sum(f func(T) int) int {
	var sum int
	for _, v := range r {
		sum += f(v)
	}
	return sum
}

func (r Collection[T]) Count() int { return len(r) }

func (r Collection[T]) Travel(sub SubCollection[T], visit CollectionVisit[T]) error {
	if visit == nil {
		return nil
	}

	var t T
	return travel(t, t, 0, sub, visit, r...)
}

func (r Collection[T]) DualTravel(t Collection[T], sub SubCollection[T], visit DualCollectionVisit[T]) error {
	if visit == nil {
		return nil
	}
	return dualTravel(r, t, sub, visit)
}

func (r Collection[T]) Sort(cmp func(a, b T) int) Collection[T] {
	if cmp == nil {
		return r
	}
	slices.SortFunc(r, cmp)
	return r
}

func travel[T any](parent, root T, depth int, sub SubCollection[T], visit CollectionVisit[T], items ...T) error {
	for i, item := range items {
		ctx := &VisitContext[T]{
			Root:         root,
			Parent:       parent,
			Current:      item,
			CurrentIndex: i,
			Depth:        depth,
		}
		stop, skip, err := visit(ctx)
		if err != nil {
			return err
		}
		if stop {
			break
		}
		if skip {
			continue
		}
		if sub == nil {
			continue
		}
		r := root
		if depth == 0 {
			r = item
		}
		if err := travel(item, r, depth+1, sub, visit, sub(item)...); err != nil {
			return err
		}
	}
	return nil
}

func dualTravel[T any](itemsC, itemsT []T, sub SubCollection[T], visit DualCollectionVisit[T]) error {
	for i := 0; i < max(len(itemsC), len(itemsT)); i++ {
		c, cok := collext.Index(itemsC, i)
		t, tok := collext.Index(itemsT, i)
		ctx := &DualVisitContext[T]{
			Current:   c,
			CurrentIn: cok,
			Target:    t,
			TargetIn:  tok,
			Index:     i,
		}
		stop, skip, err := visit(ctx)
		if err != nil {
			return err
		}
		if stop {
			break
		}
		if skip {
			continue
		}
		if sub == nil {
			continue
		}
		var subsC, subsT []T
		if cok {
			subsC = sub(c)
		}
		if tok {
			subsT = sub(t)
		}
		if err := dualTravel(subsC, subsT, sub, visit); err != nil {
			return err
		}
	}
	return nil
}
