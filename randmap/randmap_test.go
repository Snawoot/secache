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
