package channels

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	ErrExpectedClosed = errors.New("expected channel to be closed, but it was still open")
	ErrExpectedWrite  = errors.New("expected to write the channel, instead it failed to write")
)

func TestMultiGoroutineClose(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	ct, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ct.Done()
		m.Close()
	}()
	go func() {
		<-ct.Done()
		m.Close()
	}()
	cancel()
	time.Sleep(time.Second)
	if m.Open() {
		t.Error(ErrExpectedClosed)
	}
}

func TestChanReader(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	defer m.Close()

	go func(reader ChanReader[struct{}]) {
		<-reader.Read()
	}(m)

	if !m.Write(struct{}{}) {
		t.Error(ErrExpectedWrite)
	}
}

func TestChanWriter(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	defer m.Close()

	go func(writer ChanWriter[struct{}]) {
		if !writer.Write(struct{}{}) {
			t.Error(ErrExpectedWrite)
		}
	}(m)

	<-m.Read()
}

func TestChanCloser(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	defer m.Close()

	func(closer ChanCloser[struct{}]) {
		m.Close()
	}(m)

	if m.Open() {
		t.Error(ErrExpectedClosed)
	}
}

func TestChanReadCloser(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	defer m.Close()

	go func(rc ChanReadCloser[struct{}]) {
		<-rc.Read()
		rc.Close()
	}(m)

	if !m.Write(struct{}{}) {
		t.Error(ErrExpectedWrite)
	}

	if m.Open() {
		t.Error(ErrExpectedClosed)
	}
}

func TestChanWriteCloser(t *testing.T) {
	m := MulticloseMake[int](0)
	defer m.Close()

	go func(wc ChanWriteCloser[int]) {
		if !wc.Write(1) {
			t.Error(ErrExpectedWrite)
		}
		wc.Close()
	}(m)

	i := <-m.Read()
	if i != 1 {
		t.Errorf("expected read value to be 1, was: %v", i)
	}
	i, ok := <-m.Read()
	if i != 0 {
		t.Errorf("expected read value to be 0, was: %v", i)
	}
	if ok {
		t.Error("expected channel to be invalid, was valid")
	}
	if m.Open() {
		t.Error(ErrExpectedClosed)
	}
}

func TestChanReadWriter(t *testing.T) {
	// If we don't buffer the channel, the following test will deadlock
	m := MulticloseMake[int](1)
	defer m.Close()

	func(rw ChanReadWriter[int]) {
		if !rw.Write(1) {
			t.Error(ErrExpectedWrite)
		}
		i := <-rw.Read()
		if i != 1 {
			t.Errorf("expected read value to be 1, was: %v", i)
		}
	}(m)
}

func TestChanReadWriteCloser(t *testing.T) {
	// If we don't buffer the channel, the following test will deadlock
	m := MulticloseMake[int](1)
	defer m.Close()

	func(rwc ChanReadWriteCloser[int]) {
		if !rwc.Write(1) {
			t.Error(ErrExpectedWrite)
		}
		i := <-rwc.Read()
		if i != 1 {
			t.Errorf("expected read value to be 1, was: %v", i)
		}
		rwc.Close()
		if m.Open() {
			t.Error(ErrExpectedClosed)
		}
	}(m)
}

func TestChanBuffer(t *testing.T) {
	m := MulticloseMake[struct{}](2)
	defer m.Close()

	if !m.Write(struct{}{}) {
		t.Error(ErrExpectedWrite)
	}
	if !m.Write(struct{}{}) {
		t.Error(ErrExpectedWrite)
	}
	go func(r ChanReader[struct{}]) {
		for i := 0; i < 3; i++ {
			<-r.Read()
		}
	}(m)
	if !m.Write(struct{}{}) {
		t.Error(ErrExpectedWrite)
	}
}

func TestChanClosedWrite(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	m.Close()

	if m.Open() {
		t.Error(ErrExpectedClosed)
	}

	if m.Write(struct{}{}) {
		t.Error("expected to not write to the channel since it's closed, instead it wrote to the channel")
	}
}

func TestZeroChan(t *testing.T) {
	m := MulticloseChan[int]{}
	defer m.Close()

	if m.Open() {
		t.Error("expected channel to be closed, was open")
	}

	m.Make(1)

	if !m.Write(1) {
		t.Error(ErrExpectedWrite)
	}
	i := <-m.Read()
	if i != 1 {
		t.Errorf("expected read value to be 1, was: %v", i)
	}
}

func TestMultiOpen(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	defer m.Close()

	c := &m.ch
	m.Make(1)
	if &m.ch != c {
		t.Error("second Open was not a no-op as expected")
	}
}

func TestNotifyClosed(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	defer m.Close()

	var complete uint32
	wg := sync.WaitGroup{}
	dl, _ := t.Deadline()
	ct, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel()

	wg.Add(1)
	go func(complete *uint32) {
		defer wg.Done()
		select {
		case <-ct.Done():
			t.Error("failed to get notify closed within the deadline")
			return
		case <-m.Context().Done():
			atomic.StoreUint32(complete, 1)
			return
		}
	}(&complete)

	m.Close()
	wg.Wait()
	if complete != 1 {
		t.Errorf("expected complete to be 1, was %v", complete)
	}
}

func TestChanContext(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	defer m.Close()

	var complete uint32
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(complete *uint32) {
		defer func() {
			atomic.StoreUint32(complete, 1)
			wg.Done()
		}()
		<-m.Context().Done()
	}(&complete)

	m.Close()
	wg.Wait()

	if complete != 1 {
		t.Errorf("expected complete to be 1, was %v", complete)
	}
}
