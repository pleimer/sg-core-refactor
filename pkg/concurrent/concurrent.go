package concurrent

import "sync"

//Map thread safe map type
type Map struct {
	sync.RWMutex
	Items map[string]interface{} //Optimize with unsafe types?
}

//NewMap map constructor
func NewMap() *Map {
	return &Map{
		RWMutex: sync.RWMutex{},
		Items:   map[string]interface{}{},
	}
}

//Set set index in map
func (m *Map) Set(key string, value interface{}) {
	m.Lock()
	m.Items[key] = value
	m.Unlock()
}

//Get get item with key
func (m *Map) Get(key string) (interface{}, bool) {
	m.RLock()
	val, ok := m.Items[key]
	m.RUnlock()
	return val, ok
}

//Delete delete index in map
func (m *Map) Delete(key string) {
	m.Lock()
	delete(m.Items, key)
	m.Unlock()
}

//MapItem key value pair for use in iteration
type MapItem struct {
	Key   string
	Value interface{}
}

//Iter iterate with range keyword
func (m *Map) Iter() <-chan MapItem {
	c := make(chan MapItem)

	go func() {
		m.RLock()
		for k, v := range m.Items {
			c <- MapItem{k, v}
		}
		close(c)
		m.RUnlock()
	}()

	return c
}
