package collection

import (
	"container/heap"
	"math"
)

type PriorityQueue[T any] interface {
	Empty() bool
	Len() int
	PushDefault(T)
	Push(T, float64)
	Pop() (T, bool)
	Update(T, float64) bool
	Priority(T) (float64, bool)
	Contains(T) bool
}

type PriorityQueueOptions[T any] struct {
	hash    func(T) int64
	maximum bool
	cap     int
}

type PriorityQueueOption[T any] func(*PriorityQueueOptions[T])

type priorityQueue[T any] struct {
	heap     *binHeap[T]
	hash     func(T) int64
	cache    map[int64]*item[T]
	priority float64
}

type item[T any] struct {
	key      int64
	value    T
	priority float64
	index    int
}

type binHeap[T any] struct {
	items   []*item[T]
	compare func(a, n *item[T]) bool
}

func WithPQHash[T any](hash func(T) int64) PriorityQueueOption[T] {
	return func(o *PriorityQueueOptions[T]) { o.hash = hash }
}

func WithPQMaximum[T any]() PriorityQueueOption[T] {
	return func(o *PriorityQueueOptions[T]) { o.maximum = true }
}

func WithPQCap[T any](cap int) PriorityQueueOption[T] {
	return func(o *PriorityQueueOptions[T]) { o.cap = cap }
}

func NewPriorityQueue[T any](opts ...PriorityQueueOption[T]) PriorityQueue[T] {
	o := new(PriorityQueueOptions[T])
	for _, opt := range opts {
		opt(o)
	}

	return &priorityQueue[T]{
		&binHeap[T]{
			func() []*item[T] {
				if o.cap > 0 {
					return make([]*item[T], 0, o.cap)
				}
				return make([]*item[T], 0)
			}(),
			func(a, n *item[T]) bool {
				if o.maximum {
					return a.priority > n.priority
				}
				return a.priority < n.priority
			},
		},
		o.hash,
		func() map[int64]*item[T] {
			if o.cap > 0 {
				return make(map[int64]*item[T], o.cap)
			}
			return make(map[int64]*item[T])
		}(),
		func() float64 {
			if o.maximum {
				return math.Inf(-1)
			}
			return math.Inf(0)
		}(),
	}
}

func (r *priorityQueue[T]) Empty() bool { return r.heap.Len() == 0 }

func (r *priorityQueue[T]) Len() int { return r.heap.Len() }

func (r *priorityQueue[T]) PushDefault(t T) { r.Push(t, r.priority) }

func (r *priorityQueue[T]) Push(t T, priority float64) {
	var key int64
	if r.hash != nil {
		key = r.hash(t)
		if _, ok := r.cache[key]; ok {
			return
		}
	}

	item := &item[T]{value: t, priority: priority, key: key}
	heap.Push(r.heap, item)

	if r.hash != nil {
		r.cache[item.key] = item
	}
}

func (r *priorityQueue[T]) Pop() (T, bool) {
	if r.heap.Len() == 0 {
		var t T
		return t, false
	}

	item := heap.Pop(r.heap).(*item[T])
	if r.hash != nil {
		delete(r.cache, item.key)
	}
	return item.value, true
}

func (r *priorityQueue[T]) Update(t T, priority float64) bool {
	item, ok := r.get(t)
	if !ok {
		return false
	}

	item.priority = priority
	heap.Fix(r.heap, item.index)
	return true
}

func (r *priorityQueue[T]) Priority(t T) (float64, bool) {
	item, ok := r.get(t)
	if !ok {
		return 0, false
	}
	return item.priority, true
}

func (r *priorityQueue[T]) Contains(t T) bool {
	_, ok := r.get(t)
	return ok
}

func (r *priorityQueue[T]) get(t T) (*item[T], bool) {
	if r.hash == nil {
		return nil, false
	}

	item, ok := r.cache[r.hash((t))]
	if !ok {
		return nil, false
	}

	return item, true
}

func (r *binHeap[T]) Len() int { return len(r.items) }

func (r *binHeap[T]) Less(i, j int) bool { return r.compare(r.items[i], r.items[j]) }

func (r binHeap[T]) Swap(i, j int) {
	r.items[i], r.items[j] = r.items[j], r.items[i]
	r.items[i].index = i
	r.items[j].index = j
}

func (r *binHeap[T]) Push(t any) {
	item := t.(*item[T])
	item.index = r.Len()
	r.items = append(r.items, item)
}

func (r *binHeap[T]) Pop() any {
	if r.Len() == 0 {
		return nil
	}

	index := r.Len() - 1
	item := r.items[index]
	item.index = -1

	r.items[index] = nil
	r.items = r.items[:index]

	return item
}
