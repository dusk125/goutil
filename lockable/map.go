package lockable

import (
	"encoding/json"

	"golang.org/x/exp/maps"
)

type Map[K comparable, V any] struct {
	Locker[map[K]V]
}

func NewMap[K comparable, V any]() *Map[K, V] {
	m := &Map[K, V]{}
	m.UnsafeMake()
	return m
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

func (m *Map[K, V]) Clear() {
	m.Lock()
	defer m.Unlock()
	m.UnsafeClear()
}

func (m *Map[K, V]) UnsafeClear() {
	maps.Clear(m.item)
}

func (m *Map[K, V]) Clone() *Map[K, V] {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeClone()
}

func (m *Map[K, V]) UnsafeClone() *Map[K, V] {
	n := &Map[K, V]{}
	n.item = maps.Clone(m.item)
	return n
}

func (m *Map[K, V]) Copy(dst *Map[K, V]) {
	dst.Lock()
	defer dst.Unlock()
	m.RLock()
	defer m.RUnlock()
	m.UnsafeCopy(dst)
}

func (m *Map[K, V]) UnsafeCopy(dst *Map[K, V]) {
	maps.Copy(dst.item, m.item)
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

func (m *Map[K, V]) DeleteFunc(del func(K, V) bool) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeDeleteFunc(del)
}

func (m *Map[K, V]) UnsafeDeleteFunc(del func(K, V) bool) {
	maps.DeleteFunc(m.item, del)
}

func (m *Map[K, V]) EqualFunc(m2 *Map[K, V], eq func(V, V) bool) bool {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeEqualFunc(m2, eq)
}

func (m *Map[K, V]) UnsafeEqualFunc(m2 *Map[K, V], eq func(V, V) bool) bool {
	return maps.EqualFunc(m.item, m2.item, eq)
}

func (m *Map[K, V]) Keys() []K {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeKeys()
}

func (m *Map[K, V]) UnsafeKeys() []K {
	return maps.Keys(m.item)
}

func (m *Map[K, V]) Values() []V {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeValues()
}

func (m *Map[K, V]) UnsafeValues() []V {
	return maps.Values(m.item)
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
