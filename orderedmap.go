package orderedmap

import (
	"sync"

	"github.com/LPX3F8/glist"
)

// OrderedMap use List[T] to ensure order
// The actual key-value pair exists in the basic map
// All operations lock objects and use read-write mutex
// to reduce lock competition.
type OrderedMap[K, V comparable] struct {
	keys     *glist.List[K]
	elements map[K]*glist.Element[K]
	values   map[K]V
	mu       *sync.RWMutex
}

// New returns a pointer of *OrderedMap[K, V]
func New[K, V comparable]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:     glist.New[K](),
		elements: map[K]*glist.Element[K]{},
		values:   map[K]V{},
		mu:       new(sync.RWMutex),
	}
}

// Store key-value pair
func (m *OrderedMap[K, V]) Store(k K, v V) *OrderedMap[K, V] {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.values[k]; !ok {
		m.elements[k] = m.keys.PushBack(k)
	}
	m.values[k] = v
	return m
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *OrderedMap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	if actual, loaded = m.Load(k); loaded {
		return actual, loaded
	}
	m.Store(k, v)
	return actual, loaded
}

// Has return key exists
func (m *OrderedMap[K, V]) Has(k K) bool {
	_, ok := m.Load(k)
	return ok
}

// Load returns the value stored in the map for a key,
// or zero-value if no value is present, It depends on the data type.
// The ok result indicates whether value was found in the map.
func (m *OrderedMap[K, V]) Load(k K) (val V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok = m.values[k]
	return
}

// Delete removes key-value pair
func (m *OrderedMap[K, V]) Delete(k K) *OrderedMap[K, V] {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ele, ok := m.elements[k]; ok {
		m.keys.Remove(ele)
		delete(m.elements, k)
		delete(m.values, k)
	}
	return m
}

// Len return the map key size
func (m *OrderedMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.keys.Len()
}

// Slice returns the values slice
func (m *OrderedMap[K, V]) Slice() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	idx := 0
	slice := make([]V, m.keys.Len())
	for e := m.keys.Front(); e != nil; e = e.Next() {
		slice[idx] = m.values[e.Value]
		idx++
	}
	return slice
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *OrderedMap[K, V]) Range(f func(key K, val V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for e := m.keys.Front(); e != nil; e = e.Next() {
		if !f(e.Value, m.values[e.Value]) {
			break
		}
	}
	return
}
