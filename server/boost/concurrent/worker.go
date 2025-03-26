package concurrent

import (
	"sync"
)

type worker[T, R any] struct {
	taskQueue chan *task[T, R]
}

type workerPool[T, R any] struct {
	maxWorkers  int
	waitWorker  bool
	workerCount int
	workers     []*worker[T, R]
	stopc       <-chan struct{}
	ready       chan *worker[T, R]

	mu sync.Mutex
}

func newWorkerPool[T, R any](maxWorkers int, waitWorker bool, stopc <-chan struct{}) *workerPool[T, R] {
	w := &workerPool[T, R]{
		maxWorkers: maxWorkers,
		waitWorker: waitWorker,
		ready:      make(chan *worker[T, R], maxWorkers),
		stopc:      stopc,
	}
	go func() {
		<-stopc
		w.mu.Lock()
		for _, v := range w.workers {
			v.taskQueue <- nil
		}
		w.workers = nil
		w.mu.Unlock()
	}()
	return w
}

func (r *workerPool[T, R]) Serve(task *task[T, R]) bool {
	w := r.worker()
	if w == nil {
		return false
	}

	w.taskQueue <- task
	return true
}

func (r *workerPool[T, R]) worker() *worker[T, R] {
	for {
		select {
		case <-r.stopc:
			return nil
		case worker := <-r.ready:
			return worker
		default:
		}

		var w *worker[T, R]
		r.mu.Lock()
		if r.workerCount < r.maxWorkers {
			r.workerCount++
			w = &worker[T, R]{taskQueue: make(chan *task[T, R])}
			r.workers = append(r.workers, w)
			r.runWorker(w)
		}
		r.mu.Unlock()

		if w != nil {
			return w
		}
		if !r.waitWorker {
			return nil
		}

		select {
		case <-r.stopc:
			return nil
		case worker := <-r.ready:
			return worker
		}
	}
}

func (r *workerPool[T, R]) runWorker(worker *worker[T, R]) {
	go func() {
		for v := range worker.taskQueue {
			select {
			case <-r.stopc:
				return
			default:
			}
			if v == nil {
				break
			}

			select {
			case <-v.run():
				r.ready <- worker
			case <-r.stopc:
				v.discard()
			}
		}
	}()
}
