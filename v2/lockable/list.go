package lockable

import (
	"encoding/json"
	"sort"
	"sync"
)

type List[T any] struct {
	l sync.RWMutex
	e []T
}

func NewList[T any]() *List[T] {
	return &List[T]{e: make([]T, 0)}
}

func (l *List[T]) Get(i int) T {
	l.l.RLock()
	defer l.l.RUnlock()
	return l.e[i]
}

func (l *List[T]) Append(items ...T) {
	l.l.Lock()
	defer l.l.Unlock()
	l.e = append(l.e, items...)
}

func (l *List[T]) Prepend(items ...T) {
	l.l.Lock()
	defer l.l.Unlock()
	l.e = append(items, l.e...)
}

func (l *List[T]) Len() int {
	l.l.RLock()
	defer l.l.RUnlock()
	return len(l.e)
}

func (l *List[T]) Foreach(f func(i int, v T) bool) {
	l.l.RLock()
	defer l.l.RUnlock()
	for i, v := range l.e {
		f(i, v)
	}
}

func (l *List[T]) Slice(s, e int) []T {
	l.l.RLock()
	defer l.l.RUnlock()
	return l.e[s:e]
}

func (l *List[T]) Set(items []T) {
	l.l.Lock()
	defer l.l.Unlock()
	l.e = items
}

func (l *List[T]) MarshalJSON() ([]byte, error) {
	l.l.RLock()
	defer l.l.RUnlock()
	return json.Marshal(l.e)
}

func (l *List[T]) UnMarshalJSON(b []byte) error {
	l.l.Lock()
	defer l.l.Unlock()
	return json.Unmarshal(b, &l.e)
}

func (l *List[T]) SortSlice(lessF func(i, j T) bool) {
	l.l.Lock()
	defer l.l.Unlock()
	sort.Slice(l.e, func(i, j int) bool { return lessF(l.e[i], l.e[j]) })
}
