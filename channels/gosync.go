package channels

type GoSync struct {
	group *Group
}

func NewGoSync() (s *GoSync) {
	s = &GoSync{
		group: NewGroup(),
	}
	return
}

func (s *GoSync) Add(f func(stop <-chan interface{})) {
	go func(stop <-chan interface{}) {
		defer s.group.CloseAll()
		f(stop)
	}(s.group.Chan())
}

func (s *GoSync) StopAll() {
	s.group.CloseAll()
}

func (s *GoSync) Wait() {
	s.group.Wait()
}
