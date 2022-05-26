package channels

import "sync"

// ChanWriter gives permission to a callee to write to a channel
type ChanWriter[T any] interface {
	Write(v T)
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
type ChanReadWriterCloser[T any] interface {
	ChanReader[T]
	ChanWriter[T]
	ChanCloser[T]
}

// The Multiclose Chan is a channel that can only ever be closed once, additional calls to Close() will be a no-op.
type MulticloseChan[T any] struct {
	ch chan T
	o  sync.Once
}

func MulticloseMake[T any](buffered int) *MulticloseChan[T] {
	return &MulticloseChan[T]{
		ch: make(chan T, buffered),
	}
}

// Close the MulticloseChan; multiple invocations (across goroutines) are allowed.
// 	If this is the first time this has been called on this channel, close it;
// 	if this is a subsequent call to close, then no-op.
func (m *MulticloseChan[T]) Close() {
	m.o.Do(func() { close(m.ch) })
}

// Write to the underlying channel
func (m *MulticloseChan[T]) Write(v T) {
	m.ch <- v
}

// Read from the underlying channel
func (m *MulticloseChan[T]) Read() <-chan T {
	return m.ch
}
