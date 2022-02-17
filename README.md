# shardedsingleflight
 [![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)[![Go Reference](https://pkg.go.dev/badge/github.com/tarndt/shardedsingleflight.svg)](https://pkg.go.dev/github.com/tarndt/shardedsingleflight) [![Go Report Card](https://goreportcard.com/badge/github.com/tarndt/shardedsingleflight)](https://goreportcard.com/report/github.com/tarndt/shardedsingleflight)

A sharded version of `golang.org/x/sync/singleflight` ([code](https://github.com/golang/sync/tree/master/singleflight), [docs](https://pkg.go.dev/golang.org/x/sync/singleflight)) for high contention/concurrency environments.

## Why does this exist?

1. [Brad Fitzpatrick](https://en.wikipedia.org/wiki/Brad_Fitzpatrick)'s *singleflight* library is amazingly useful! It is a robust and simple way to counter the [thundering herd problem](https://en.wikipedia.org/wiki/Thundering_herd_problem) in many cases.
2. A number of times in [my career](https://www.linkedin.com/in/tylorarndt/) I have have encountered problems using *singleflight* on machines with many cores/goroutines due to contention for the internal mutexes used by *singleflight* Groups.
3. I have written less robust versions of the sharded solution in this repo too many times and would like to spend my time on more interesting problems in the future.
4. If you face a similar challenge, I hope you can benefit from this solution as well.

### Show me the money!

*shardedsingleflight* allows configuring both the shard count and shard mapping ([hash](https://pkg.go.dev/hash#Hash64)) algorithm to be specified. Below is a comparison of parallel vanilla *singleflight* (`noshard-24`) vs. *shardedsingleflight* on a 24 logical-core machine using various hash algorithms and the default shard count heuristic (`nextPrime(v-cores * 7)`). On this machine, *shardedsingleflight* using [FNV-64](https://pkg.go.dev/hash/fnv) is about **9x faster** than vanilla *singleflight*. As always test on your own hardware and using your own software to valid this is worth using over vanilla *singleflight*.
```
go test -test.bench=.*
goos: linux
goarch: amd64
pkg: github.com/tarndt/shardedsingleflight
cpu: AMD Ryzen 9 3900X 12-Core Processor
BenchmarkDo/noshard-24     			    	 2070744	       583.30 ns/op	      94 B/op	       0 allocs/op
BenchmarkDo/shard-hash-fnv64-24         	18465124	        65.65 ns/op	     119 B/op	       2 allocs/op
BenchmarkDo/shard-hash-fnv64a-24        	16838563	        69.95 ns/op	     119 B/op	       2 allocs/op
BenchmarkDo/shard-hash-crc-iso-24       	15008778	        77.21 ns/op	     127 B/op	       2 allocs/op
BenchmarkDo/shard-hash-crc-ecma-24      	14828358	        79.10 ns/op	     127 B/op	       2 allocs/op
BenchmarkDo/shard-hash-maphash-24       	 9129664	       133.1 ns/op	     272 B/op	       3 allocs/op
PASS
ok  	github.com/tarndt/shardedsingleflight	3.162s
```
*Note: As seen above the extra complexity involved in sharding comes with a memory allocation tradeoff, but still is favorable in terms of execution time in a high-contention environment nonetheless.*

### Contributing

Issues, PRs and feedback are welcome!