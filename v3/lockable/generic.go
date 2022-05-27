package lockable

import (
	"sync"
)

type Locker[T any] struct {
	l sync.RWMutex
	i T
}

func New[T any](i T) *Locker[T] {
	return &Locker[T]{i: i}
}

func (l *Locker[T]) Get(f func(item T)) {
	l.l.RLock()
	defer l.l.RUnlock()
	f(l.i)
}

func (l *Locker[T]) Set(f func(item *T)) {
	l.l.Lock()
	defer l.l.Unlock()
	f(&l.i)
}
