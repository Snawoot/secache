package randmap_test

import (
	"fmt"

	"github.com/Snawoot/secache/randmap"
)

func ExampleRandMap_Range() {
	m := randmap.Wrap(map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	})
	for k, v := range m.Range {
		fmt.Println(k, v)
	}
	// Output:
	// a 1
	// b 2
	// c 3
}
