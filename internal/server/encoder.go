package server

import (
	"fmt"
	"time"
	"strings"
	"hash/crc32"
)

func CRCString(data string) string {
	crc := crc32.ChecksumIEEE([]byte(data))
	return fmt.Sprintf("%08x", crc) // hex string
}
func EncodeRecord(data string)string{
	crc := CRCString(data)
	timestamp := time.Now().UnixNano()
	parts := strings.SplitN(data, " ", 3)
	key := parts[1]
	value := strings.Join(parts[2:], " ")
	return fmt.Sprintf("[%s][%016x][%016x][%016x][%s][%s]\n",
    crc,
    timestamp,
    len(key),
    len(value),
    key,
    value,
)
}