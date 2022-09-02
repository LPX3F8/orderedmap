package orderedmap

import "github.com/LPX3F8/glist"

type Item[K comparable, V any] struct {
	k K
	v V
	e *glist.Element[K]
	m *OrderedMap[K, V]
}

func newItem[K comparable, V any](k K, v V, e *glist.Element[K], m *OrderedMap[K, V]) *Item[K, V] {
	return &Item[K, V]{k: k, v: v, e: e, m: m}
}

func (i *Item[K, V]) Value() V {
	return i.v
}

func (i *Item[K, V]) Key() K {
	return i.k
}

func (i *Item[K, V]) elements() *glist.Element[K] {
	return i.e
}

func (i *Item[K, V]) Prev() *Item[K, V] {
	i.m.mu.RLock()
	defer i.m.mu.RUnlock()
	if n := i.elements().Prev(); n != nil {
		return i.m.items[n.Value]
	}
	return nil
}

func (i *Item[K, V]) Next() *Item[K, V] {
	i.m.mu.RLock()
	defer i.m.mu.RUnlock()
	if n := i.elements().Next(); n != nil {
		return i.m.items[n.Value]
	}
	return nil
}
