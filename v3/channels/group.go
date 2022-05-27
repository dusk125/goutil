package channels

import "sync"

type Group struct {
	sync.Mutex
	sync.WaitGroup
	chans []chan struct{}
}

func NewGroup() (g *Group) {
	g = &Group{
		chans: make([]chan struct{}, 0),
	}
	return
}

func (g *Group) Chan() (c chan struct{}) {
	g.Lock()
	defer g.Unlock()

	c = make(chan struct{})
	g.chans = append(g.chans, c)
	g.Add(1)
	go func(c <-chan struct{}) {
		defer g.Done()
		<-c
	}(c)

	return
}

func (g *Group) CloseAll() {
	g.Lock()
	defer g.Unlock()

	for _, c := range g.chans {
		close(c)
	}
	g.chans = g.chans[:0]
}
