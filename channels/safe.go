package channels

import (
	"log"
	"os"

	"github.com/tevino/abool"
)

type SafeChanOpts struct {
	Debug bool
}

type SafeChan struct {
	ch     chan interface{}
	closed *abool.AtomicBool
	opts   SafeChanOpts
	logger *log.Logger
}

func MakeSafeChan(ch chan interface{}, opts ...SafeChanOpts) (safe *SafeChan) {
	var (
		opt SafeChanOpts
	)
	if len(opts) > 0 {
		opt = opts[0]
	}
	safe = &SafeChan{
		ch:     ch,
		opts:   opt,
		closed: abool.New(),
		logger: log.New(os.Stdout, "[SafeChannel]", log.Flags()),
	}
	return safe
}

func NewSafeChan(opts ...SafeChanOpts) (safe *SafeChan) {
	return MakeSafeChan(make(chan interface{}), opts...)
}

func (s *SafeChan) log(msg ...interface{}) {
	if s.opts.Debug {
		s.logger.Println(msg...)
	}
}

func (s *SafeChan) Closed() bool {
	return s.closed.IsSet()
}

func (s *SafeChan) Read() (ch chan interface{}) {
	return s.ch
}

func (s *SafeChan) Send(msg interface{}) (sent bool) {
	if s.closed.IsSet() {
		s.log("Attempt to send on closed channel")
		return
	}

	s.ch <- msg

	return true
}

func (s *SafeChan) Close() {
	if s.closed.IsNotSet() {
		s.closed.Set()
		close(s.ch)
		return
	}

	s.log("Attempt to close channel twice")
}
