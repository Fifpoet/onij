package collection

import (
	"bytes"
	"encoding/gob"
	"github.com/cespare/xxhash/v2"
)

type Set[T any] interface {
	Add(...T)
	Remove(...T)
	RemoveAll(func(T) bool)
	Contains(T) bool
	Clear()
	Reset()
	Len() int
	Empty() bool
	Elements() Collection[T]
	Union(...Set[T]) Set[T]
	Intersect(...Set[T]) Set[T]
	Subtract(...Set[T]) Set[T]
}

func NewSet[K comparable](elements ...K) Set[K] {
	return newKeySet(elements...)
}

func NewAnySet[T any](elements ...T) Set[T] {
	return NewAnySetWithKeyFunc(hash, elements...)
}

func NewAnySetWithKeyFunc[T any, K comparable](fn func(T) K, elements ...T) Set[T] {
	return newMapSet(fn, elements...)
}

type keySet[K comparable] struct {
	*set[struct{}, K]
}

func newKeySet[K comparable](elements ...K) *keySet[K] {
	set := &keySet[K]{newSet[struct{}, K]()}
	set.Add(elements...)
	return set
}

func (r *keySet[K]) Add(elements ...K) {
	for _, v := range elements {
		r.add(v, struct{}{})
	}
}

func (r *keySet[K]) Remove(elements ...K) { r.remove(elements...) }

func (r *keySet[K]) RemoveAll(fn func(K) bool) { r.removeAll(fn) }

func (r *keySet[K]) Contains(k K) bool { return r.contains(k) }

func (r *keySet[K]) Clear() { r.clear() }

func (r *keySet[K]) Reset() { r.reset() }

func (r *keySet[K]) Len() int { return r.len() }

func (r *keySet[K]) Empty() bool { return r.empty() }

func (r *keySet[K]) Elements() Collection[K] { return r.keys() }

func (r *keySet[K]) Union(sets ...Set[K]) Set[K] {
	return &keySet[K]{r.unionSets(r.bases()...)}
}

func (r *keySet[K]) Intersect(sets ...Set[K]) Set[K] {
	return &keySet[K]{r.intersectSets(r.bases()...)}
}

func (r *keySet[K]) Subtract(sets ...Set[K]) Set[K] {
	return &keySet[K]{r.subtractSets(r.bases()...)}
}

func (r *keySet[K]) bases(sets ...Set[K]) []*set[struct{}, K] {
	bases := make([]*set[struct{}, K], 0, len(sets))
	for _, v := range sets {
		if m, ok := v.(*keySet[K]); ok {
			bases = append(bases, m.set)
			continue
		}
		bases = append(bases, newKeySet(v.Elements()...).set)
	}
	return bases
}

type mapSet[T any, K comparable] struct {
	*set[T, K]

	hash func(T) K
}

func newMapSet[T any, K comparable](fn func(T) K, elements ...T) *mapSet[T, K] {
	set := &mapSet[T, K]{hash: fn, set: newSet[T, K]()}
	set.Add(elements...)
	return set
}

func (r *mapSet[T, K]) Add(elements ...T) {
	for _, v := range elements {
		r.add(r.hash(v), v)
	}
}

func (r *mapSet[T, K]) Remove(elements ...T) {
	for _, v := range elements {
		r.remove(r.hash(v))
	}
}

func (r *mapSet[T, K]) RemoveAll(fn func(T) bool) {
	r.removeAll(func(k K) bool {
		v, _ := r.get(k)
		return fn(v)
	})
}

func (r *mapSet[T, K]) Contains(t T) bool { return r.contains(r.hash(t)) }

func (r *mapSet[T, K]) Clear() { r.clear() }

func (r *mapSet[T, K]) Reset() { r.reset() }

func (r *mapSet[T, K]) Len() int { return r.len() }

func (r *mapSet[T, K]) Empty() bool { return r.empty() }

func (r *mapSet[T, K]) Elements() Collection[T] { return r.values() }

func (r *mapSet[T, K]) Union(sets ...Set[T]) Set[T] {
	return &mapSet[T, K]{hash: r.hash, set: r.unionSets(r.bases()...)}
}

func (r *mapSet[T, K]) Intersect(sets ...Set[T]) Set[T] {
	return &mapSet[T, K]{hash: r.hash, set: r.intersectSets(r.bases()...)}
}

func (r *mapSet[T, K]) Subtract(sets ...Set[T]) Set[T] {
	return &mapSet[T, K]{hash: r.hash, set: r.subtractSets(r.bases()...)}
}

func (r *mapSet[T, K]) bases(sets ...Set[T]) []*set[T, K] {
	bases := make([]*set[T, K], 0, len(sets))
	for _, v := range sets {
		if m, ok := v.(*mapSet[T, K]); ok {
			bases = append(bases, m.set)
			continue
		}
		bases = append(bases, newMapSet(r.hash, v.Elements()...).set)
	}
	return bases
}

type set[T any, K comparable] struct {
	kv *KV[K, T]
}

func newSet[T any, K comparable]() *set[T, K] {
	return &set[T, K]{NewKV[K, T]()}
}

func (r *set[T, K]) get(k K) (T, bool) { return r.kv.Get(k) }

func (r *set[T, K]) add(k K, v T) { r.kv.Add(k, v) }

func (r *set[T, K]) remove(ks ...K) { r.kv.Remove(ks...) }

func (r *set[T, K]) removeAll(fn func(K) bool) { r.kv.RemoveAll(fn) }

func (r *set[T, K]) contains(k K) bool { return r.kv.Contains(k) }

func (r *set[T, K]) clear() { r.kv.Clear() }

func (r *set[T, K]) reset() { r.kv.Reset() }

func (r *set[T, K]) len() int { return r.kv.Len() }

func (r *set[T, K]) empty() bool { return r.kv.Empty() }

func (r *set[T, K]) keys() []K { return r.kv.Keys() }

func (r *set[T, K]) values() []T { return r.kv.Values() }

func (r *set[T, K]) unionSets(sets ...*set[T, K]) *set[T, K] {
	newSet := newSet[T, K]()
	for _, set := range append(append(make([]*set[T, K], 0, len(sets)+1), r), sets...) {
		for k, v := range set.kv.Pairs() {
			newSet.add(k, v)
		}
	}
	return newSet
}

func (r *set[T, K]) intersectSets(sets ...*set[T, K]) *set[T, K] {
	if len(sets) == 0 {
		return newSet[T, K]()
	}

	var newSet *set[T, K]
	for i, set := range sets {
		if i == 0 {
			newSet = r.intersect(set)
			continue
		}
		if newSet.kv.Empty() {
			break
		}
		newSet = newSet.intersect(set)
	}
	return newSet
}

func (r *set[T, K]) subtractSets(sets ...*set[T, K]) *set[T, K] {
	newSet := newSet[T, K]()
	intersections := r.intersectSets(sets...)
	for k, v := range r.kv.Pairs() {
		if intersections.kv.Contains(k) {
			continue
		}
		newSet.kv.Add(k, v)
	}
	return newSet
}

func (r *set[T, K]) intersect(set *set[T, K]) *set[T, K] {
	newSet := newSet[T, K]()
	for k, v := range set.kv.Pairs() {
		if r.kv.Contains(k) {
			newSet.kv.Add(k, v)
		}
	}
	return newSet
}

func hash[T any](t T) uint64 {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(t)
	return xxhash.Sum64(buf.Bytes())
}
