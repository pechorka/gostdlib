package syncx_test

import (
	"maps"
	"testing"

	"github.com/pechorka/gostdlib/pkg/syncx"
	"github.com/pechorka/gostdlib/pkg/testing/require"
)

func TestMap_StoreLoad(t *testing.T) {
	t.Run("store and load basic types", func(t *testing.T) {
		m := syncx.NewMap[string, int]()

		m.Store("one", 1)
		m.Store("two", 2)

		v, ok := m.Load("one")
		require.Equal(t, true, ok)
		require.Equal(t, 1, v)

		v, ok = m.Load("two")
		require.Equal(t, true, ok)
		require.Equal(t, 2, v)

		v, ok = m.Load("three")
		require.Equal(t, false, ok)
		require.Equal(t, 0, v) // Check zero value for int
	})

	t.Run("store and load struct types", func(t *testing.T) {
		type customValue struct {
			ID   int
			Name string
		}
		m := syncx.NewMap[int, customValue]()

		val1 := customValue{ID: 1, Name: "Alice"}
		val2 := customValue{ID: 2, Name: "Bob"}

		m.Store(1, val1)
		m.Store(2, val2)

		v, ok := m.Load(1)
		require.Equal(t, true, ok)
		require.Equal(t, val1, v)

		v, ok = m.Load(2)
		require.Equal(t, true, ok)
		require.Equal(t, val2, v)

		v, ok = m.Load(3)
		require.Equal(t, false, ok)
		require.Equal(t, customValue{}, v) // Check zero value for struct
	})

	t.Run("overwrite existing key", func(t *testing.T) {
		m := syncx.NewMap[string, string]()
		m.Store("key", "initial")

		v, ok := m.Load("key")
		require.Equal(t, true, ok)
		require.Equal(t, "initial", v)

		m.Store("key", "updated")

		v, ok = m.Load("key")
		require.Equal(t, true, ok)
		require.Equal(t, "updated", v)
	})
}

func TestMap_LoadOrStore(t *testing.T) {
	m := syncx.NewMap[string, int]()

	t.Run("store new key", func(t *testing.T) {
		actual, loaded := m.LoadOrStore("one", 1)
		require.Equal(t, false, loaded)
		require.Equal(t, 1, actual)

		v, ok := m.Load("one")
		require.Equal(t, true, ok)
		require.Equal(t, 1, v)
	})

	t.Run("load existing key", func(t *testing.T) {
		// Ensure key "one" exists from previous subtest or add it
		m.Store("one", 1)

		actual, loaded := m.LoadOrStore("one", 100) // Try to store a different value
		require.Equal(t, true, loaded)
		require.Equal(t, 1, actual) // Should return the existing value (1)

		// Verify the value wasn't overwritten
		v, ok := m.Load("one")
		require.Equal(t, true, ok)
		require.Equal(t, 1, v)
	})
}

func TestMap_LoadAndDelete(t *testing.T) {
	m := syncx.NewMap[string, int]()
	m.Store("delete_me", 10)

	m.Store("keep_me", 20)

	t.Run("delete existing key", func(t *testing.T) {
		v, loaded := m.LoadAndDelete("delete_me")
		require.Equal(t, true, loaded)
		require.Equal(t, 10, v)

		// Verify it's deleted
		_, ok := m.Load("delete_me")
		require.Equal(t, false, ok)
	})

	t.Run("delete non-existent key", func(t *testing.T) {
		v, loaded := m.LoadAndDelete("non_existent")
		require.Equal(t, false, loaded)
		require.Equal(t, 0, v) // Should return zero value
	})

	t.Run("ensure other keys remain", func(t *testing.T) {
		v, ok := m.Load("keep_me")
		require.Equal(t, true, ok)
		require.Equal(t, 20, v)
	})
}

func TestMap_Delete(t *testing.T) {
	m := syncx.NewMap[string, int]()
	m.Store("delete_me", 10)
	m.Store("keep_me", 20)

	t.Run("delete existing key", func(t *testing.T) {
		m.Delete("delete_me")

		// Verify it's deleted
		_, ok := m.Load("delete_me")
		require.Equal(t, false, ok)
	})

	t.Run("delete non-existent key", func(t *testing.T) {
		// Deleting a non-existent key should be a no-op
		m.Delete("non_existent")
		_, ok := m.Load("non_existent")
		require.Equal(t, false, ok)
	})

	t.Run("ensure other keys remain", func(t *testing.T) {
		v, ok := m.Load("keep_me")
		require.Equal(t, true, ok)
		require.Equal(t, 20, v)
	})
}

