package channels

import "context"

func RecvWithContext[T any](timeout context.Context, ch <-chan T) (v T, has bool) {
	select {
	case <-timeout.Done():
		return
	case v, has = <-ch:
		return
	}
}
