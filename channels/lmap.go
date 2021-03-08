package channels

import (
	"sync"
)

type SafeChanMap struct {
	sync.RWMutex
	m map[string]*SafeChan
}

func (m *SafeChanMap) Make() {
	m.Lock()
	defer m.Unlock()
	m.m = make(map[string]*SafeChan)
}

func (m *SafeChanMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.m)
}

func (m *SafeChanMap) Has(k string) (has bool) {
	m.RLock()
	defer m.RUnlock()
	_, has = m.m[k]
	return
}

func (m *SafeChanMap) Put(k string, v *SafeChan) (o *SafeChan, had bool) {
	m.Lock()
	defer m.Unlock()
	o, had = m.m[k]
	m.m[k] = v
	return
}

func (m *SafeChanMap) Get(k string) (v *SafeChan, has bool) {
	m.RLock()
	defer m.RUnlock()
	v, has = m.m[k]
	return
}

func (m *SafeChanMap) GetOrPut(k string, p *SafeChan) (v *SafeChan, put bool) {
	m.Lock()
	defer m.Unlock()
	var (
		has bool
	)

	if v, has = m.m[k]; !has {
		m.m[k] = p
		v = p
		put = true
	}
	return
}

func (m *SafeChanMap) Range(h func(string, *SafeChan) bool, keys ...string) {
	m.Lock()
	defer m.Unlock()

	if len(keys) == 0 {
		for _, key := range keys {
			if v, has := m.m[key]; has {
				if needDelete := h(key, v); needDelete {
					delete(m.m, key)
				}
			}
		}
	} else {
		for k, v := range m.m {
			if needDelete := h(k, v); needDelete {
				delete(m.m, k)
			}
		}
	}
}

func (m *SafeChanMap) Values() (vals []*SafeChan) {
	m.RLock()
	defer m.RUnlock()

	vals = make([]*SafeChan, len(m.m))
	i := 0
	for _, v := range m.m {
		vals[i] = v
		i++
	}
	return
}

func (m *SafeChanMap) Clear(h func(k string, v *SafeChan)) {
	m.Range(func(k string, v *SafeChan) bool {
		h(k, v)
		return true
	})
}

func (m *SafeChanMap) Delete(k string) (v *SafeChan, had bool) {
	m.Lock()
	defer m.Unlock()
	v, had = m.m[k]
	delete(m.m, k)
	return
}
