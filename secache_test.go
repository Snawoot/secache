package secache

import (
	"math"
	"sync"
	"testing"

	"github.com/Snawoot/secache/randmap"
)

func TestNew(t *testing.T) {
	f := func(k int, v int) bool { return true }
	c := New(2, f)
	if c.n != 2 {
		t.Errorf("expected n=2, got %d", c.n)
	}

	c = New(1, f)
	if c.n != 2 {
		t.Errorf("expected min n=2, got %d", c.n)
	}
}

func TestLen(t *testing.T) {
	f := func(k int, v int) bool { return true }
	c := New(2, f)
	if c.Len() != 0 {
		t.Error("expected empty cache")
	}
}

func TestSetGet(t *testing.T) {
	f := func(k int, v int) bool { return true }
	c := New(2, f)
	c.Set(1, 10)
	v, ok := c.Get(1)
	if !ok || v != 10 {
		t.Error("expected to get set value")
	}
	if c.Len() != 1 {
		t.Error("expected len=1")
	}
}

func TestFlush(t *testing.T) {
	f := func(k int, v int) bool { return true }
	c := New(2, f)
	c.Set(1, 10)
	c.Flush()
	if c.Len() != 0 {
		t.Error("expected empty after flush")
	}
}

func TestGetValidOrDelete(t *testing.T) {
	forceValid := true
	f := func(k int, v int) bool { return v > 0 || forceValid }
	c := New(2, f)
	c.Set(1, 10)
	v, ok := c.GetValidOrDelete(1)
	if !ok || v != 10 {
		t.Error("expected valid value")
	}

	c.Set(2, -5)
	forceValid = false
	v, ok = c.Get(2)
	if !ok || v != -5 {
		t.Errorf("expected Get to return even invalid: v = %d; ok = %t", v, ok)
	}
	v, ok = c.GetValidOrDelete(2)
	if ok {
		t.Error("expected not ok for invalid")
	}
	_, ok = c.Get(2)
	if ok {
		t.Error("expected deleted after GetValidOrDelete")
	}
}

func TestGetOrCreate(t *testing.T) {
	f := func(k int, v int) bool { return v > 0 }
	c := New(2, f)
	called := 0
	v := c.GetOrCreate(1, func() int {
		called++
		return 10
	})
	if v != 10 || called != 1 {
		t.Error("expected new value created")
	}

	v = c.GetOrCreate(1, func() int {
		called++
		return 20
	})
	if v != 10 || called != 1 {
		t.Error("expected existing valid value used")
	}

	// Make invalid
	c.Do(func(m *randmap.RandMap[int, int]) {
		m.Set(1, -1)
	})

	v = c.GetOrCreate(1, func() int {
		called++
		return 30
	})
	if v != 30 || called != 2 {
		t.Error("expected new value for invalid")
	}
}

func TestConcurrent(t *testing.T) {
	f := func(k int, v int) bool { return true }
	c := New(2, f)
	var wg sync.WaitGroup
	const num = 100
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(i, i*10)
		}(i)
	}
	wg.Wait()

	for i := 0; i < num; i++ {
		v, ok := c.Get(i)
		if !ok || v != i*10 {
			t.Errorf("expected %d, got %d", i*10, v)
		}
	}
	if c.Len() != num {
		t.Errorf("expected len=%d, got %d", num, c.Len())
	}
}

func TestEviction(t *testing.T) {
	type K = int
	type V = int // generation number
	currentGen := 0
	ttl := 60
	f := func(k K, elemGen V) bool {
		return currentGen-elemGen < ttl
	}
	const n = 10
	const genItems = 1000
	const generations = 1000
	c := New[K, V](n, f)

	for currentGen = range generations {
		for i := range genItems {
			c.Set(currentGen*genItems+i, currentGen)
		}
	}

	var total int
	var invalid int
	c.Do(func(m *randmap.RandMap[K, V]) {
		for k, v := range m.Range {
			total++
			if !f(k, v) {
				invalid++
			}
		}
	})
	fractionInv := float64(invalid) / float64(total)
	expected := 1.0 / float64(n)
	if math.Abs(fractionInv-expected) > 0.1 {
		t.Errorf("invalid fraction %f not close to expected %f", fractionInv, expected)
	}
	t.Logf("tested Cache(%d, |ttl = %d|) with %d generations, %d elements in each generation: len = %d, %.2f%% expired",
		n, ttl, generations, genItems, c.Len(), fractionInv*100)
}
