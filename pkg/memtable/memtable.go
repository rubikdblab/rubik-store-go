package memtable

import "rubik-store-go/pkg/wal"
import "rubik-store-go/pkg/skiplist"

// MemTable ensures that writes made to the MemTabel is persisted to disk and is ready for query once
// it has returned success.
type MemTable struct {
	sl             *skiplist.Skiplist
	wal            *wal.WAL
	currentSize    int64
	flushThreshold int64
}
