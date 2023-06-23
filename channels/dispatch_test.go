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

	onev, hasOne := channels.RecvWithContext(to, oneCh)
	twov, hasTwo := channels.RecvWithContext(to, twoCh)

	switch {
	case one == two:
		t.Errorf("channel ids shouldn't be equal. one: %v, two: %v", one, two)
	case !hasOne:
		t.Error("channel one didn't recieve in time")
	case !hasTwo:
		t.Error("channel two didn't recieve in time")
	case onev != twov:
		t.Error("Values should be equal")
	}
}
