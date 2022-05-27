package channels

import (
	"context"
	"errors"
	"testing"
	"time"
)

var (
	ErrExpectedClosed = errors.New("expected channel to be closed, but it was still open")
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

	m.Write(struct{}{})
}

func TestChanWriter(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	defer m.Close()

	go func(writer ChanWriter[struct{}]) {
		writer.Write(struct{}{})
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

	m.Write(struct{}{})

	if m.Open() {
		t.Error(ErrExpectedClosed)
	}
}

func TestChanWriteCloser(t *testing.T) {
	m := MulticloseMake[int](0)
	defer m.Close()

	go func(wc ChanWriteCloser[int]) {
		wc.Write(1)
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
		rw.Write(1)
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
		rwc.Write(1)
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

	m.Write(struct{}{})
	m.Write(struct{}{})
	go func(r ChanReader[struct{}]) {
		for i := 0; i < 3; i++ {
			<-r.Read()
		}
	}(m)
	m.Write(struct{}{})
}

func TestChanClosedWrite(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	m.Close()

	if m.Open() {
		t.Error(ErrExpectedClosed)
	}

	m.Write(struct{}{})
}

func TestZeroChan(t *testing.T) {
	m := MulticloseChan[int]{}
	defer m.Close()

	if m.Open() {
		t.Error("expected channel to be closed, was open")
	}

	m.Make(1)

	m.Write(1)
	i := <-m.Read()
	if i != 1 {
		t.Errorf("expected read value to be 1, was: %v", i)
	}
}

func TestMultiOpen(t *testing.T) {
	m := MulticloseMake[struct{}](0)
	defer m.Close()

	m.Open()
	c := &m.ch
	m.Open()
	if &m.ch != c {
		t.Error("second Open was not a no-op as expected")
	}
}
