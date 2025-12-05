// Package secache implements sampling eviction cache, a generic cache with
// arbitrary expiration criteria defined by validity function provided by user.
// It offers O(1) time complexity for all operations and very flexible notion
// of element validity, which is useful when usual approaches based on item age
// do not provide reasonable approximation or fit at all.
package secache

import (
	"sync"

	"github.com/Snawoot/secache/randmap"
)

// ValidityFunc is a function which used to check if element should stay in cache.
type ValidityFunc[K comparable, V any] = func(K, V) bool

// Cache implements a generic multipurpose cache which uses sampling eviction to
// maintain stable ratio of evictable and valid elements.
//
// Each time new element is added, cache makes n attempts to pick random cache
// element, test its validity and remove invalid one. This way it maintains
// dynamic equilibrium around certain rate of evictable elements, trading space
// for eviction efforts and flexibility. In many cases, however, space saved by
// more accurate custom eviction criteria may make up for space overhead
// compared to classical TTL cache which approximate object validity with
// its age.
//
// All cache operations have O(1) am. time complexity.
//
// Cache object is safe for concurrent use by multiple goroutines.
type Cache[K comparable, V any] struct {
	mux sync.Mutex
	m   *randmap.RandMap[K, V]
	f   ValidityFunc[K, V]
	n   int
}

// MinN is the minimal number of sampling evictions per element addition to
// ensure stable cache overhead.
const MinN = 2

// New creates new cache instance with n sampling eviction attempts per
// element addition. Validity of sampled elements is tested with function f.
//
// In practical terms, n corresponds to cache overhead goal, how many invalid
// elements on average will be there in a long run.
//
// Useful n values are:
//
//	2  - ~50% of invalid elements,		~100% overhead
//	3  - ~33.(3)% of invalid elements,	~50% overhead
//	4  - ~25% of invalid elements,		~33.(3)% overhead
//	5  - ~20% of invalid elements,		~25% overhead
//	6  - ~16.(6)% of invalid elements,	~20% overhead
//	7  - ~14.28% of invalid elements,	~16.(6)% overhead
//	8  - ~12.5% of invalid elements,	~14.28% overhead
//	9  - ~11.(1)% of invalid elements,	~12.5% overhead
//	10 - ~10% of invalid elements,		~11.(1)% overhead
//	11 - ~9.(09)% of invalid elements,	~10% overhead
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
// f should not operate on cache object, but only on provided storage.
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

// Get lookups key in cache, valid or not.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.Do(func(m *randmap.RandMap[K, V]) {
		value, ok = m.Get(key)
	})
	return
}

// GetValidOrDelete fetches valid key from cache or deletes it if it was
// found, but not valid.
func (c *Cache[K, V]) GetValidOrDelete(key K) (value V, ok bool) {
	c.Do(func(m *randmap.RandMap[K, V]) {
		value, ok = m.Get(key)
		if !ok {
			return
		}
		if !c.f(key, value) {
			ok = false
			m.Delete(key)
		}
	})
	return
}

// GetOrCreate fetches valid key from cache or creates new one with provided function.
func (c *Cache[K, V]) GetOrCreate(key K, newValFunc func() V) (value V) {
	c.Do(func(m *randmap.RandMap[K, V]) {
		var ok bool
		value, ok = m.Get(key)
		if !ok || !c.f(key, value) {
			value = newValFunc()
			c.SetLocked(m, key, value)
		}
	})
	return
}

// SetLocked is an utility function which adds or updates key with proper
// expiration logic. It is intended to be used within Do(f) transaction.
func (c *Cache[K, V]) SetLocked(m *randmap.RandMap[K, V], key K, value V) {
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
}

// Set adds new item to cache or updates existing one and then runs
// sampling eviction if new item was added.
func (c *Cache[K, V]) Set(key K, value V) {
	c.Do(func(m *randmap.RandMap[K, V]) {
		c.SetLocked(m, key, value)
	})
}
