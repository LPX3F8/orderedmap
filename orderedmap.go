package orderedmap

import (
	"bytes"
	"sync"

	"github.com/LPX3F8/glist"
	j "github.com/json-iterator/go"
)

var _customJson = j.ConfigCompatibleWithStandardLibrary

type Filter[K comparable, V any] func(idx int, key K, val V) (want bool)
type Visitor[K comparable, V any] func(idx int, key K, val V) (skip bool)
type TravelMode uint

const (
	Forward TravelMode = iota
	Reverse
)

// OrderedMap use List[T] to ensure order
// The actual key-value pair exists in the basic map
// All operations lock objects and use read-write mutex
// to reduce lock competition.
type OrderedMap[K comparable, V any] struct {
	keys     *glist.List[K]
	elements map[K]*glist.Element[K]
	values   map[K]V
	mu       *sync.RWMutex
}

// New returns a pointer of *OrderedMap[K, V]
func New[K comparable, V any]() *OrderedMap[K, V] {
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

// Slice returns the elements slice
func (m *OrderedMap[K, V]) Slice(filters ...Filter[K, V]) []V {
	return m.slice(Forward, filters...)
}

// Reverse the elements in the array.
func (m *OrderedMap[K, V]) Reverse(filters ...Filter[K, V]) []V {
	return m.slice(Reverse, filters...)
}

// Slice returns the elements slice
func (m *OrderedMap[K, V]) slice(mode TravelMode, filters ...Filter[K, V]) []V {
	slice := make([]V, m.keys.Len())
	num := 0
	m.Travel(mode, func(idx int, key K, val V) bool {
		slice[num] = val
		num++
		return false
	}, filters...)
	return slice[:num]
}

func (m *OrderedMap[K, V]) TravelForward(f Visitor[K, V], filters ...Filter[K, V]) {
	m.Travel(Forward, f, filters...)
}

func (m *OrderedMap[K, V]) TravelReverse(f Visitor[K, V], filters ...Filter[K, V]) {
	m.Travel(Reverse, f, filters...)
}

func (m *OrderedMap[K, V]) Travel(mode TravelMode, f Visitor[K, V], filters ...Filter[K, V]) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var (
		idx  int
		skip bool
		drop bool
		hf   func() *glist.Element[K]
		nf   func(e *glist.Element[K]) *glist.Element[K]
		key  K
		val  V
	)

	if mode == Forward {
		hf = m.keys.Front
		nf = func(e *glist.Element[K]) *glist.Element[K] { return e.Next() }
	} else {
		hf = m.keys.Back
		nf = func(e *glist.Element[K]) *glist.Element[K] { return e.Prev() }
	}

	for e := hf(); e != nil; e = nf(e) {
		idx++
		key = e.Value
		val = m.values[e.Value]
		for _, filter := range filters {
			if drop = !filter(idx-1, key, val); drop {
				break
			}
		}
		if drop {
			continue
		}
		if skip = f(idx-1, key, val); skip {
			break
		}
	}
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
// Deprecated: Please replace it with the TravelForward.
func (m *OrderedMap[K, V]) Range(f func(key K, val V) bool) {
	m.TravelForward(func(idx int, k K, val V) bool { return f(k, val) })
	return
}

// Clear empty saved elements
func (m *OrderedMap[K, V]) Clear() *OrderedMap[K, V] {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.keys = glist.New[K]()
	m.elements = map[K]*glist.Element[K]{}
	m.values = map[K]V{}
	return m
}

// MarshalJSON implement the json.Marshaler interface.
// the interface implemented by types that can marshal themselves into valid JSON.
func (m *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	var err error
	var keyBytes, valBytes []byte
	buf := bytes.NewBuffer(nil)
	buf.WriteRune('{')
	for e := m.keys.Front(); e != nil; e = e.Next() {
		if keyBytes, err = _customJson.Marshal(e.Value); err != nil {
			break
		}
		if valBytes, err = _customJson.Marshal(m.values[e.Value]); err != nil {
			break
		}
		buf.Write(keyBytes)
		buf.WriteRune(':')
		buf.Write(valBytes)
		if e.Next() != nil {
			buf.WriteRune(',')
		}
	}
	if err != nil {
		return nil, err
	}
	buf.WriteRune('}')
	return buf.Bytes(), nil
}
