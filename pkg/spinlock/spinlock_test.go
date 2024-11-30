package spinlock

import (
	"sync"
	"testing"

	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func TestLock(t *testing.T) {
	t.Run("count", func(t *testing.T) {
		const threadCount = 8
		wg := &sync.WaitGroup{}
		wg.Add(threadCount)

		l := &Lock{}
		count := 0
		const n = 1_000_000

		for range threadCount {
			go func() {
				for range n {
					l.Lock()
					count++
					l.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait()

		require.Equal(t, n*threadCount, count)
	})
}
