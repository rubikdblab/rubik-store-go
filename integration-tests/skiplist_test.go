package integration_test

import (
	"bytes"
	"testing"

	"rubik-store-go/pkg/skiplist"
)

type op struct {
	action string // "insert", "remove", "search"
	key    string
	value  string // used by insert; expected value for search ("" means expect nil)
}

func runOps(t *testing.T, s *skiplist.Skiplist, ops []op) {
	t.Helper()
	for _, o := range ops {
		switch o.action {
		case "insert":
			s.Insert(&skiplist.Entry{Key: []byte(o.key), Value: []byte(o.value)})
		case "remove":
			s.Remove([]byte(o.key))
		case "search":
			got := s.Search([]byte(o.key))
			if o.value == "" {
				if got != nil {
					t.Errorf("Search(%q) = %q, want nil", o.key, got.Value)
				}
			} else {
				if got == nil {
					t.Errorf("Search(%q) = nil, want %q", o.key, o.value)
				} else if !bytes.Equal(got.Value, []byte(o.value)) {
					t.Errorf("Search(%q) = %q, want %q", o.key, got.Value, o.value)
				}
			}
		}
	}
}

func TestSkiplist(t *testing.T) {
	tests := []struct {
		name     string
		maxLevel int
		ops      []op
	}{
		{
			name:     "search on empty list",
			maxLevel: 4,
			ops: []op{
				{action: "search", key: "a", value: ""},
			},
		},
		{
			name:     "insert and search single key",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "a", value: "1"},
				{action: "search", key: "a", value: "1"},
			},
		},
		{
			name:     "insert multiple keys in unsorted order",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "c", value: "3"},
				{action: "insert", key: "a", value: "1"},
				{action: "insert", key: "d", value: "4"},
				{action: "insert", key: "b", value: "2"},
				{action: "search", key: "a", value: "1"},
				{action: "search", key: "b", value: "2"},
				{action: "search", key: "c", value: "3"},
				{action: "search", key: "d", value: "4"},
			},
		},
		{
			name:     "search for missing key",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "a", value: "1"},
				{action: "insert", key: "c", value: "3"},
				{action: "search", key: "b", value: ""},
				{action: "search", key: "z", value: ""},
			},
		},
		{
			name:     "update existing key",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "a", value: "old"},
				{action: "insert", key: "a", value: "new"},
				{action: "search", key: "a", value: "new"},
			},
		},
		{
			name:     "remove middle key",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "a", value: "1"},
				{action: "insert", key: "b", value: "2"},
				{action: "insert", key: "c", value: "3"},
				{action: "remove", key: "b"},
				{action: "search", key: "b", value: ""},
				{action: "search", key: "a", value: "1"},
				{action: "search", key: "c", value: "3"},
			},
		},
		{
			name:     "remove first key",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "a", value: "1"},
				{action: "insert", key: "b", value: "2"},
				{action: "remove", key: "a"},
				{action: "search", key: "a", value: ""},
				{action: "search", key: "b", value: "2"},
			},
		},
		{
			name:     "remove last key",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "a", value: "1"},
				{action: "insert", key: "b", value: "2"},
				{action: "remove", key: "b"},
				{action: "search", key: "b", value: ""},
				{action: "search", key: "a", value: "1"},
			},
		},
		{
			name:     "remove only key",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "a", value: "1"},
				{action: "remove", key: "a"},
				{action: "search", key: "a", value: ""},
			},
		},
		{
			name:     "remove non-existent key",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "a", value: "1"},
				{action: "insert", key: "c", value: "3"},
				{action: "remove", key: "b"},
				{action: "search", key: "a", value: "1"},
				{action: "search", key: "c", value: "3"},
			},
		},
		{
			name:     "remove from empty list",
			maxLevel: 4,
			ops: []op{
				{action: "remove", key: "a"},
			},
		},
		{
			name:     "insert after remove",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "a", value: "1"},
				{action: "remove", key: "a"},
				{action: "insert", key: "a", value: "2"},
				{action: "search", key: "a", value: "2"},
			},
		},
		{
			name:     "interleaved inserts and removes",
			maxLevel: 4,
			ops: []op{
				{action: "insert", key: "b", value: "2"},
				{action: "insert", key: "a", value: "1"},
				{action: "insert", key: "d", value: "4"},
				{action: "remove", key: "b"},
				{action: "insert", key: "c", value: "3"},
				{action: "remove", key: "a"},
				{action: "search", key: "a", value: ""},
				{action: "search", key: "b", value: ""},
				{action: "search", key: "c", value: "3"},
				{action: "search", key: "d", value: "4"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := skiplist.New(tt.maxLevel)
			runOps(t, s, tt.ops)
		})
	}
}
