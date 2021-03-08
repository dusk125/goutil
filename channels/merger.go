package channels

import (
	"sync"

	"github.com/dusk125/goutil/snowflake"
)

type Merger struct {
	Aggr    chan interface{}
	entries SafeChanMap
	wg      sync.WaitGroup
}

func NewMerger() (m *Merger) {
	m = &Merger{
		Aggr: make(chan interface{}),
	}
	m.entries.Make()
	return
}

func (m *Merger) Add(c *SafeChan) {
	m.wg.Add(1)
	var (
		id = snowflake.New("chan")
	)
	m.entries.Put(id, c)
	go func() {
		defer func() {
			m.entries.Delete(id)
			m.wg.Done()
		}()
		for item := range c.Read() {
			m.Aggr <- item
		}
	}()
}

func (m *Merger) Close() {
	m.entries.Clear(func(k string, v *SafeChan) {
		v.Close()
	})
	m.wg.Wait()
	close(m.Aggr)
}
