package randmap

import (
	"testing"
)

func TestMake(t *testing.T) {
	m := Make[int, string]()
	if m.Len() != 0 {
		t.Errorf("expected len 0, got %d", m.Len())
	}
	_, ok := m.GetRandom()
	if ok {
		t.Error("expected no value from empty map")
	}
}

func TestSetGet(t *testing.T) {
	m := Make[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	if m.Len() != 3 {
		t.Errorf("expected len 3, got %d", m.Len())
	}

	val, ok := m.Get("a")
	if !ok || val != 1 {
		t.Errorf("expected 1 for 'a', got %v, %v", val, ok)
	}

	val, ok = m.Get("d")
	if ok {
		t.Errorf("expected no value for 'd', got %v", val)
	}
}

func TestUpdate(t *testing.T) {
	m := Make[string, int]()
	m.Set("a", 1)
	if val, ok := m.Get("a"); !ok || val != 1 {
		t.Errorf("expected 1 for 'a', got %v, %v", val, ok)
	}

	m.Set("a", 10)
	if val, ok := m.Get("a"); !ok || val != 10 {
		t.Errorf("expected 10 for 'a', got %v, %v", val, ok)
	}
	if m.Len() != 1 {
		t.Errorf("expected len 1 after update, got %d", m.Len())
	}
}

func TestDelete(t *testing.T) {
	m := Make[int, string]()
	m.Set(1, "one")
	m.Set(2, "two")
	m.Set(3, "three")

	m.Delete(2)
	if m.Len() != 2 {
		t.Errorf("expected len 2, got %d", m.Len())
	}
	_, ok := m.Get(2)
	if ok {
		t.Error("expected '2' to be deleted")
	}

	// Check remaining
	if val, ok := m.Get(1); !ok || val != "one" {
		t.Errorf("expected 'one' for 1, got %v, %v", val, ok)
	}
	if val, ok := m.Get(3); !ok || val != "three" {
		t.Errorf("expected 'three' for 3, got %v, %v", val, ok)
	}
}

func TestDeleteNonexistent(t *testing.T) {
	m := Make[string, int]()
	m.Set("a", 1)
	m.Delete("b")
	if m.Len() != 1 {
		t.Errorf("expected len 1, got %d", m.Len())
	}
}

func TestDeleteLast(t *testing.T) {
	m := Make[int, string]()
	m.Set(1, "one")
	m.Set(2, "two")
	m.Delete(2) // Delete last
	if m.Len() != 1 {
		t.Errorf("expected len 1, got %d", m.Len())
	}
	if val, ok := m.Get(1); !ok || val != "one" {
		t.Errorf("expected 'one' for 1, got %v, %v", val, ok)
	}
}

func TestGetRandom(t *testing.T) {
	m := Make[string, int]()
	_, ok := m.GetRandom()
	if ok {
		t.Error("expected no value from empty map")
	}

	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// Call multiple times to check validity
	seen := make(map[int]struct{})
	for i := 0; i < 100; i++ {
		val, ok := m.GetRandom()
		if !ok {
			t.Error("expected value from non-empty map")
		}
		if val != 1 && val != 2 && val != 3 {
			t.Errorf("unexpected value %d", val)
		}
		seen[val] = struct{}{}
	}
	if len(seen) != 3 {
		t.Errorf("expected to see all 3 values in 100 calls, saw %d", len(seen))
	}
}

func TestGetRandomAfterDelete(t *testing.T) {
	m := Make[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.Delete("b")

	// Call multiple times
	seen := make(map[int]struct{})
	for i := 0; i < 100; i++ {
		val, ok := m.GetRandom()
		if !ok {
			t.Error("expected value")
		}
		if val != 1 && val != 3 {
			t.Errorf("unexpected value %d", val)
		}
		seen[val] = struct{}{}
	}
	if len(seen) != 2 {
		t.Errorf("expected to see 2 values, saw %d", len(seen))
	}
}

func TestInternalConsistency(t *testing.T) {
	m := Make[int, string]()
	for i := 1; i <= 10; i++ {
		m.Set(i, "val")
	}
	for i := 5; i <= 8; i++ {
		m.Delete(i)
	}
	if m.Len() != 6 {
		t.Errorf("expected len 6, got %d", m.Len())
	}

	// Check ik and ki consistency
	for idx, key := range m.ik {
		if gotIdx, ok := m.ki[key]; !ok || gotIdx != idx {
			t.Errorf("inconsistency for key %d: expected idx %d, got %d, %v", key, idx, gotIdx, ok)
		}
	}
	for key, idx := range m.ki {
		if gotKey, ok := m.ik[idx]; !ok || gotKey != key {
			t.Errorf("inconsistency for idx %d: expected key %d, got %d, %v", idx, key, gotKey, ok)
		}
	}
}

func TestRange(t *testing.T) {
	m := Make[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	t.Run("IterateAll", func(t *testing.T) {
		seen := make(map[string]int)
		for k, v := range m.Range {
			seen[k] = v
		}
		if len(seen) != 3 {
			t.Errorf("expected to see 3 items, saw %d", len(seen))
		}
		for k, v := range m.kv {
			if seenV, ok := seen[k]; !ok || seenV != v {
				t.Errorf("mismatch for %s: expected %d, got %d", k, v, seenV)
			}
		}
	})

	t.Run("EarlyStop", func(t *testing.T) {
		count := 0
		for _, _ = range m.Range {
			count++
			break
		}
		if count != 1 {
			t.Errorf("expected to iterate 1 time, got %d", count)
		}
	})

	t.Run("EmptyMap", func(t *testing.T) {
		empty := Make[string, int]()
		called := false
		for _, _ = range empty.Range {
			called = true
			break
		}
		if called {
			t.Error("expected no calls on empty map")
		}
	})
}

func TestWrapEmpty(t *testing.T) {
	orig := map[string]int{}
	m := Wrap(orig)
	if m.Len() != 0 {
		t.Errorf("expected len 0, got %d", m.Len())
	}
	_, ok := m.GetRandom()
	if ok {
		t.Error("expected no value from empty wrapped map")
	}
}

func TestWrap(t *testing.T) {
	orig := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	m := Wrap(orig)

	if m.Len() != 3 {
		t.Errorf("expected len 3, got %d", m.Len())
	}

	val, ok := m.Get("a")
	if !ok || val != 1 {
		t.Errorf("expected 1 for 'a', got %v, %v", val, ok)
	}
	val, ok = m.Get("b")
	if !ok || val != 2 {
		t.Errorf("expected 2 for 'b', got %v, %v", val, ok)
	}
	val, ok = m.Get("c")
	if !ok || val != 3 {
		t.Errorf("expected 3 for 'c', got %v, %v", val, ok)
	}

	// Check internal consistency
	if len(m.ik) != 3 || len(m.ki) != 3 {
		t.Errorf("expected ik and ki len 3, got %d, %d", len(m.ik), len(m.ki))
	}
	seenKeys := make(map[string]bool)
	for i := 0; i < 3; i++ {
		key, ok := m.ik[i]
		if !ok {
			t.Errorf("missing ik[%d]", i)
		}
		seenKeys[key] = true
		idx, ok := m.ki[key]
		if !ok || idx != i {
			t.Errorf("ki[%s] expected %d, got %d, %v", key, i, idx, ok)
		}
	}
	if len(seenKeys) != 3 {
		t.Errorf("expected 3 unique keys in ik, got %d", len(seenKeys))
	}
}

func TestWrapGetRandom(t *testing.T) {
	orig := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	m := Wrap(orig)

	// Call multiple times to check validity
	seen := make(map[int]bool)
	for i := 0; i < 100; i++ {
		val, ok := m.GetRandom()
		if !ok {
			t.Error("expected value from non-empty map")
		}
		if val != 1 && val != 2 && val != 3 {
			t.Errorf("unexpected value %d", val)
		}
		seen[val] = true
	}
	if len(seen) != 3 {
		t.Errorf("expected to see all 3 values in 100 calls, saw %d", len(seen))
	}
}

func TestWrapOperations(t *testing.T) {
	orig := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	m := Wrap(orig)

	// Set new
	m.Set("d", 4)
	if m.Len() != 4 {
		t.Errorf("expected len 4, got %d", m.Len())
	}
	val, ok := m.Get("d")
	if !ok || val != 4 {
		t.Errorf("expected 4 for 'd', got %v, %v", val, ok)
	}

	// Update existing
	m.Set("a", 10)
	if val, ok := m.Get("a"); !ok || val != 10 {
		t.Errorf("expected 10 for 'a', got %v, %v", val, ok)
	}
	if m.Len() != 4 {
		t.Errorf("expected len 4 after update, got %d", m.Len())
	}

	// Delete
	m.Delete("b")
	if m.Len() != 3 {
		t.Errorf("expected len 3 after delete, got %d", m.Len())
	}
	_, ok = m.Get("b")
	if ok {
		t.Error("expected 'b' to be deleted")
	}

	// GetRandom after operations
	seen := make(map[int]bool)
	for i := 0; i < 100; i++ {
		val, ok := m.GetRandom()
		if !ok {
			t.Error("expected value")
		}
		if val != 10 && val != 3 && val != 4 {
			t.Errorf("unexpected value %d", val)
		}
		seen[val] = true
	}
	if len(seen) != 3 {
		t.Errorf("expected to see 3 values, saw %d", len(seen))
	}
}
