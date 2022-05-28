package channels

import (
	"sync"
)

// The Merger allows multiple writers to a single channel without worrying about writing to a closed channel
type Merger[T any] struct {
	aggr *MulticloseChan[T]
	wg   sync.WaitGroup
}

func NewMerger[T any]() (m *Merger[T]) {
	m = &Merger[T]{
		aggr: MulticloseMake[T](1),
	}
	return
}

// Provides an channel reader to consume the aggregate stream.
func (m *Merger[T]) Aggr() ChanReader[T] {
	return m.aggr
}

// Adds a channel to the stream and returns it; will call 'onstop' when the Merger is closed.
func (m *Merger[T]) Add(onstop func()) ChanWriter[T] {
	c := MulticloseMake[T](0)
	m.wg.Add(1)
	go func(c ChanReadCloser[T], n ChanWriter[T], onstop func()) {
		defer func() {
			c.Close()
			if onstop != nil {
				onstop()
			}
			m.wg.Done()
		}()
		for {
			select {
			case <-m.aggr.Context().Done():
				return
			case i := <-c.Read():
				// If we failed to write, then the channel is closed and we should be done
				if !n.Write(i) {
					return
				}
			}
		}
	}(c, m.aggr, onstop)
	return c
}

// Close the Merger and stop all goroutines associated with it.
func (m *Merger[T]) Close() {
	m.aggr.Close()
	m.wg.Wait()
}
