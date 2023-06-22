package channels

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMerger(t *testing.T) {
	m := NewMerger[int]()
	defer m.Close()

	// First, we need to create multiple writers to the Merger and have them write a value into it.

	var complete uint32
	wg := sync.WaitGroup{}
	dl, _ := t.Deadline()
	ct, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel()
	for i := 1; i < 3; i++ {
		wg.Add(1)
		go func(ct context.Context, i int, complete *uint32) {
			c := m.Add(func() {
				atomic.AddUint32(complete, 1)
				wg.Done()
			})
			c.Write(i)
			<-ct.Done()
		}(ct, i, &complete)
	}

	// Next, we need to read those values and ensure that we got the correct number of them.

	to, toCancel := context.WithTimeout(ct, time.Second*5)
	defer toCancel()
	total := 0
Loop:
	for {
		select {
		case <-to.Done():
			t.Error("expected to read values before timeout, instead timed out")
			return
		case i := <-m.Aggr().Read():
			if total += i; total == 3 {
				select {
				case <-m.Aggr().Read():
					t.Error("expected the read channel to block, instead it had a value ready")
				default:
				}
				break Loop
			}
		}
	}

	m.Close()
	wg.Wait()
	if complete != 2 {
		t.Errorf("expected number of complete gorountines to be 2, was: %v", complete)
	}
}
