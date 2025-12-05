package secache_test

import (
	"fmt"

	"github.com/Snawoot/secache"
	"github.com/Snawoot/secache/randmap"
)

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
