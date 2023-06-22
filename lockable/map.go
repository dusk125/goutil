package lockable

import (
	"encoding/json"
)

type Map[K comparable, V any] struct {
	Locker[map[K]V]
}

func (m *Map[K, V]) Make() {
	m.Lock()
	defer m.Unlock()
	m.UnsafeMake()
}

func (m *Map[K, V]) UnsafeMake() {
	m.item = make(map[K]V)
}

func (m *Map[K, V]) Nil() bool {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeNil()
}

func (m *Map[K, V]) UnsafeNil() bool {
	return m.item == nil
}

func (m *Map[K, V]) Get(k K) (v V, has bool) {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeGet(k)
}

func (m *Map[K, V]) UnsafeGet(k K) (v V, has bool) {
	v, has = m.item[k]
	return
}

func (m *Map[K, V]) Set(k K, v V) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeSet(k, v)
}

func (m *Map[K, V]) UnsafeSet(k K, v V) {
	m.item[k] = v
}

func (m *Map[K, V]) Has(k K) bool {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeHas(k)
}

func (m *Map[K, V]) UnsafeHas(k K) bool {
	_, has := m.item[k]
	return has
}

func (m *Map[K, V]) Len() int {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeLen()
}

func (m *Map[K, V]) UnsafeLen() int {
	return len(m.item)
}

func (m *Map[K, V]) Delete(k K) (v V, has bool) {
	m.Lock()
	defer m.Unlock()
	return m.UnsafeDelete(k)
}

func (m *Map[K, V]) UnsafeDelete(k K) (v V, has bool) {
	v, has = m.item[k]
	delete(m.item, k)
	return
}

func (m *Map[K, V]) Keys() []K {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeKeys()
}

func (m *Map[K, V]) UnsafeKeys() []K {
	ks := make([]K, 0, len(m.item))
	for k := range m.item {
		ks = append(ks, k)
	}
	return ks
}

func (m *Map[K, V]) Values() []V {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeValues()
}

func (m *Map[K, V]) UnsafeValues() []V {
	vs := make([]V, 0, len(m.item))
	for _, v := range m.item {
		vs = append(vs, v)
	}
	return vs
}

func (m *Map[K, V]) Foreach(f func(k K, v V)) {
	m.RLock()
	defer m.RUnlock()
	m.UnsafeForeach(f)
}

func (m *Map[K, V]) UnsafeForeach(f func(k K, v V)) {
	for k, v := range m.item {
		f(k, v)
	}
}

func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeMarshalJSON()
}

func (m *Map[K, V]) UnsafeMarshalJSON() ([]byte, error) {
	return json.Marshal(m.item)
}

func (m *Map[K, V]) UnmarshalJSON(b []byte) error {
	m.Lock()
	defer m.Unlock()
	return m.UnsafeUnmarshalJSON(b)
}

func (m *Map[K, V]) UnsafeUnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.item)
}
