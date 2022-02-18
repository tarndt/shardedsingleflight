# shardedsingleflight
 [![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)[![Go Reference](https://pkg.go.dev/badge/github.com/tarndt/shardedsingleflight.svg)](https://pkg.go.dev/github.com/tarndt/shardedsingleflight) [![Go Report Card](https://goreportcard.com/badge/github.com/tarndt/shardedsingleflight)](https://goreportcard.com/report/github.com/tarndt/shardedsingleflight)

A sharded wrapper for `golang.org/x/sync/singleflight` ([code](https://github.com/golang/sync/tree/master/singleflight), [docs](https://pkg.go.dev/golang.org/x/sync/singleflight)) for high contention/concurrency environments.

## What is singleflight?
If you are not familiar, *singleflight* is a package created by [Brad Fitzpatrick](https://en.wikipedia.org/wiki/Brad_Fitzpatrick) that addresses the [thundering herd problem](https://en.wikipedia.org/wiki/Thundering_herd_problem) by assigning every operation a key and de-duplicating concurrently invoked operations based on that key. So for example, if you have a function that reads a file from disk and you wrap that function with singleflight, if the function is invoked twice, the second caller will get the same result returned to the first caller and the file will only be read once.

```go
//Not thundering herd safe!
func readFile(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	return string(data), err
}

var g singeflight.Group //normally a field in a struct

//Safe from the herd!
func ReadFile(filepath string) (string, error) {
	result, err, _ := g.Do(filepath, func() (interface{}, error) {
		return readFile(filepath)
	})
	return result.(string), err
}
```
*See [singleflight.Do](https://pkg.go.dev/golang.org/x/sync/singleflight#Group.Do) and [singleflight.DoChan](https://pkg.go.dev/golang.org/x/sync/singleflight#Group.DoChan) for more details.*

When duplicate operations are expensive and results can be shared, *singeflight*'s reduction in duplicate computation and/or I/O almost always warrants the overhead; I have found that on systems running very concurrent workloads that the [mutex](https://pkg.go.dev/sync#Mutex) contention [internal to singleflight](https://github.com/golang/sync/blob/master/singleflight/singleflight.go#L69) can quickly become significant.

This package creates shard's that each have their own [singleflight.Group](https://pkg.go.dev/golang.org/x/sync/singleflight#Group) and use a ([hash function](https://pkg.go.dev/hash#Hash64)) to distribute the keys over them. The result is a drastic reduction of [mutex](https://pkg.go.dev/sync#Mutex) contention as demonstrated by the package's benchmarks that compare it to a vanilla *singleflight* Group. How drastic a reduction depends on factors including how well distributed your hash function is, the number of shards provisioned, how many cores your system has, and the level of concurrency.

## So- The short version: Why does this exist?
1. [Brad Fitzpatrick](https://en.wikipedia.org/wiki/Brad_Fitzpatrick)'s *singleflight* library is amazingly useful! It is a robust and simple way to counter the [thundering herd problem](https://en.wikipedia.org/wiki/Thundering_herd_problem) in many cases.
2. A number of times in [my career](https://www.linkedin.com/in/tylorarndt/) I have have encountered problems using *singleflight* on machines with many cores/goroutines due to contention for the internal mutexes used by *singleflight* Groups.
3. I have written less robust versions of the sharded solution in this repo too many times and would like to spend my time on more interesting problems in the future.
4. If you face a similar challenge, I hope you can benefit from this solution as well!

### Show me the money!
*shardedsingleflight* allows configuring both the shard count and shard mapping ([hash](https://pkg.go.dev/hash#Hash64)) algorithm to be overridden. Below is a comparison of parallel vanilla *singleflight* (`noshard-24`) vs. *shardedsingleflight* on a 24 logical-core machine using various hash algorithms and the default shard count heuristic (`nextPrime(logical-cores * 7)`). On this machine, *shardedsingleflight* using [FNV-64](https://pkg.go.dev/hash/fnv) is about **9x faster** than vanilla *singleflight*. As always test on your own hardware and using your own software to validate this is worth using over vanilla *singleflight*. Software engineering is always about tradeoffs, we are trading a little extra memory and hash computation for less [mutex](https://pkg.go.dev/sync#Mutex) contention internal to singleflight, but that only pays off with enough concurrency (and thus contention).
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