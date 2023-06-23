package channels_test

import (
	"context"
	"testing"
	"time"

	"github.com/dusk125/goutil/v4/channels"
)

func TestBroadcaster(t *testing.T) {
	broadcast := channels.Broadcaster[int]{}
	one, oneCh := broadcast.Add()
	two, twoCh := broadcast.Add()

	broadcast.Write(1)

	to, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	onev, _, oneCancelled := channels.RecvWithContext(to, oneCh)
	twov, _, twoCancelled := channels.RecvWithContext(to, twoCh)

	switch {
	case one == two:
		t.Errorf("channel ids shouldn't be equal. one: %v, two: %v", one, two)
	case !oneCancelled:
		t.Error("channel one didn't recieve in time")
	case !twoCancelled:
		t.Error("channel two didn't recieve in time")
	case onev != twov:
		t.Error("Values should be equal")
	}
}
