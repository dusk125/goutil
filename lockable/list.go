package lockable

import (
	"encoding/json"

	"golang.org/x/exp/slices"
)

// A lockable list that provides some list-like methods for convienence
type List[T any] struct {
	Locker[[]T]
}

func NewList[T any](length int, capacity ...int) *List[T] {
	l := &List[T]{}
	l.UnsafeMake(length, capacity...)
	return l
}

func (l *List[T]) Make(length int, capacity ...int) {
	l.Lock()
	defer l.Unlock()
	l.UnsafeMake(length, capacity...)
}

func (l *List[T]) UnsafeMake(length int, capacity ...int) {
	var cap int
	if len(capacity) > 0 {
		cap = capacity[0]
	}
	l.item = make([]T, length, cap)
}

func (l *List[T]) Nil() bool {
	l.RLock()
	defer l.RUnlock()
	return l.UnsafeNil()
}

func (l *List[T]) UnsafeNil() bool {
	return l.item == nil
}

// Get returns a copy of the internal list
func (l *List[T]) Get() []T {
	l.RLock()
	defer l.RUnlock()
	return l.UnsafeGet()
}

func (l *List[T]) UnsafeGet() []T {
	o := make([]T, len(l.item))
	copy(o, l.item)
	return o
}

// At returns the item at the given index
func (l *List[T]) At(index int) T {
	l.RLock()
	defer l.RUnlock()
	return l.UnsafeAt(index)
}

func (l *List[T]) UnsafeAt(index int) T {
	return l.item[index]
}

// Append adds the given items to the back of the list
func (l *List[T]) Append(items ...T) {
	l.Lock()
	defer l.Unlock()
	l.UnsafeAppend(items...)
}

func (l *List[T]) UnsafeAppend(items ...T) {
	l.item = append(l.item, items...)
}

// Prepend add the given items to the front of the list
func (l *List[T]) Prepend(items ...T) {
	l.Lock()
	defer l.Unlock()
	l.UnsafePrepend(items...)
}

func (l *List[T]) UnsafePrepend(items ...T) {
	l.item = append(items, l.item...)
}

func (l *List[T]) Len() int {
	l.RLock()
	defer l.RUnlock()
	return l.UnsafeLen()
}

func (l *List[T]) UnsafeLen() int {
	return len(l.item)
}

func (l *List[T]) Cap() int {
	l.RLock()
	defer l.RUnlock()
	return l.UnsafeCap()
}

func (l *List[T]) UnsafeCap() int {
	return cap(l.item)
}

// Find returns the first index where the given function returns true.
// Will return length of the list if item was not found
func (l *List[T]) Find(f func(index int, val T) bool) int {
	l.RLock()
	defer l.RUnlock()
	return l.UnsafeFind(f)
}

func (l *List[T]) UnsafeFind(f func(index int, val T) bool) int {
	for i, v := range l.item {
		if f(i, v) {
			return i
		}
	}
	return len(l.item)
}

// Delete removes the given index from the list
func (l *List[T]) Delete(index int) {
	l.Lock()
	defer l.Unlock()
	l.UnsafeDelete(index)
}

func (l *List[T]) UnsafeDelete(index int) {
	l.item = append(l.item[:index], l.item[index+1:]...)
}

// FindAndDelete deletes the first item where the given function returns true.
func (l *List[T]) FindAndDelete(f func(i int, v T) bool) (deleted bool) {
	l.Lock()
	defer l.Unlock()
	return l.UnsafeFindAndDelete(f)
}

func (l *List[T]) UnsafeFindAndDelete(f func(i int, v T) bool) (deleted bool) {
	i := l.UnsafeFind(f)
	if deleted = i < len(l.item); deleted {
		l.UnsafeDelete(i)
	}
	return
}

// Foreach runs the given function for every item in the list
func (l *List[T]) Foreach(f func(index int, val T) (shouldBreak bool)) {
	l.RLock()
	defer l.RUnlock()
	l.UnsafeForeach(f)
}

func (l *List[T]) UnsafeForeach(f func(index int, val T) (shouldBreak bool)) {
	for i, v := range l.item {
		if f(i, v) {
			break
		}
	}
}

// Slice returns a copy of the list sliced between start and end
func (l *List[T]) Slice(start, end int) []T {
	l.RLock()
	defer l.RUnlock()
	return l.UnsafeSlice(start, end)
}

func (l *List[T]) UnsafeSlice(start, end int) []T {
	o := make([]T, end-start)
	copy(o, l.item[start:end])
	return o
}

func (l *List[T]) SliceInPlace(start, end int) {
	l.Lock()
	defer l.Unlock()
	l.UnsafeSliceInPlace(start, end)
}

func (l *List[T]) UnsafeSliceInPlace(start, end int) {
	l.item = l.item[start:end]
}

// Set replaces the internal list with the given list
func (l *List[T]) Set(items []T) {
	l.Lock()
	defer l.Unlock()
	l.UnsafeSet(items)
}

func (l *List[T]) UnsafeSet(items []T) {
	l.item = items
}

func (l *List[T]) MarshalJSON() ([]byte, error) {
	l.RLock()
	defer l.RUnlock()
	return l.UnsafeMarshalJSON()
}

func (l *List[T]) UnsafeMarshalJSON() ([]byte, error) {
	return json.Marshal(l.item)
}

func (l *List[T]) UnmarshalJSON(b []byte) error {
	l.Lock()
	defer l.Unlock()
	return l.UnsafeUnmarshalJSON(b)
}

func (l *List[T]) UnsafeUnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &l.item)
}

// Sort sorts the list in place
func (l *List[T]) Sort(less func(i, j T) bool) {
	l.Lock()
	defer l.Unlock()
	l.UnsafeSort(less)
}

func (l *List[T]) UnsafeSort(less func(i, j T) bool) {
	slices.SortFunc[T](l.item, less)
}

func (l *List[T]) Read(f func()) {
	l.RLock()
	defer l.RUnlock()
	f()
}

func (l *List[T]) Write(f func()) {
	l.Lock()
	defer l.Unlock()
	f()
}
