package channels

import (
	"sync/atomic"

	"github.com/dusk125/goutil/v4/lockable"
)

type Broadcaster[T any] struct {
	counter atomic.Uint64
	members lockable.Map[uint64, chan T]
}

func (b *Broadcaster[T]) Write(m T) {
	if b.members.Nil() {
		b.members.Make()
	}
	b.members.Foreach(func(k uint64, v chan T) {
		v <- m
	})
}

func (b *Broadcaster[T]) Add() (id uint64, ch <-chan T) {
	c := make(chan T, 1)
	id = b.counter.Add(1)
	b.members.Safe(true, func() {
		if b.members.UnsafeNil() {
			b.members.UnsafeMake()
		}
		b.members.UnsafeSet(id, c)
	})
	return id, c
}

func (b *Broadcaster[T]) Remove(id uint64) {
	b.members.Safe(true, func() {
		if b.members.UnsafeNil() {
			b.members.UnsafeMake()
		}
		b.members.UnsafeDelete(id)
	})
}
