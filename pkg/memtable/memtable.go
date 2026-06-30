package memtable

import (
	"errors"
	"rubik-store-go/pkg/config"
	"rubik-store-go/pkg/skiplist"
	"rubik-store-go/pkg/wal"
)

// MemTable ensures that writes made to the MemTabel is persisted to disk and is ready for query once
// it has returned success.
type MemTable struct {
	sl             *skiplist.Skiplist
	wal            *wal.WAL
	currentSize    int64
	flushThreshold int64
}

func NewMemTable(storageDir string, o config.Options) (*MemTable, error) {
	// TODO:
	return nil, errors.New("not implemented")
}
