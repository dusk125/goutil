// This package allows for easily setting up a "ctrl+c" interrupt waiter. By importing this package, the waiter is setup so the only thing you have to do is call 'keyinterrupt.Wait()' and your code will block until it is interrupted.
package keyinterrupt

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	wait   chan os.Signal
	ctx    context.Context
	cancel context.CancelFunc
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
	wait = make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-wait
		cancel()
	}()
}

func Cancel() {
	cancel()
}

func Wait() {
	<-ctx.Done()
}

func Context() context.Context {
	return ctx
}
