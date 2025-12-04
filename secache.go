// Package secache implements generic cache with arbitrary expiration criteria
// defined by arbitrary validity function.
package secache

import (
	"sync"

	"github.com/Snawoot/secache/randmap"
)

// ValidityFunc is a function which used to check if element should stay in cache.
type ValidityFunc[K comparable, V any] = func(K, V) bool

type Cache[K comparable, V any] struct {
	mux sync.Mutex
	m   *randmap.RandMap[K, V]
	f   ValidityFunc[K, V]
	n   int
}

const MinN = 2

func New[K comparable, V any](n int, f ValidityFunc[K, V]) *Cache[K, V] {
	return &Cache[K, V]{
		m: randmap.Make[K, V](),
		n: max(n, MinN),
		f: f,
	}
}

// Flush empties cache.
func (c *Cache[K, V]) Flush() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.m = randmap.Make[K, V]()
}

// Do acquires lock and exposes storage to a provided function f.
// f should not operate on cache object, but only on exposed storage.
// Provided storage reference is valid only within f.
func (c *Cache[K, V]) Do(f func(*randmap.RandMap[K, V])) {
	c.mux.Lock()
	defer c.mux.Unlock()
	f(c.m)
}

// Len returns number of items in cache.
func (c *Cache[K, V]) Len() (l int) {
	c.Do(func(m *randmap.RandMap[K, V]) {
		l = m.Len()
	})
	return
}

// Get fetches key from cache.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.Do(func(m *randmap.RandMap[K, V]) {
		value, ok = m.Get(key)
	})
	return
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.Do(func(m *randmap.RandMap[K, V]) {
		oldLen := m.Len()
		m.Set(key, value)
		if newLen := m.Len(); newLen > oldLen {
			// new element was added, run eviction attempts
			for i := 0; i < c.n; i++ {
				ck, cv, ok := m.GetRandom()
				if !ok {
					// cache is empty
					break
				}
				if !c.f(ck, cv) {
					m.Delete(ck)
				}
			}
		}
	})
}
