# orderedmap
 🧑‍💻 Implementation of ordered map in golang. Fast, thread-safe and generic support

# Install
- go version >= 1.18
```bash
go get -u github.com/LPX3F8/orderedmap
```

# Features
- Support conversion to slices
- Support JSON marshaler
- Support ordered traversal
- Support filter
- Thread safety
- Generics support

# Example
```go
import "github.com/LPX3F8/orderedmap"

func main() {
	om := orderedmap.New[string, int]()
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
		fmt.Printf("idx: %v, key: %v, val: %v\n", idx, k, v)
		return false
	}, filter3)

	// JSON Marshal
	// output: {"k1":1,"k2":2,"k3":3,"k4":4,"k5":5}
	jBytes, _ := json.Marshal(om)
	fmt.Println(string(jBytes))
```

# Benchmark
```text
# orderedmap basic test
BenchmarkOrderedMap-10                	 3441459	       340.9 ns/op	      32 B/op	       1 allocs/op
BenchmarkOrderedMapSlack-10           	 3154348	       379.2 ns/op	      32 B/op	       1 allocs/op
BenchmarkOrderedMapWork-10            	 3171066	       380.2 ns/op	      32 B/op	       1 allocs/op
BenchmarkOrderedMapWorkSlack-10       	 2857524	       421.1 ns/op	      32 B/op	       1 allocs/op

# golang sync.Map
BenchmarkNativeSyncMap_Store-10             	 1813066	       764.5 ns/op	     192 B/op	       5 allocs/op
BenchmarkNativeSyncMap_LoadOrStore-10       	 1691062	       686.5 ns/op	     117 B/op	       4 allocs/op
BenchmarkNativeSyncMap_Delete-10            	 1000000            2913 ns/op	       0 B/op	       0 allocs/op

# orderedmap
BenchmarkOrderedMap_Store-10          	 2492340	       512.3 ns/op	     171 B/op	       1 allocs/op
BenchmarkOrderedMap_LoadOrStore-10    	 2141635	       576.0 ns/op	     194 B/op	       1 allocs/op
BenchmarkOrderedMap_Delete-10         	 6182010	       174.1 ns/op	       0 B/op	       0 allocs/op
```