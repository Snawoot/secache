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
