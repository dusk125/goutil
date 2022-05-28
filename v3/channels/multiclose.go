package channels

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// ChanWriter gives permission to a callee to write to a channel.
// 	Implicitly gives permission to the callee to be notified when the channel is closed.
type ChanWriter[T any] interface {
	Write(v T) bool
}

// ChanReader gives permission to a callee to read from a channel.
type ChanReader[T any] interface {
	Read() <-chan T
}

// ChanCloser gives permission to a callee to close a channel.
type ChanCloser[T any] interface {
	Close()
}

// Gives both read and close permissions.
type ChanReadCloser[T any] interface {
	ChanReader[T]
	ChanCloser[T]
}

// Gives both write and close permissions.
type ChanWriteCloser[T any] interface {
	ChanWriter[T]
	ChanCloser[T]
}

// Gives both read and write permissions.
type ChanReadWriter[T any] interface {
	ChanReader[T]
	ChanWriter[T]
}

// Gives all read, write, and close permissions.
type ChanReadWriteCloser[T any] interface {
	ChanReadWriter[T]
	ChanCloser[T]
}

// The Multiclose Chan is a channel that can only ever be closed once, additional calls to Close() will be a no-op.
type MulticloseChan[T any] struct {
	ch     chan T
	closed chan struct{}
	open   int32
	o, c   sync.Once
}

// Allocates and makes a channel.
func MulticloseMake[T any](buffered int) *MulticloseChan[T] {
	c := MulticloseChan[T]{}
	c.Make(buffered)
	return &c
}

// Make allocates the underlying members of the channel for operation.
func (m *MulticloseChan[T]) Make(buffered int) {
	m.o.Do(func() {
		m.closed = make(chan struct{})
		m.ch = make(chan T, buffered)
		atomic.StoreInt32(&m.open, 1)
	})
}

// Tells the caller whether the channel is open or not.
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
		close(m.closed)
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

// Returns a context.Context that will fire when this channel is closed.
func (m *MulticloseChan[T]) Context() context.Context {
	return ctxt(m.closed)
}

type ctxt <-chan struct{}

func (c ctxt) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c ctxt) Done() <-chan struct{} {
	return c
}

func (c ctxt) Err() error {
	return nil
}

func (c ctxt) Value(key any) any {
	return nil
}
