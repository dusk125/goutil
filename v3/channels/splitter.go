package channels

type splitadd[T any] struct {
	id string
	ch chan T
}

type Splitter[T any] struct {
	Dispatch   chan T
	entries    map[string]chan T
	closeCheck chan struct{}
	ar         chan splitadd[T]
	OnEmpty    func()
}

func NewSplitter[T any]() (s *Splitter[T]) {
	s = &Splitter[T]{
		Dispatch:   make(chan T),
		closeCheck: make(chan struct{}),
		entries:    make(map[string]chan T),
		ar:         make(chan splitadd[T]),
		OnEmpty:    func() {},
	}
	go s.handle()
	return
}

func (s *Splitter[T]) handle() {
	for {
		select {
		case ar, ok := <-s.ar:
			if !ok {
				return
			}
			if ch, has := s.entries[ar.id]; has {
				close(ch)
				delete(s.entries, ar.id)
				s.OnEmpty()
			} else if ar.ch != nil {
				s.entries[ar.id] = ar.ch
			}
		case item, ok := <-s.Dispatch:
			if !ok {
				return
			}
			for _, entry := range s.entries {
				entry <- item
			}
		}
	}
}

func (s *Splitter[T]) Add(id string) <-chan T {
	c := make(chan T)
	s.ar <- splitadd[T]{
		id: id,
		ch: c,
	}
	return c
}

func (s *Splitter[T]) Remove(id string) {
	s.ar <- splitadd[T]{
		id: id,
	}
}

func (s *Splitter[T]) IsOpen() bool {
	select {
	case <-s.closeCheck:
		return false
	default:
		return true
	}
}

func (s *Splitter[T]) Close() {
	close(s.closeCheck)
	close(s.Dispatch)
	for _, ch := range s.entries {
		close(ch)
	}
}
