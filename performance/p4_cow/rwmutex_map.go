package cow

import "sync"

type RWMutexMap struct {
	mu    sync.RWMutex
	dirty map[interface{}]interface{}
}

func (m *RWMutexMap) Load(key interface{}) (interface{}, bool) {
	m.mu.RLock()
	value, ok := m.dirty[key]
	m.mu.RUnlock()
	return value, ok
}

func (m *RWMutexMap) Store(key, value interface{}) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[interface{}]interface{})
	}
	m.dirty[key] = value
	m.mu.Unlock()
}

func (m *RWMutexMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	m.mu.Lock()
	actual, loaded = m.dirty[key]
	if !loaded {
		actual = value
		if m.dirty == nil {
			m.dirty = make(map[interface{}]interface{})
		}
		m.dirty[key] = value
	}
	m.mu.Unlock()
	return actual, loaded
}

func (m *RWMutexMap) Delete(key interface{}) {
	m.mu.Lock()
	delete(m.dirty, key)
	m.mu.Unlock()
}

func (m *RWMutexMap) Range(f func(key, value interface{}) (shouldContinue bool)) {
	m.mu.RLock()
	keys := make([]interface{}, 0, len(m.dirty))
	for k := range m.dirty {
		keys = append(keys, k)
	}
	m.mu.RUnlock()

	for _, k := range keys {
		v, ok := m.Load(k)
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}
