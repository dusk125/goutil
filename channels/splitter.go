package channels

type Splitter struct {
	Dispatch   chan interface{}
	entries    SafeChanMap
	closeCheck chan struct{}
}

func NewSplitter() (s *Splitter) {
	s = &Splitter{
		Dispatch:   make(chan interface{}),
		closeCheck: make(chan struct{}),
	}
	s.entries.Make()
	go s.handle()
	return
}

func (s *Splitter) handle() {
	for item := range s.Dispatch {
		s.entries.Range(func(k string, v *SafeChan) bool {
			return !v.Send(item)
		})
	}
}

func (s *Splitter) Add(id string) (c *SafeChan) {
	c = NewSafeChan()
	s.entries.Put(id, c)
	return c
}

func (s *Splitter) Count() int {
	return s.entries.Len()
}

func (s *Splitter) IsOpen() bool {
	select {
	case <-s.closeCheck:
		return false
	default:
		return true
	}
}

func (s *Splitter) Close() {
	close(s.closeCheck)
	close(s.Dispatch)
	s.entries.Clear(func(k string, v *SafeChan) {
		v.Close()
	})
}
