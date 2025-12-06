// Package randmap implements key-value map capable of uniform sampling
// of keys in it.
package randmap

// RandMap represents key-value map which allows to fetch random key from map,
// granting equal probability of selection among keys. Time complexity of all
// operations is O(1) am.
type RandMap[K comparable, V any] struct {
	kv map[K]V
}

// Make creates empty map.
func Make[K comparable, V any]() *RandMap[K, V] {
	return &RandMap[K, V]{
		kv: make(map[K]V),
	}
}

// Wrap indexes and wraps existing standard map into a *RandMap instance.
// Original map should not be modified directly after that.
func Wrap[K comparable, V any](m map[K]V) *RandMap[K, V] {
	rm := &RandMap[K, V]{
		kv: m,
	}
	return rm
}

// Get retrieves key from map.
func (m *RandMap[K, V]) Get(key K) (val V, ok bool) {
	item, ok := m.kv[key]
	return item, ok
}

// Set adds or updates key-value pair in map.
func (m *RandMap[K, V]) Set(key K, item V) {
	m.kv[key] = item
}

// Delete removes key from map.
func (m *RandMap[K, V]) Delete(key K) {
	delete(m.kv, key)
}

// GetRandom retrieves uniformly-distributed random key-value pair from map,
// if it's not empty.
func (m *RandMap[K, V]) GetRandom() (k K, v V, ok bool) {
	for k, v = range m.kv {
		ok = true
		return
	}
	return
}

// Len returns number of key-value pairs in map.
func (m *RandMap[K, V]) Len() int {
	return len(m.kv)
}

// Range iterates over all map elements.
func (m *RandMap[K, V]) Range(f func(key K, value V) bool) {
	for k, v := range m.kv {
		if !f(k, v) {
			return
		}
	}
}
