package collection

import "container/list"

type Queue[T any] interface {
	Dequeue() (T, bool)
	Empty() bool
	Enqueue(T)
	Len() int
	Peek() (T, bool)
	Clear()
}

func NewQueue[T any]() Queue[T] { return &queue[T]{list.New()} }

type queue[T any] struct {
	l *list.List
}

func (r *queue[T]) Empty() bool { return r.Len() == 0 }

func (r *queue[T]) Len() int { return r.l.Len() }

func (r *queue[T]) Enqueue(t T) { r.l.PushBack(t) }

func (r *queue[T]) Dequeue() (T, bool) {
	var t T
	if r.Empty() {
		return t, false
	}

	e := r.l.Front()
	t = e.Value.(T)

	r.l.Remove(e)

	return t, true
}

func (r *queue[T]) Peek() (T, bool) {
	if r.Empty() {
		var t T
		return t, false
	}

	return r.l.Front().Value.(T), true
}

func (r *queue[T]) Clear() {
	for {
		if _, ok := r.Dequeue(); !ok {
			break
		}
	}
}