func TestMap_Swap(t *testing.T) {
	m := syncx.NewMap[string, int]()
	m.Store("key1", 1)

	t.Run("swap existing key", func(t *testing.T) {
		previous, loaded := m.Swap("key1", 10)
		require.Equal(t, true, loaded)
		require.Equal(t, 1, previous)

		// Verify the new value
		v, ok := m.Load("key1")
		require.Equal(t, true, ok)
		require.Equal(t, 10, v)
	})

	t.Run("swap non-existent key", func(t *testing.T) {
		previous, loaded := m.Swap("key2", 20)
		require.Equal(t, false, loaded)
		require.Equal(t, 0, previous) // Should return zero value

		// Verify the new value was stored
		v, ok := m.Load("key2")
		require.Equal(t, true, ok)
		require.Equal(t, 20, v)
	})

	t.Run("swap again on previously swapped key", func(t *testing.T) {
		previous, loaded := m.Swap("key1", 100)
		require.Equal(t, true, loaded)
		require.Equal(t, 10, previous) // Previous value was 10

		// Verify the new value
		v, ok := m.Load("key1")
		require.Equal(t, true, ok)
		require.Equal(t, 100, v)
	})
}

func TestMap_CompareAndSwap(t *testing.T) {
	m := syncx.NewMap[string, int]()
	m.Store("cas_key", 5)

	t.Run("successful CAS", func(t *testing.T) {
		swapped := m.CompareAndSwap("cas_key", 5, 50)
		require.Equal(t, true, swapped)

		// Verify the new value
		v, ok := m.Load("cas_key")
		require.Equal(t, true, ok)
		require.Equal(t, 50, v)
	})

	t.Run("failed CAS with wrong old value", func(t *testing.T) {
		// Current value is 50
		swapped := m.CompareAndSwap("cas_key", 5, 500) // old value is 5, should fail
		require.Equal(t, false, swapped)

		// Verify value didn't change
		v, ok := m.Load("cas_key")
		require.Equal(t, true, ok)
		require.Equal(t, 50, v)
	})

	t.Run("failed CAS on non-existent key", func(t *testing.T) {
		swapped := m.CompareAndSwap("non_existent", 0, 100)
		require.Equal(t, false, swapped)

		// Verify key wasn't created
		_, ok := m.Load("non_existent")
		require.Equal(t, false, ok)
	})
}

func TestMap_CompareAndDelete(t *testing.T) {
	m := syncx.NewMap[string, int]()
	m.Store("cad_key", 15)
	m.Store("keep_key", 25)

	t.Run("successful CAD", func(t *testing.T) {
		deleted := m.CompareAndDelete("cad_key", 15)
		require.Equal(t, true, deleted)

		// Verify key is deleted
		_, ok := m.Load("cad_key")
		require.Equal(t, false, ok)
	})

	t.Run("failed CAD with wrong old value", func(t *testing.T) {
		// Store key again for this test
		m.Store("cad_key", 15)
		deleted := m.CompareAndDelete("cad_key", 100) // old value is 100, should fail
		require.Equal(t, false, deleted)

		// Verify key was not deleted
		v, ok := m.Load("cad_key")
		require.Equal(t, true, ok)
		require.Equal(t, 15, v)

		// Clean up for next tests if needed (or rely on test isolation)
		m.Delete("cad_key")
	})

	t.Run("failed CAD on non-existent key", func(t *testing.T) {
		deleted := m.CompareAndDelete("non_existent", 0)
		require.Equal(t, false, deleted)
	})

	t.Run("ensure other keys remain", func(t *testing.T) {
		v, ok := m.Load("keep_key")
		require.Equal(t, true, ok)
		require.Equal(t, 25, v)
	})
}

func TestMap_Range(t *testing.T) {
	m := syncx.NewMap[string, int]()
	items := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
	}
	for k, v := range items {
		m.Store(k, v)
	}

	t.Run("range over all items", func(t *testing.T) {
		count := 0
		seen := make(map[string]int)

		m.Range(func(key string, value int) bool {
			seen[key] = value
			count++
			return true // Continue ranging
		})

		require.Equal(t, len(items), count)
		require.Equal(t, len(items), len(seen))
		for k, vExpected := range items {
			vActual, ok := seen[k]
			require.Equal(t, true, ok)
			require.Equal(t, vExpected, vActual)
		}
	})

	t.Run("range and stop early", func(t *testing.T) {
		count := 0
		m.Range(func(key string, value int) bool {
			count++
			return false
		})

		require.True(t, count == 1)
	})

	t.Run("range on empty map", func(t *testing.T) {
		emptyMap := syncx.NewMap[int, bool]()
		called := false
		emptyMap.Range(func(key int, value bool) bool {
			called = true
			return false // Should not matter
		})
		require.Equal(t, false, called)
	})

	t.Run("range with type conversion check", func(t *testing.T) {
		type myKey struct{ id int }
		type myVal struct{ name string }
		mTyped := syncx.NewMap[myKey, myVal]()
		k1 := myKey{id: 1}
		v1 := myVal{name: "A"}
		k2 := myKey{id: 2}
		v2 := myVal{name: "B"}
		mTyped.Store(k1, v1)
		mTyped.Store(k2, v2)

		seen := make(map[myKey]myVal)
		mTyped.Range(func(key myKey, value myVal) bool {
			seen[key] = value
			return true
		})

		expected := map[myKey]myVal{k1: v1, k2: v2}
		require.True(t, maps.Equal(expected, seen))
	})
}
