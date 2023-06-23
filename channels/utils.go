package channels

import "context"

func RecvWithContext[T any](timeout context.Context, ch <-chan T) (v T, open, cancelled bool) {
	select {
	case <-timeout.Done():
		cancelled = true
		return
	case v, open = <-ch:
		return
	}
}
