package shardedsingleflight

import (
	"golang.org/x/sync/singleflight"
)

//ShardedGroup is a sharded singleflight Group. See singleflight.Group
type ShardedGroup struct {
	newHash NewHash
	shardc  uint64
	shards  []singleflight.Group
}

//NewShardedGroup is the sharded singleflight Group contstructor. Any provided
//options, which may be absent, will be applied.
func NewShardedGroup(opts ...GroupOptions) *ShardedGroup {
	sg := &ShardedGroup{
		newHash: DefHashFunc,
		shardc:  DefShardCount,
	}

	for _, opt := range opts {
		opt.apply(sg)
	}
	sg.shards = make([]singleflight.Group, sg.shardc)

	return sg
}

//Do maps the key to a shard and runs the singleflight shard's Do.
// See singleflight.Do
func (sg *ShardedGroup) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	return sg.shards[sg.shardIdx(key)].Do(key, fn)
}

//DoChan maps the key to a shard and runs the singleflight shard's DoChan.
// See singleflight.DoChan
func (sg *ShardedGroup) DoChan(key string, fn func() (interface{}, error)) <-chan singleflight.Result {
	return sg.shards[sg.shardIdx(key)].DoChan(key, fn)
}

//Forget maps the key to a shard and runs the singleflight shard's Forget.
// See singleflight.Forget
func (sg *ShardedGroup) Forget(key string) {
	sg.shards[sg.shardIdx(key)].Forget(key)
}

func (sg *ShardedGroup) shardIdx(key string) uint64 {
	hasher := sg.newHash()
	hasher.Write([]byte(key))
	return hasher.Sum64() % sg.shardc
}
