package orderedmap

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderedMap(t *testing.T) {
	a := assert.New(t)
	m := New[int, int]()

	testKVList := make([][]int, 0)
	for i := 1; i <= 26; i++ {
		testKVList = append(testKVList, []int{i, i})
		m.Store(i, i)
	}

	idx := 0
	m.Range(func(key int, value int) bool {
		a.Equal(testKVList[idx][0], key)
		a.Equal(testKVList[idx][1], value)
		idx++
		return true
	})

	m.Delete(20)
	val, has := m.Load(20)
	a.Equal(0, val)
	a.False(has)

	wg := new(sync.WaitGroup)
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			hammerOrderedMap(t, m, 1000, wg)
		}()
	}
	wg.Wait()
	a.Equal(25, m.Len())

	l := m.Slice()
	a.Equal(len(l), m.Len())
	rl := m.Reverse()
	a.Equal(len(rl), m.Len())
	a.Equal(0, len(m.Clear().Slice()))
	a.Equal(0, len(m.Clear().Reverse()))
}

func hammerOrderedMap(t *testing.T, m *OrderedMap[int, int], loops int, group *sync.WaitGroup) {
	var key int
	var v int
	for i := 0; i < loops; i++ {
		rand.Seed(time.Now().UnixNano())
		key = rand.Int()
		v = rand.Int()
		m.Store(key, v)
		if val, ok := m.Load(key); !ok {
			t.Fatalf("key not found")
		} else {
			assert.Equal(t, val, v)
		}
		m.Delete(key)
		if m.Has(key) {
			t.Fatalf("key found ")
		}
	}
	group.Done()
}

func TestOrderedMap_MarshalJSON(t *testing.T) {
	a := assert.New(t)
	bytes1 := []byte(`{"key1":1,"key2":2,"key3":3}`)
	m := New[string, int]()
	m.Store("key1", 1)
	m.Store("key2", 2)
	m.Store("key3", 3)
	d, err := json.Marshal(m)
	a.Equal(d, bytes1)
	a.NoError(err)

	m2 := New[struct {
		Key     string `json:"key"`
		KeyInfo string `json:"keyInfo"`
	}, int]()

	m2.Store(struct {
		Key     string `json:"key"`
		KeyInfo string `json:"keyInfo"`
	}{Key: "key1", KeyInfo: "k1info"}, 1).
		Store(struct {
			Key     string `json:"key"`
			KeyInfo string `json:"keyInfo"`
		}{Key: "key2", KeyInfo: "k2info"}, 2).
		Store(struct {
			Key     string `json:"key"`
			KeyInfo string `json:"keyInfo"`
		}{Key: "key3", KeyInfo: "k3info"}, 3)
	d, err = json.Marshal(m2)
	a.Error(err)
}

func TestOrderedMap_Example(t *testing.T) {
	om := New[string, int]()
	om.Store("k1", 1).Store("k2", 2).Store("k3", 3).
		Store("k4", 4).Store("k5", 5)
	om.Load("k5")                // return 5, true
	om.LoadOrStore("k5", 55)     // return 5, true
	om.LoadOrStore("k6", 6)      // return 6, false
	om.Delete("k6").Delete("k7") // 'k6' be removed, Deleting a non-existent key will not report an error.
	om.Has("k6")                 // return false

	// use filter func to filter items
	filter1 := func(idx int, k string, v int) (want bool) { return v > 1 }
	filter2 := func(idx int, k string, v int) (want bool) { return v < 5 }
	filter3 := func(idx int, k string, v int) (want bool) { return v%2 == 0 }
	s0 := om.Slice()
	fmt.Println(s0) // out: [1 2 3 4 5]
	s1 := om.Slice(filter1, filter2, filter3)
	fmt.Println(s1) // out: [2 4]

	// use a filter to filter the key value when travel items
	om.TravelForward(func(idx int, k string, v int) (skip bool) {
		fmt.Printf("[NOFILTER] idx: %v, key: %v, val: %v\n", idx, k, v)
		return false
	})
	om.TravelForward(func(idx int, k string, v int) (skip bool) {
		fmt.Printf("[FILTER] idx: %v, key: %v, val: %v\n", idx, k, v)
		return false
	}, filter3)

	// JSON Marshal
	// output: {"k1":1,"k2":2,"k3":3,"k4":4,"k5":5}
	jBytes, _ := json.Marshal(om)
	fmt.Println(string(jBytes))
}
