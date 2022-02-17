package shardedsingleflight

import (
	"crypto/rand"
	"fmt"
	"io"
	"testing"

	"golang.org/x/sync/singleflight"
)

type doer interface {
	Do(string, func() (interface{}, error)) (interface{}, error, bool)
}

func BenchmarkDo(b *testing.B) {
	keys := randKeys(b, 1024, 10)

	b.Run("noshard", func(b *testing.B) {
		benchDo(b, new(singleflight.Group), keys)
	})

	hashes := []struct {
		constr NewHash
		name   string
	}{
		{FNV64, "fnv64"}, {FNV64a, "fnv64a"},
		{CRCISO, "crc-iso"}, {CRCECMA, "crc-ecma"},
		{MapHash, "maphash"},
	}

	for _, hash := range hashes {
		b.Run(fmt.Sprintf("shard-hash-"+hash.name), func(b *testing.B) {
			benchDo(b, NewShardedGroup(WithHashFunc(hash.constr)), keys)
		})
	}
}

func benchDo(b *testing.B, g doer, keys []string) {
	keyc := len(keys)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			g.Do(keys[i%keyc], func() (interface{}, error) {
				return nil, nil
			})
		}
	})
}

func randKeys(b *testing.B, count, length uint) []string {
	keys := make([]string, 0, count)
	key := make([]byte, length)

	for i := uint(0); i < count; i++ {
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			b.Fatalf("Failed to generate random key %d of %d of length %d: %s", i+1, count, length, err)
		}
		keys = append(keys, string(key))
	}
	return keys
}
