# orderedmap
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/LPX3F8/orderedmap)](https://goreportcard.com/report/github.com/LPX3F8/orderedmap)
![Coverage](https://img.shields.io/badge/Coverage-94.4%25-brightgreen)
[![](https://img.shields.io/badge/README-中文-yellow.svg)](https://github.com/LPX3F8/orderedmap/blob/main/README_CN.md)

 🧑‍💻 Implementation of ordered map in golang. Fast, thread-safe and generic support

## Install
- go version >= 1.18
```bash
go get -u github.com/LPX3F8/orderedmap
```

## Features
- Support conversion to slices
- Support JSON marshaler
- Support ordered traversal
- Support filter
- Thread safety
- Generics support

## Example
```go
import "github.com/LPX3F8/orderedmap"

func main() {
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

	// travel items
	for i := om.Front(); i != nil; i = i.Next() {
		fmt.Println("[TEST FRONT]", i.Key(), i.Value())
	}
	for i := om.Back(); i != nil; i = i.Prev() {
		fmt.Println("[TEST BACK]", i.Key(), i.Value())
	}
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
```

## Benchmark
```text
goos: darwin
goarch: arm64
pkg: github.com/LPX3F8/orderedmap

# Basic test
BenchmarkOrderedMap-10                   	 3498038	       338.5 ns/op	      64 B/op	       2 allocs/op
BenchmarkOrderedMapSlack-10              	 3410408	       352.6 ns/op	      64 B/op	       2 allocs/op
BenchmarkOrderedMapWork-10               	 3167127	       378.6 ns/op	      64 B/op	       2 allocs/op
BenchmarkOrderedMapWorkSlack-10          	 3039068	       394.3 ns/op	      64 B/op	       2 allocs/op

# Native Sync.Map test
BenchmarkNativeSyncMap_Store-10          	 1510597	       668.7 ns/op	     140 B/op	       5 allocs/op
BenchmarkNativeSyncMap_LoadOrStore-10    	 1749106	       689.8 ns/op	     181 B/op	       4 allocs/op
BenchmarkNativeSyncMap_Delete-10         	 1000000	      2203 ns/op	       0 B/op	       0 allocs/op

# OrderedMap test
BenchmarkOrderedMap_Store-10             	 3161652	       379.7 ns/op	     120 B/op	       2 allocs/op
BenchmarkOrderedMap_LoadOrStore-10       	 2854708	       421.1 ns/op	     125 B/op	       2 allocs/op
BenchmarkOrderedMap_Delete-10            	 8021584	       144.9 ns/op	       0 B/op	       0 allocs/op
```