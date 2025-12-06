package secache_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/Snawoot/secache"
	"github.com/Snawoot/secache/randmap"
)

func Example() {
	// demonstrates use of cache as a usual TTL cache
	const TTL = 1 * time.Minute
	type CacheItem struct {
		expires time.Time
		value   string
	}
	c := secache.New[string, *CacheItem](3, func(key string, item *CacheItem) bool {
		return time.Now().Before(item.expires)
	})

	key := "some key"
	item, ok := c.GetValidOrDelete(key)
	fmt.Println(item)
	if !ok {
		c.Set(key, &CacheItem{
			expires: time.Now().Add(TTL),
			value:   strings.ToTitle(key),
		})
	}

	item, ok = c.GetValidOrDelete(key)
	fmt.Printf("%q %t", item.value, ok)
	// Output:
	// <nil>
	// "SOME KEY" true
}

func ExampleCache_Do() {
	// demonstrates cache item increment
	c := secache.New[string, int](2, func(_ string, _ int) bool {
		return true
	})
	c.Set("a", 1)
	c.Set("b", 2)
	incrKey := "c"
	c.Do(func(m *randmap.RandMap[string, int]) {
		old, _ := m.Get(incrKey)
		m.Set(incrKey, old+1)
	})
	val, _ := c.Get(incrKey)
	fmt.Printf("c[%q] = %d\n", incrKey, val)
	// Output:
	// c["c"] = 1
}
