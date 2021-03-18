package channels

import "sync"

type Merger struct {
	Aggr chan interface{}
	stop chan interface{}
	wg   sync.WaitGroup
}

func NewMerger() (m *Merger) {
	m = &Merger{
		Aggr: make(chan interface{}),
		stop: make(chan interface{}),
	}
	return
}

func (m *Merger) Add(c chan interface{}, onstop func()) {
	m.wg.Add(1)
	go func(stop, c <-chan interface{}) {
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

func (m *Merger) Close() {
	close(m.stop)
	m.wg.Wait()
	close(m.Aggr)
}
