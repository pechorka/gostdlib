package syncx

import (
	"context"
	"errors"
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

func TestErrorPool(t *testing.T) {
	t.Run("Get from empty pool calls fabric", func(t *testing.T) {
		var fabricCalled atomic.Bool
		var createdValue int = 123
		fabric := func(ctx context.Context) (int, error) {
			fabricCalled.Store(true)
			return createdValue, nil
		}
		p := NewErrorPool(fabric)

		val, err := p.Get(context.Background())
		require.NoError(t, err)
		require.Equal(t, createdValue, val)
		require.True(t, fabricCalled.Load())
	})

	t.Run("Get from non-empty pool does not call fabric", func(t *testing.T) {
		var fabricCalled atomic.Bool
		fabric := func(ctx context.Context) (int, error) {
			fabricCalled.Store(true)
			return 0, errors.New("should not be called")
		}
		p := NewErrorPool(fabric)
		putValue := 456
		p.Put(putValue)

		val, err := p.Get(context.Background())
		require.NoError(t, err)
		require.Equal(t, putValue, val)
		require.False(t, fabricCalled.Load())
	})

	t.Run("Put and Get sequence", func(t *testing.T) {
		var fabricCallCount atomic.Int32
		createdValue := 789
		fabric := func(ctx context.Context) (int, error) {
			fabricCallCount.Add(1)
			return createdValue, nil
		}
		p := NewErrorPool(fabric)
		putValue := 101

		// 1. Put an item
		p.Put(putValue)
		require.Equal(t, int32(0), fabricCallCount.Load())

		// 2. Get the put item
		val1, err1 := p.Get(context.Background())
		require.NoError(t, err1)
		require.Equal(t, putValue, val1)
		require.Equal(t, int32(0), fabricCallCount.Load()) // Fabric still not called

		// 3. Get again, pool is empty, fabric should be called
		val2, err2 := p.Get(context.Background())
		require.NoError(t, err2)
		require.Equal(t, createdValue, val2)
		require.Equal(t, int32(1), fabricCallCount.Load()) // Fabric called once
	})

	t.Run("Fabric returns error", func(t *testing.T) {
		var fabricCalled atomic.Bool
		fabricErr := errors.New("fabric error")
		fabric := func(ctx context.Context) (int, error) {
			fabricCalled.Store(true)
			return 0, fabricErr
		}
		p := NewErrorPool(fabric)

		val, err := p.Get(context.Background())
		require.ErrorIs(t, err, fabricErr)
		require.Equal(t, 0, val) // Check for zero value of T
		require.True(t, fabricCalled.Load())
	})

	t.Run("Fabric context propagates", func(t *testing.T) {
		type ctxKey string
		key := ctxKey("testKey")
		expectedValue := "testValue"
		var receivedValue any

		fabric := func(ctx context.Context) (int, error) {
			receivedValue = ctx.Value(key)
			return 1, nil
		}
		p := NewErrorPool(fabric)
		ctx := context.WithValue(context.Background(), key, expectedValue)

		_, err := p.Get(ctx)
		require.NoError(t, err)
		require.NotNil(t, receivedValue)
		receivedString, ok := receivedValue.(string)
		require.True(t, ok)
		require.Equal(t, expectedValue, receivedString)
	})
}
