package channels

type splitadd struct {
	id string
	ch chan interface{}
}

type Splitter struct {
	Dispatch   chan interface{}
	entries    map[string]chan interface{}
	closeCheck chan struct{}
	ar         chan splitadd
}

func NewSplitter() (s *Splitter) {
	s = &Splitter{
		Dispatch:   make(chan interface{}),
		closeCheck: make(chan struct{}),
		entries:    make(map[string]chan interface{}),
		ar:         make(chan splitadd),
	}
	go s.handle()
	return
}

func (s *Splitter) handle() {
	for {
		select {
		case item, ok := <-s.Dispatch:
			if !ok {
				return
			}
			for _, entry := range s.entries {
				entry <- item
			}
		case ar, ok := <-s.ar:
			if !ok {
				return
			}
			if ch, has := s.entries[ar.id]; has {
				close(ch)
				delete(s.entries, ar.id)
			} else {
				s.entries[ar.id] = ar.ch
			}
		}
	}
}

func (s *Splitter) Add(id string) <-chan interface{} {
	c := make(chan interface{})
	s.ar <- splitadd{
		id: id,
		ch: c,
	}
	return c
}

func (s *Splitter) Remove(id string) {
	s.ar <- splitadd{
		id: id,
	}
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
	for _, ch := range s.entries {
		close(ch)
	}
}
