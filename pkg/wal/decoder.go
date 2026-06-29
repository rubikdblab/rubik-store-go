package wal

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
)

// Decode reads and verifies a single Record from r.
//
// It returns io.EOF if r ends exactly at a record boundary, or
// io.ErrUnexpectedEOF if r ends partway through a record. The latter
// indicates a torn write left by a crash; callers replaying a WAL should
// stop there and may truncate the file before appending new records.
func Decode(r io.Reader) (*Record, error) {
	lengthBuf := make([]byte, lengthSize)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lengthBuf)
	if length < crcSize+versionSize {
		return nil, fmt.Errorf("wal: invalid record length %d", length)
	}

	body := make([]byte, length)
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, err
	}

	storedCRC := binary.LittleEndian.Uint32(body[:crcSize])
	payload := body[crcSize:]
	if crc32.ChecksumIEEE(payload) != storedCRC {
		return nil, fmt.Errorf("wal: record checksum mismatch")
	}

	switch version := payload[0]; version {
	case version1:
		record, err := decodeV1(payload[versionSize:])
		if err != nil {
			return nil, err
		}
		record.CRC = storedCRC
		return record, nil
	default:
		return nil, fmt.Errorf("wal: unsupported record version %d", version)
	}
}

func decodeV1(b []byte) (*Record, error) {
	if len(b) < v1HeaderSize {
		return nil, fmt.Errorf("wal: truncated record header")
	}

	seqNum := binary.LittleEndian.Uint64(b[0:8])
	recordType := RecordType(b[8])
	keyLen := binary.LittleEndian.Uint32(b[9:13])
	valueLen := binary.LittleEndian.Uint32(b[13:17])

	rest := b[v1HeaderSize:]
	if uint32(len(rest)) != keyLen+valueLen {
		return nil, fmt.Errorf("wal: record length mismatch")
	}

	key := make([]byte, keyLen)
	copy(key, rest[:keyLen])
	value := make([]byte, valueLen)
	copy(value, rest[keyLen:])

	return &Record{
		SeqNum:   seqNum,
		Type:     recordType,
		KeyLen:   keyLen,
		ValueLen: valueLen,
		Key:      key,
		Value:    value,
	}, nil
}
