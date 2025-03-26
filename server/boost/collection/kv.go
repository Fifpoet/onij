package collection

import (
	"iter"
	"maps"
)

type KV[K comparable, V any] struct {
	m map[K]V
}

func NewKV[K comparable, V any]() *KV[K, V] {
	return &KV[K, V]{m: make(map[K]V)}
}

func (r *KV[K, V]) Get(k K) (v V, ok bool) {
	v, ok = r.m[k]
	return
}

func (r *KV[K, V]) Add(k K, v V) {
	r.m[k] = v
}

func (r *KV[K, V]) Remove(keys ...K) {
	for _, key := range keys {
		delete(r.m, key)
	}
}

func (r *KV[K, V]) RemoveAll(fn func(K) bool) {
	if fn == nil {
		return
	}

	var keys []K
	for k := range r.m {
		if fn(k) {
			keys = append(keys, k)
		}
	}
	for _, key := range keys {
		delete(r.m, key)
	}
}

func (r *KV[K, V]) Clear() {
	clear(r.m)
}

func (r *KV[K, V]) Reset() {
	r.m = make(map[K]V)
}

func (r *KV[K, V]) Len() int {
	return len(r.m)
}

func (r *KV[K, V]) Empty() bool {
	return r.Len() == 0
}

func (r *KV[K, V]) Contains(k K) bool {
	if _, ok := r.m[k]; ok {
		return true
	}
	return false
}

func (r *KV[K, V]) Keys() []K {
	keys := make([]K, 0, len(r.m))
	for k := range r.m {
		keys = append(keys, k)
	}
	return keys
}

func (r *KV[K, V]) Values() []V {
	values := make([]V, 0, len(r.m))
	for _, v := range r.m {
		values = append(values, v)
	}
	return values
}

func (r *KV[K, V]) Pairs() iter.Seq2[K, V] {
	return maps.All(r.m)
}
