package channels

import (
	"sync"
	"sync/atomic"
)

// ChanWriter gives permission to a callee to write to a channel
type ChanWriter[T any] interface {
	Write(v T) bool
}

// ChanReader gives permission to a callee to read from a channel
type ChanReader[T any] interface {
	Read() <-chan T
}

// ChanCloser gives permission to a callee to close a channel
type ChanCloser[T any] interface {
	Close()
}

// Gives both read and close permissions
type ChanReadCloser[T any] interface {
	ChanReader[T]
	ChanCloser[T]
}

// Gives both write and close permissions
type ChanWriteCloser[T any] interface {
	ChanWriter[T]
	ChanCloser[T]
}

// Gives both read and write permissions
type ChanReadWriter[T any] interface {
	ChanReader[T]
	ChanWriter[T]
}

// Gives all read, write, and close permissions
type ChanReadWriteCloser[T any] interface {
	ChanReader[T]
	ChanWriter[T]
	ChanCloser[T]
}

// The Multiclose Chan is a channel that can only ever be closed once, additional calls to Close() will be a no-op.
type MulticloseChan[T any] struct {
	ch   chan T
	open int32
	o, c sync.Once
}

func MulticloseMake[T any](buffered int) *MulticloseChan[T] {
	c := MulticloseChan[T]{}
	c.Make(buffered)
	return &c
}

func (m *MulticloseChan[T]) Make(buffered int) {
	m.o.Do(func() {
		m.ch = make(chan T, buffered)
		atomic.StoreInt32(&m.open, 1)
	})
}

func (m *MulticloseChan[T]) Open() bool {
	return atomic.LoadInt32(&m.open) == 1
}

// Close the MulticloseChan; multiple invocations (across goroutines) are allowed.
// 	If this is the first time this has been called on this channel, close it;
// 	if this is a subsequent call to close, then no-op.
func (m *MulticloseChan[T]) Close() {
	m.c.Do(func() {
		atomic.StoreInt32(&m.open, 0)
		close(m.ch)
	})
}

// Write to the underlying channel, if the underlying channel is closed, no-op.
func (m *MulticloseChan[T]) Write(v T) bool {
	if m.Open() {
		m.ch <- v
		return true
	}
	return false
}

// Read from the underlying channel
func (m *MulticloseChan[T]) Read() <-chan T {
	return m.ch
}
