package wal

import (
	"encoding/binary"
	"hash/crc32"
)

/**
*	V1 Record Format
*	Offset  Size  Field      Notes
*	------  ----  ---------  -----
*	0       4     Length     uint32 LE — byte count of everything below (CRC..Value)
*	4       4     CRC32      uint32 LE — checksum over Version..Value
*	8       1     Version    uint8  — record format version, starts at 1
*	9       8     SeqNum     uint64 LE
*	17      1     Type       uint8  — RecordPut | RecordDelete
*	18      4     KeyLen     uint32 LE
*	22      4     ValueLen   uint32 LE
*	26      ...   Key        KeyLen bytes
*	26+KL   ...   Value      ValueLen bytes
 */

const (
	version1       uint8 = 1
	currentVersion       = version1

	lengthSize  = 4
	crcSize     = 4
	versionSize = 1

	// v1 header fields after the version byte: SeqNum(8) + Type(1) + KeyLen(4) + ValueLen(4)
	v1HeaderSize = 8 + 1 + 4 + 4
)

// Encode serializes a Record using the current WAL record format version.
func Encode(record *Record) []byte {
	keyLen := uint32(len(record.Key))
	valueLen := uint32(len(record.Value))

	payload := make([]byte, versionSize+v1HeaderSize+int(keyLen)+int(valueLen))
	payload[0] = currentVersion
	binary.LittleEndian.PutUint64(payload[1:9], record.SeqNum)
	payload[9] = byte(record.Type)
	binary.LittleEndian.PutUint32(payload[10:14], keyLen)
	binary.LittleEndian.PutUint32(payload[14:18], valueLen)
	copy(payload[18:], record.Key)
	copy(payload[18+keyLen:], record.Value)

	buf := make([]byte, lengthSize+crcSize+len(payload))
	binary.LittleEndian.PutUint32(buf[0:lengthSize], uint32(crcSize+len(payload)))
	binary.LittleEndian.PutUint32(buf[lengthSize:lengthSize+crcSize], crc32.ChecksumIEEE(payload))
	copy(buf[lengthSize+crcSize:], payload)

	return buf
}
