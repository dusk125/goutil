package channels

import "sync"

type Group struct {
	sync.Mutex
	sync.WaitGroup
	chans []chan interface{}
}

func NewGroup() (g *Group) {
	g = &Group{
		chans: make([]chan interface{}, 0),
	}
	return
}

func (g *Group) Chan() (c chan interface{}) {
	g.Lock()
	defer g.Unlock()

	c = make(chan interface{})
	g.chans = append(g.chans, c)
	g.Add(1)
	go func(c <-chan interface{}) {
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
