// Package randmap implements key-value map capable of uniform sampling
// of keys in it.
package randmap

import (
	"math/rand/v2"
)

// RandMap represents key-value map which allows to fetch random key from map,
// granting equal probability of selection among keys. Time complexity of all
// operations is O(1) am.
type RandMap[K comparable, V any] struct {
	kv map[K]V
	ik map[int]K
	ki map[K]int
}

// Make creates empty map.
func Make[K comparable, V any]() *RandMap[K, V] {
	return &RandMap[K, V]{
		kv: make(map[K]V),
		ik: make(map[int]K),
		ki: make(map[K]int),
	}
}

// Wrap indexes and wraps existing standard map into a *RandMap instance.
// Original map should not be modified directly after that.
func Wrap[K comparable, V any](m map[K]V) *RandMap[K, V] {
	rm := &RandMap[K, V]{
		kv: m,
		ik: make(map[int]K),
		ki: make(map[K]int),
	}
	i := 0
	for k := range m {
		rm.ik[i] = k
		rm.ki[k] = i
		i++
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
	oldLen := len(m.kv)
	m.kv[key] = item
	if newLen := len(m.kv); newLen > oldLen {
		// adding new element and updating indexes
		m.ik[newLen-1] = key
		m.ki[key] = newLen - 1
	}
}

// Delete removes key from map.
func (m *RandMap[K, V]) Delete(key K) {
	deletedIdx, ok := m.ki[key]
	if !ok {
		return
	}
	oldLen := m.Len()

	delete(m.kv, key)
	delete(m.ki, key)
	if deletedIdx == oldLen-1 {
		delete(m.ik, deletedIdx)
	} else {
		relocatedKey := m.ik[oldLen-1]
		m.ki[relocatedKey] = deletedIdx
		m.ik[deletedIdx] = relocatedKey
		delete(m.ik, oldLen-1)
	}
}

// GetRandom retrieves uniformly-distributed random key-value pair from map,
// if it's not empty.
func (m *RandMap[K, V]) GetRandom() (val V, ok bool) {
	var emptyV V
	l := m.Len()
	if l == 0 {
		return emptyV, false
	}
	item, ok := m.kv[m.ik[rand.IntN(l)]]
	return item, ok
}

// Len returns number of key-value pairs in map.
func (m *RandMap[K, V]) Len() int {
	return len(m.kv)
}

// Range iterates over all map elements.
//
// Example:
//
//	m := Make()
//	for k, v := range m.Range {
//		fmt.Println(k, v)
//	}
func (m *RandMap[K, V]) Range(f func(key K, value V) bool) {
	for k, v := range m.kv {
		if !f(k, v) {
			return
		}
	}
}
