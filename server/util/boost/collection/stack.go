package collection

import "container/list"

type Stack[T any] interface {
	Empty() bool
	Len() int
	Peek() (T, bool)
	Pop() (T, bool)
	Push(T)
	Clear()
}

func NewStack[T any]() Stack[T] { return &stack[T]{list.New()} }

type stack[T any] struct {
	l *list.List
}

func (r *stack[T]) Empty() bool { return r.Len() == 0 }

func (r *stack[T]) Len() int { return r.l.Len() }

func (r *stack[T]) Push(t T) { r.l.PushBack(t) }

func (r *stack[T]) Peek() (T, bool) {
	if r.Empty() {
		var t T
		return t, false
	}
	return r.l.Back().Value.(T), true
}

func (r *stack[T]) Pop() (T, bool) {
	var t T
	if r.Empty() {
		return t, false
	}
	e := r.l.Back()
	t = e.Value.(T)
	r.l.Remove(e)
	return t, true
}

func (r *stack[T]) Clear() {
	for {
		if _, ok := r.Pop(); !ok {
			break
		}
	}
}
