package channels

import "sync"

type ChanWriter[T any] interface {
	Write(v T)
}

type ChanReader[T any] interface {
	Read() (v T, ok bool)
}

type ChanCloser[T any] interface {
	Close()
}

type ChanReadCloser[T any] interface {
	ChanReader[T]
	ChanCloser[T]
}

type ChanWriteCloser[T any] interface {
	ChanWriter[T]
	ChanCloser[T]
}

type ChanReadWriter[T any] interface {
	ChanReader[T]
	ChanWriter[T]
}

type ChanReadWriterCloser[T any] interface {
	ChanReader[T]
	ChanWriter[T]
	ChanCloser[T]
}

type MulticloseChan[T any] struct {
	ch chan T
	o  sync.Once
}

func MulticloseMake[T any](buffered int) *MulticloseChan[T] {
	return &MulticloseChan[T]{
		ch: make(chan T, buffered),
	}
}

func (m *MulticloseChan[T]) Close() {
	m.o.Do(func() { close(m.ch) })
}

func (m *MulticloseChan[T]) Write(v T) {
	m.ch <- v
}

func (m *MulticloseChan[T]) Read() (v T, ok bool) {
	v, ok = <-m.ch
	return
}
