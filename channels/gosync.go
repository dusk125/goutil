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

func (s *GoSync) Add(fs ...func(stop <-chan interface{})) {
	for _, f := range fs {
		go func(stop <-chan interface{}, f func(stop <-chan interface{})) {
			defer s.group.CloseAll()
			f(stop)
		}(s.group.Chan(), f)
	}
}

func (s *GoSync) StopAll() {
	s.group.CloseAll()
}

func (s *GoSync) Wait() {
	s.group.Wait()
}
