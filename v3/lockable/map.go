package lockable

import (
	"encoding/json"
	"sync"
)

type Map[K comparable, V any] struct {
	l sync.RWMutex
	e map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{e: make(map[K]V)}
}

func (m *Map[K, V]) Make() {
	m.l.Lock()
	defer m.l.Unlock()
	m.e = make(map[K]V)
}

func (m *Map[K, V]) Get(k K) (v V, has bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	v, has = m.e[k]
	return
}

func (m *Map[K, V]) Set(k K, v V) {
	m.l.Lock()
	defer m.l.Unlock()
	m.e[k] = v
}

func (m *Map[K, V]) Len() int {
	m.l.RLock()
	defer m.l.RUnlock()
	return len(m.e)
}

func (m *Map[K, V]) Delete(k K) (v V, has bool) {
	m.l.Lock()
	defer m.l.Unlock()
	v, has = m.e[k]
	delete(m.e, k)
	return
}

func (m *Map[K, V]) Foreach(f func(k K, v V)) {
	m.l.RLock()
	defer m.l.RUnlock()
	for k, v := range m.e {
		f(k, v)
	}
}

func (l *Map[K, V]) MarshalJSON() ([]byte, error) {
	l.l.RLock()
	defer l.l.RUnlock()
	return json.Marshal(l.e)
}

func (l *Map[K, V]) UnMarshalJSON(b []byte) error {
	l.l.Lock()
	defer l.l.Unlock()
	return json.Unmarshal(b, &l.e)
}
