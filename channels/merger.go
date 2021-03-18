package channels

import (
	"sync"

	"github.com/dusk125/goutil/logger"
	"github.com/dusk125/goutil/snowflake"
)

type Merger struct {
	Aggr    chan interface{}
	entries sync.Map
	wg      sync.WaitGroup
}

func NewMerger() (m *Merger) {
	m = &Merger{
		Aggr: make(chan interface{}),
	}
	// m.entries.Make()
	return
}

func (m *Merger) Add(c chan interface{}) {
	m.wg.Add(1)
	var (
		id = snowflake.New("chan")
	)
	// m.entries.Put(id, c)
	m.entries.Store(id, c)
	go func() {
		defer func() {
			logger.Debug("Calling delete")
			m.entries.Delete(id)
			logger.Debug("delete done")
			m.wg.Done()
			logger.Debug("done done")
		}()
		for item := range c {
			m.Aggr <- item
		}
	}()
}

func (m *Merger) Close() {
	logger.Debug("Calling range")
	m.entries.Range(func(key, value interface{}) bool {
		close(value.(chan interface{}))
		return true
	})
	// m.entries.Range(func(key string, c *SafeChan) (needDelete bool) {
	// 	c.Close()
	// 	return
	// })
	logger.Debug("range done")
	m.wg.Wait()
	logger.Debug("done waiting")
	close(m.Aggr)
}
