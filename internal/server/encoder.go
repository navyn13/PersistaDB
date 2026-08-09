package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"time"
)

type Header struct {
	CRC       uint32
	Timestamp int64
	KeyLen    uint32
	ValueLen  uint32
}

type RecordEncoder interface {
	Encode() []byte
}
type SetRecordEncoder struct {
	key   string
	value string
}
type DeleteRecordEncoder struct {
	key string
}

func CRCString(data string) string {
	crc := crc32.ChecksumIEEE([]byte(data))
	return fmt.Sprintf("%08x", crc) // hex string
}

func (e *SetRecordEncoder) Encode() []byte {
	header := Header{
		CRC:       crc32.ChecksumIEEE([]byte(e.key + e.value)),
		Timestamp: time.Now().UnixNano(),
		KeyLen:    uint32(len(e.key)),
		ValueLen:  uint32(len(e.value)),
	}

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, header)
	buf.WriteString(e.key)
	buf.WriteString(e.value)

	return buf.Bytes()
}
func (e *DeleteRecordEncoder) Encode() []byte {
	header := Header{
		CRC:       crc32.ChecksumIEEE([]byte(e.key)),
		Timestamp: time.Now().UnixNano(),
		KeyLen:    uint32(len(e.key)),
		ValueLen:  tombstoneValueLen,
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, header)
	buf.WriteString(e.key)
	return buf.Bytes()
}
