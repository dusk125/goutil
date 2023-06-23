package lockable

import (
	"sync"
)

type Locker[T any] struct {
	sync.RWMutex
	item T
}

func New[T any](i T) *Locker[T] {
	return &Locker[T]{item: i}
}

func (l *Locker[T]) Get() T {
	l.RLock()
	defer l.RUnlock()
	return l.item
}

func (l *Locker[T]) Set(item T) {
	l.Lock()
	defer l.Unlock()
	l.item = item
}

// Safe allows you to call multiple unsafe methods with only a single lock
// write is true if you want to obtain a write lock
func (l *Locker[T]) Safe(write bool, f func()) {
	var locker, unlocker func()
	if write {
		locker = l.Lock
		unlocker = l.Unlock
	} else {
		locker = l.RLock
		unlocker = l.RUnlock
	}
	locker()
	defer unlocker()
	f()
}
