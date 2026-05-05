package skiplist

import (
	"bytes"
	"testing"
)

func makeEntry(key, value string) *Entry {
	return &Entry{Key: []byte(key), Value: []byte(value)}
}

func TestNew(t *testing.T) {
	s := New(4)
	if s.sentinel == nil {
		t.Fatal("sentinel is nil")
	}
	if len(s.sentinel.forward) != 4 {
		t.Fatalf("sentinel.forward len = %d, want 4", len(s.sentinel.forward))
	}
	if s.level != 0 {
		t.Fatalf("level = %d, want 0", s.level)
	}
	if s.length != 0 {
		t.Fatalf("length = %d, want 0", s.length)
	}
}

func TestSearchEmpty(t *testing.T) {
	s := New(4)
	if s.Search([]byte("a")) != nil {
		t.Fatal("expected nil on empty skiplist")
	}
}

func TestInsertAndSearch(t *testing.T) {
	s := New(4)
	entries := []struct{ key, val string }{
		{"b", "2"},
		{"a", "1"},
		{"d", "4"},
		{"c", "3"},
	}
	for _, e := range entries {
		s.Insert(makeEntry(e.key, e.val))
	}
	for _, e := range entries {
		got := s.Search([]byte(e.key))
		if got == nil {
			t.Fatalf("Search(%q) = nil, want %q", e.key, e.val)
		}
		if !bytes.Equal(got.Value, []byte(e.val)) {
			t.Fatalf("Search(%q).Value = %q, want %q", e.key, got.Value, e.val)
		}
	}
}

func TestSearchMiss(t *testing.T) {
	s := New(4)
	s.Insert(makeEntry("a", "1"))
	s.Insert(makeEntry("c", "3"))
	if s.Search([]byte("b")) != nil {
		t.Fatal("expected nil for key between existing keys")
	}
	if s.Search([]byte("z")) != nil {
		t.Fatal("expected nil for key beyond range")
	}
}

func TestInsertUpdatesExistingKey(t *testing.T) {
	s := New(4)
	s.Insert(makeEntry("a", "old"))
	s.Insert(makeEntry("a", "new"))
	got := s.Search([]byte("a"))
	if got == nil {
		t.Fatal("Search returned nil after update")
	}
	if !bytes.Equal(got.Value, []byte("new")) {
		t.Fatalf("got value %q, want \"new\"", got.Value)
	}
	if s.length != 1 {
		t.Fatalf("length = %d after update, want 1", s.length)
	}
}

func TestInsertLength(t *testing.T) {
	s := New(4)
	keys := []string{"e", "b", "a", "d", "c"}
	for i, k := range keys {
		s.Insert(makeEntry(k, k))
		if s.length != i+1 {
			t.Fatalf("after %d inserts, length = %d, want %d", i+1, s.length, i+1)
		}
	}
}

func TestRemoveExisting(t *testing.T) {
	s := New(4)
	s.Insert(makeEntry("a", "1"))
	s.Insert(makeEntry("b", "2"))
	s.Insert(makeEntry("c", "3"))
	s.Remove([]byte("b"))
	if s.Search([]byte("b")) != nil {
		t.Fatal("Search returned non-nil after Remove")
	}
	if s.Search([]byte("a")) == nil {
		t.Fatal("key \"a\" missing after removing \"b\"")
	}
	if s.Search([]byte("c")) == nil {
		t.Fatal("key \"c\" missing after removing \"b\"")
	}
	if s.length != 2 {
		t.Fatalf("length = %d after remove, want 2", s.length)
	}
}

func TestRemoveFirst(t *testing.T) {
	s := New(4)
	s.Insert(makeEntry("a", "1"))
	s.Insert(makeEntry("b", "2"))
	s.Remove([]byte("a"))
	if s.Search([]byte("a")) != nil {
		t.Fatal("first element still searchable after remove")
	}
	if s.Search([]byte("b")) == nil {
		t.Fatal("second element missing after removing first")
	}
	if s.length != 1 {
		t.Fatalf("length = %d, want 1", s.length)
	}
}

func TestRemoveLast(t *testing.T) {
	s := New(4)
	s.Insert(makeEntry("a", "1"))
	s.Insert(makeEntry("b", "2"))
	s.Remove([]byte("b"))
	if s.Search([]byte("b")) != nil {
		t.Fatal("last element still searchable after remove")
	}
	if s.Search([]byte("a")) == nil {
		t.Fatal("first element missing after removing last")
	}
	if s.length != 1 {
		t.Fatalf("length = %d, want 1", s.length)
	}
}

func TestRemoveNonExistent(t *testing.T) {
	s := New(4)
	s.Insert(makeEntry("a", "1"))
	s.Insert(makeEntry("c", "3"))
	s.Remove([]byte("b"))
	if s.length != 2 {
		t.Fatalf("length = %d after removing non-existent key, want 2", s.length)
	}
	if s.Search([]byte("a")) == nil || s.Search([]byte("c")) == nil {
		t.Fatal("existing keys missing after removing non-existent key")
	}
}

func TestRemoveEmpty(t *testing.T) {
	s := New(4)
	s.Remove([]byte("a")) // must not panic
}

// TestForwardChainSingleLevel uses maxLevel=1 to disable level promotion,
// giving a deterministic level-0 chain we can walk directly.
func TestForwardChainSingleLevel(t *testing.T) {
	s := New(1)
	for _, k := range []string{"c", "a", "b"} {
		s.Insert(makeEntry(k, k))
	}

	// Expect: sentinel -> a -> b -> c -> nil
	want := []string{"a", "b", "c"}
	n := s.sentinel
	for i, key := range want {
		n = n.forward[0]
		if n == nil {
			t.Fatalf("chain too short at position %d, want key %q", i, key)
		}
		if !bytes.Equal(n.entry.Key, []byte(key)) {
			t.Fatalf("position %d: got %q, want %q", i, n.entry.Key, key)
		}
	}
	if n.forward[0] != nil {
		t.Fatal("chain longer than expected after inserts")
	}

	s.Remove([]byte("b"))

	// Expect: sentinel -> a -> c -> nil
	want = []string{"a", "c"}
	n = s.sentinel
	for i, key := range want {
		n = n.forward[0]
		if n == nil {
			t.Fatalf("after remove, chain too short at position %d, want key %q", i, key)
		}
		if !bytes.Equal(n.entry.Key, []byte(key)) {
			t.Fatalf("after remove, position %d: got %q, want %q", i, n.entry.Key, key)
		}
	}
	if n.forward[0] != nil {
		t.Fatal("chain longer than expected after remove")
	}
}
