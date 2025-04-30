package syncx

import (
	"sync/atomic"
	"testing"

	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func TestPool(t *testing.T) {
	t.Run("GetPut", func(t *testing.T) {
		var count int32
		p := NewPool(func() int {
			return int(atomic.AddInt32(&count, 1))
		})

		v1 := p.Get()
		require.Equal(t, 1, v1)

		p.Put(v1)

		v2 := p.Get()
		require.Equal(t, 1, v2)
		require.Equal(t, int32(1), atomic.LoadInt32(&count))

		v3 := p.Get()
		require.Equal(t, 2, v3)
	})

	t.Run("Concurrency", func(t *testing.T) {
		var count int32
		p := NewPool(func() int {
			return int(atomic.AddInt32(&count, 1))
		})

		numGoroutines := 100
		numOps := 1000
		results := make(chan int, numGoroutines*numOps)

		for range numGoroutines {
			go func() {
				for range numOps {
					v := p.Get()
					results <- v
					p.Put(v)
				}
			}()
		}

		// Collect results - we don't check specific values due to concurrency
		// but we ensure all operations completed
		for range numGoroutines * numOps {
			<-results
		}

		// Check the maximum value created by New - it should be less than total ops
		// because items are reused.
		maxCreated := atomic.LoadInt32(&count)
		require.True(t, maxCreated > 0)
		t.Logf("Max items created: %d for %d ops by %d goroutines", maxCreated, numOps*numGoroutines, numGoroutines)
		// In practice, maxCreated will be much lower than numGoroutines * numOps
		// but we check it's less than total operations as a basic sanity check.
		require.True(t, maxCreated < int32(numGoroutines*numOps))
	})
}
