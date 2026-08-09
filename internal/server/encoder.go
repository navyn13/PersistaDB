package server

import (
	"fmt"
	"hash/crc32"
	"time"
)

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
	crc := CRCString(e.key + e.value)
	timestamp := time.Now().UnixNano()
	keyLen := len(e.key)
	valueLen := len(e.value)

	record := fmt.Sprintf("%s%016x%08x%08x%s%s",
		crc,
		timestamp,
		keyLen,
		valueLen,
		e.key,
		e.value,
	)
	return []byte(record)
}
func (e *DeleteRecordEncoder) Encode() []byte {
	crc := CRCString(e.key)
	timestamp := time.Now().UnixNano()
	keyLen := len(e.key)
	valueLen := -1
	value := ""
	///valaue len must be -1 for delete
	record := fmt.Sprintf("[%s][%016x][%016x][%d][%s][%s]",
		crc,
		timestamp,
		keyLen,
		valueLen,
		e.key,
		value,
	)
	return []byte(record)
}
