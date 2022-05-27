package channels

import "sync"

type Merger[T any] struct {
	Aggr chan T
	stop chan struct{}
	wg   sync.WaitGroup
}

func NewMerger[T any]() (m *Merger[T]) {
	m = &Merger[T]{
		Aggr: make(chan T),
		stop: make(chan struct{}),
	}
	return
}

func (m *Merger[T]) Add(c <-chan T, onstop func()) {
	m.wg.Add(1)
	go func(stop chan struct{}, c <-chan T) {
		defer m.wg.Done()
		for {
			select {
			case <-stop:
				onstop()
				return
			case item, ok := <-c:
				if !ok {
					return
				}

				m.Aggr <- item
			}
		}
	}(m.stop, c)
}

func (m *Merger[T]) Close() {
	close(m.stop)
	m.wg.Wait()
	close(m.Aggr)
}
