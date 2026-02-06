package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const benchRecordCount = 10_000

var storageProviders = []struct {
	name    string
	factory func(b *testing.B) Storage
}{
	{
		name: "memory",
		factory: func(b *testing.B) Storage {
			b.Helper()

			return NewMemStore()
		},
	},
	{
		name: "spill",
		factory: func(b *testing.B) Storage {
			b.Helper()

			config := DefaultSpillStoreConfig()
			config.MaxHotBytes = 100 * 1024
			config.CompactionThreshold = 20 * 1024
			s, err := NewSpillStore(config)
			require.NoError(b, err)
			return s
		},
	},
}

var benchmarkValue []byte

func init() {
	benchmarkValue = make([]byte, 2048)
	for i := range benchmarkValue {
		benchmarkValue[i] = byte(i)
	}
}

func BenchmarkStorage_AddGetRemove(b *testing.B) {
	for _, sp := range storageProviders {
		b.Run(sp.name, func(b *testing.B) {
			storage := sp.factory(b)
			b.Cleanup(func() { _ = storage.Close() })

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < benchRecordCount; j++ {
					_ = storage.Set(fmt.Sprintf("key-%d", j), benchmarkValue)
				}
				for j := 0; j < benchRecordCount; j++ {
					_, _, _ = storage.Get(fmt.Sprintf("key-%d", j))
				}
				for j := 0; j < benchRecordCount/2; j++ {
					storage.Delete(fmt.Sprintf("key-%d", j))
				}
				storage.DeleteAll()
			}
		})
	}
}

func BenchmarkStorage_ConcurrentAddGetRemove(b *testing.B) {
	for _, sp := range storageProviders {
		b.Run(sp.name, func(b *testing.B) {
			storage := sp.factory(b)
			b.Cleanup(func() { _ = storage.Close() })

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var wg sync.WaitGroup
				wg.Add(3)
				go func() {
					defer wg.Done()

					for j := 0; j < benchRecordCount; j++ {
						switch j {
						case benchRecordCount / 2: // halfway point start loading
							go func() {
								defer wg.Done()

								for j := 0; j < benchRecordCount; j++ {
									_, _, _ = storage.Get(fmt.Sprintf("key-%d", j))
								}
							}()
						case benchRecordCount * .75: // 3/4 point start deleting
							go func() {
								defer wg.Done()

								for j := 0; j < benchRecordCount/2; j++ {
									storage.Delete(fmt.Sprintf("key-%d", j))
								}
							}()
						}

						_ = storage.Set(fmt.Sprintf("key-%d", j), benchmarkValue)
					}
				}()
				wg.Wait()
				storage.DeleteAll()
			}
		})
	}
}
