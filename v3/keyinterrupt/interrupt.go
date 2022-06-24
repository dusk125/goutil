// This package allows for easily setting up a "ctrl+c" interrupt waiter. By importing this package, the waiter is setup so the only thing you have to do is call 'keyinterrupt.Wait()' and your code will block until it is interrupted.
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
	<-wait
}
