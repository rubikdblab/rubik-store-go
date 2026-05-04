package wal

import (
	"bufio"
	"os"
	"sync"
)

type RecordType uint8

const (
	RecordPut    RecordType = 0
	RecordDelete RecordType = 1
)

type Record struct {
	CRC      uint32
	SeqNum   uint64
	Type     RecordType
	KeyLen   uint32
	ValueLen uint32
	Key      []byte
	Value    []byte
}

type WAL struct {
	sync.Mutex
	file   *os.File
	writer *bufio.Writer
	size   int64
}
