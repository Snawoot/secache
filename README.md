# secache

Sampling Eviction Cache

[![Go Reference](https://pkg.go.dev/badge/github.com/Snawoot/secache.svg)](https://pkg.go.dev/github.com/Snawoot/secache)

## Features

* Policy-agnostic: item validity is determined by user-provided function. That allows:
  * Use of common expiration strategies such as TTL, LRU, ...
  * Use of validity criterias specific for particular use case, such as item internal state, item usage statistics and so on.
* No full cache sweeps, no background goroutines for cleanup: expiration is handled probabilistically with certain dirty ratio guarantee.
* O(1) amortized time complexity for all operations.
* Adjustable space overhead for tradeoff between time and space.
* Simple and clear implementation well within 200 LOC.
* Support for complex operations within single critical section.

## Example

```go
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Snawoot/secache"
)

func main() {
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
}
```
