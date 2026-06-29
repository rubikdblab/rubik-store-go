package wal

import (
	"fmt"
	"os"
	"sync"
)

type RecordType uint8

const (
	RecordPut    RecordType = 0
	RecordDelete RecordType = 1
)

// Record is a single WAL entry representing one Put or Delete operation.
type Record struct {
	CRC      uint32 // populated by Decode; ignored by Encode
	SeqNum   uint64
	Type     RecordType
	KeyLen   uint32
	ValueLen uint32
	Key      []byte
	Value    []byte
}

type WAL struct {
	sync.Mutex
	file *os.File
	size int64
}

func NewWAL(filename string) (*WAL, error) {
	flag := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile(filename, flag, 0644)
	if err != nil {
		return nil, fmt.Errorf("failure in opening WAL: %w", err)
	}
	return &WAL{
		file: file,
	}, nil
}

// Write appends a record to the WAL.
func (w *WAL) Write(record *Record) error {
	w.Lock()
	defer w.Unlock()

	buf := Encode(record)
	n, err := w.file.Write(buf)
	if err != nil {
		return fmt.Errorf("wal: failed to write record: %w", err)
	}
	w.file.Sync()
	w.size += int64(n)
	return nil
}
