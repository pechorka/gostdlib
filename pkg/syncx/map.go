package syncx

import "sync"

// Map is a generic wrapper around sync.Map.
type Map[Key comparable, Val any] struct {
	data *sync.Map
}

// NewMap creates a new Map.
func NewMap[Key comparable, Val any]() *Map[Key, Val] {
	return &Map[Key, Val]{data: &sync.Map{}}
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *Map[Key, Val]) Load(key Key) (zero Val, ok bool) {
	v, ok := m.data.Load(key)
	if !ok {
		return zero, false
	}
	return v.(Val), true
}

// Store sets the value for a key.
func (m *Map[Key, Val]) Store(key Key, value Val) {
	m.data.Store(key, value)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map[Key, Val]) LoadOrStore(key Key, value Val) (Val, bool) {
	a, ok := m.data.LoadOrStore(key, value)
	return a.(Val), ok
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[Key, Val]) LoadAndDelete(key Key) (zero Val, ok bool) {
	v, ok := m.data.LoadAndDelete(key)
	if !ok {
		return zero, false
	}
	return v.(Val), true
}

// Delete deletes the value for a key.
func (m *Map[Key, Val]) Delete(key Key) {
	m.data.Delete(key)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[Key, Val]) Swap(key Key, value Val) (zero Val, ok bool) {
	p, ok := m.data.Swap(key, value)
	if !ok {
		return zero, false
	}
	return p.(Val), true
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (m *Map[Key, Val]) CompareAndSwap(key Key, old, new Val) (swapped bool) {
	// Note: This requires Val to be comparable for the 'old' parameter comparison.
	// The standard sync.Map CompareAndSwap does not have this constraint directly
	// on the type parameter but relies on the actual value comparison.
	// We will rely on the underlying sync.Map's behavior, which uses interface{} comparison.
	// The type assertion for 'old' might not be strictly necessary depending on usage,
	// but it maintains consistency with the generic type.
	// If Val is not comparable, this method might behave unexpectedly or panic
	// depending on the underlying implementation of interface comparison.
	return m.data.CompareAndSwap(key, old, new)
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
func (m *Map[Key, Val]) CompareAndDelete(key Key, old Val) bool {
	// Note: This requires Val to be comparable for the 'old' parameter comparison.
	// The standard sync.Map CompareAndDelete does not have this constraint directly
	// on the type parameter but relies on the actual value comparison.
	// We will rely on the underlying sync.Map's behavior, which uses interface{} comparison.
	// The type assertion for 'old' might not be strictly necessary depending on usage,
	// but it maintains consistency with the generic type.
	// If Val is not comparable, this method might behave unexpectedly or panic
	// depending on the underlying implementation of interface comparison.
	return m.data.CompareAndDelete(key, old)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Map[Key, Val]) Range(f func(key Key, value Val) bool) {
	m.data.Range(func(k, v any) bool {
		return f(k.(Key), v.(Val))
	})
}
