package indexer

import (
	"sync"

	"github.com/pkg/errors"
)

// errors
var (
	ErrEmptyQueue = errors.New("pop from empty queue")
)

type queue[T any] struct {
	arr []T
	mx  *sync.RWMutex
}

func newQueue[T any]() queue[T] {
	return queue[T]{
		arr: make([]T, 0),
		mx:  new(sync.RWMutex),
	}
}

// Push -
func (q *queue[T]) Push(el T) {
	q.mx.Lock()
	q.arr = append(q.arr, el)
	q.mx.Unlock()
}

// Pop -
func (q *queue[T]) Pop() (T, error) {
	q.mx.Lock()
	defer q.mx.Unlock()

	if len(q.arr) == 0 {
		var el T
		return el, ErrEmptyQueue
	}
	last := q.arr[0]
	q.arr = q.arr[1:]
	return last, nil
}

// First -
func (q *queue[T]) First() (T, error) {
	q.mx.RLock()
	defer q.mx.RUnlock()

	if len(q.arr) == 0 {
		var el T
		return el, ErrEmptyQueue
	}

	return q.arr[0], nil
}

// Size -
func (q *queue[T]) Size() int {
	q.mx.RLock()
	defer q.mx.RUnlock()

	return len(q.arr)
}

// IsEmpty -
func (q *queue[T]) IsEmpty() bool {
	q.mx.RLock()
	defer q.mx.RUnlock()
	return len(q.arr) == 0
}
