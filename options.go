package shardedsingleflight

import (
	"hash"
	"hash/crc64"
	"hash/fnv"
	"hash/maphash"
	"runtime"
)

//GroupOptions represents an option that can be provided to the ShardGroup constructor
type GroupOptions interface {
	apply(g *ShardedGroup)
}

var (
	_ GroupOptions = WithShardCount(0)
	_ GroupOptions = WithHashFunc(DefHashFunc)
)

//WithShardCount is an option to override the default shard count with an explicit count
type WithShardCount uint64

func (shardCount WithShardCount) apply(g *ShardedGroup) {
	g.shardc = uint64(shardCount)
}

//DefShardCount is the heuristically determined default shard count. It will vary
// by machine taking into consideration factors such as the number of cores present
var DefShardCount = nextPrime(uint64(runtime.NumCPU() * 7))

//NewHash is the type for constructors of hash.Hash64 instances
type NewHash func() hash.Hash64

//WithHashFunc is an option to override the default hash.Hash64 constructor used
// by a ShardGroup
type WithHashFunc NewHash

func (hashFunc WithHashFunc) apply(g *ShardedGroup) {
	g.newHash = NewHash(hashFunc)
}

//The following are the hash.Hash64 constructors provided by various Go stdlib packages
var (
	DefHashFunc = FNV64a
	FNV64a      = fnv.New64a
	FNV64       = fnv.New64
	CRCISO      = func() hash.Hash64 { return crc64.New(crc64.MakeTable(crc64.ISO)) }
	CRCECMA     = func() hash.Hash64 { return crc64.New(crc64.MakeTable(crc64.ECMA)) }
	MapHash     = func() hash.Hash64 { return new(maphash.Hash) }
)
