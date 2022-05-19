package redis

import (
	"sync"
	"testing"
)

const numberDBs = 16
const numberConnectionsPerDB = 10

// Redis running in Docker locally
// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkSet-12    	   47932	     25566 ns/op	      88 B/op	       4 allocs/op
// over the network running Redis on LXC with 7 cores, 10 GB RAM, 2 GB Swap:
// BenchmarkSet-12    	    1639	    767060 ns/op	      88 B/op	       4 allocs/op
func BenchmarkSet(b *testing.B) {
	kv := KV{
		key:   "1",
		value: "x",
	}

	pool, _ := NewPool(_sock)

	pool.deleteByDB()
	defer pool.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pool.Set(kv)
	}
}

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkSetMany-12    	     338	   7552079 ns/op	  579989 B/op	    2528 allocs/op
// aka: 47200 ns/op
// over the network:
// BenchmarkSetMany-12    	     145	  12219460 ns/op	  613630 B/op	    2631 allocs/op
// over the network, 1 Gb RAM:
// BenchmarkSetMany-12    	     132	  16249872 ns/op	  620893 B/op	    2654 allocs/op
// over the network, 1 Gb RAM, 4 cores:
// BenchmarkSetMany-12    	     129	  21844816 ns/op	  492221 B/op	    2255 allocs/op
func BenchmarkSetMany(b *testing.B) {
	kv := KV{
		key:   "1",
		value: "x",
	}

	pools := make([]*Pool, numberDBs)

	for i := 0; i < numberDBs; i++ {
		pool, _ := NewPool(_sock, WithDatabaseNumber(uint(i)))
		defer pool.Close()

		pools[i] = pool
	}

	defer pools[0].deleteALL()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(numberConnectionsPerDB * numberDBs)

		for _, pool := range pools {
			for j := 0; j < numberConnectionsPerDB; j++ {
				go func() {
					pool.Set(kv)
					wg.Done()
				}()
			}
		}

		wg.Wait()
	}
}

// Redis running in Docker locally
// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkGet-12    	   46378	     24412 ns/op	     112 B/op	       7 allocs/op
// over the network:
// BenchmarkGet-12    	    1618	    756434 ns/op	     112 B/op	       7 allocs/op
func BenchmarkGet(b *testing.B) {
	kv := KV{
		key:   "1",
		value: "x",
	}

	pool, _ := NewPool(_sock)

	pool.deleteByDB()
	defer pool.Close()

	pool.Set(kv)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pool.Get(kv.key)
	}
}

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkGetMany-12    	     326	   7933118 ns/op	  601847 B/op	    2451 allocs/op
// aka: 49582 ns/op
// over the network:
// BenchmarkGetMany-12    	     120	  13801612 ns/op	  638768 B/op	    2566 allocs/op
func BenchmarkGetMany(b *testing.B) {
	kv := KV{
		key:   "1",
		value: "x",
	}

	pools := make([]*Pool, numberDBs)

	for i := 0; i < numberDBs; i++ {
		pool, _ := NewPool(_sock, WithDatabaseNumber(uint(i)))
		defer pool.Close()

		pools[i] = pool
	}

	defer pools[0].deleteALL()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(numberConnectionsPerDB * numberDBs)

		for _, pool := range pools {
			for j := 0; j < numberConnectionsPerDB; j++ {
				go func() {
					pool.Get(kv.key)
					wg.Done()
				}()
			}
		}

		wg.Wait()
	}
}
