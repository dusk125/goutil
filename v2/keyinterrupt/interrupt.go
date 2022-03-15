package keyinterrupt

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	wait chan os.Signal
)

func init() {
	wait = make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)
}

func Wait() {
	defer close(wait)
	<-wait
}
