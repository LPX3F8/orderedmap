package orderedmap

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkOrderedMap(b *testing.B) {
	benchmarkOrderedMap(b, false, false)
}

func BenchmarkOrderedMapSlack(b *testing.B) {
	benchmarkOrderedMap(b, true, false)
}

func BenchmarkOrderedMapWork(b *testing.B) {
	benchmarkOrderedMap(b, false, true)
}

func BenchmarkOrderedMapWorkSlack(b *testing.B) {
	benchmarkOrderedMap(b, true, true)
}

func benchmarkOrderedMap(b *testing.B, slack, work bool) {
	m := New[int, int]()
	if slack {
		b.SetParallelism(10)
	}
	rand.Seed(time.Now().UnixNano())
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		foo := 0
		for pb.Next() {
			k := rand.Int()
			v := rand.Int()
			m.Store(k, v)
			if nv, _ := m.Load(k); nv != v {
				panic("NOT EQ!!")
			}
			m.Delete(k)
			if m.Has(k) {
				panic("KEY EXISTS")
			}
			if work {
				for i := 0; i < 100; i++ {
					foo *= 2
					foo /= 2
				}
			}
		}
		_ = foo
	})
}

func BenchmarkNativeSyncMap_Store(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	b.ReportAllocs()

	m := new(sync.Map)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := rand.Int()
			v := rand.Int()
			m.Store(k, v)
			if _, ok := m.Load(k); !ok {
				panic("KEY NOT SET")
			}
		}
	})
}

func BenchmarkOrderedMap_Store(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	b.ReportAllocs()

	m := New[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := rand.Int()
			v := rand.Int()
			m.Store(k, v)
			if !m.Has(k) {
				panic("KEY NOT SET")
			}
		}
	})
}

func BenchmarkNativeSyncMap_LoadOrStore(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	b.ReportAllocs()

	m := new(sync.Map)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := rand.Int()
			v := rand.Int()
			m.LoadOrStore(k, v)
			if _, ok := m.Load(k); !ok {
				panic("KEY NOT SET")
			}
		}
	})
}

func BenchmarkOrderedMap_LoadOrStore(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	b.ReportAllocs()

	m := New[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := rand.Int()
			v := rand.Int()
			m.LoadOrStore(k, v)
			if !m.Has(k) {
				panic("KEY NOT SET")
			}
		}
	})
}

func BenchmarkNativeSyncMap_Delete(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	b.ReportAllocs()

	m := new(sync.Map)
	for i := 0; i < 1000000; i++ {
		k := rand.Int()
		v := rand.Int()
		m.Store(k, v)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var k any
			m.Range(func(key, value any) bool {
				k = key
				return false
			})
			m.Delete(k)
		}
	})
}

func BenchmarkOrderedMap_Delete(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	b.ReportAllocs()

	m := New[int, int]()
	for i := 0; i < 1000000; i++ {
		k := rand.Int()
		v := rand.Int()
		m.Store(k, v)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var key int
			m.Range(func(k int, v int) bool {
				key = k
				return false
			})
			m.Delete(key)
		}
	})
}
