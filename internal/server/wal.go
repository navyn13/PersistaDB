package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type WAL struct {
	activeFile *os.File
	nextFileID int
	mu         sync.Mutex
	keyDir     map[string]KeyDirEntry
}

const HEADER_SIZE = 20

type KeyDirEntry struct {
	FileName  string
	Offset    int64
	ValueSize uint32
}

func NewWAL() *WAL {
	w := &WAL{
		keyDir: make(map[string]KeyDirEntry),
	}
	// Ensure data directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil
	}
	return w
}
func (w *WAL) BuildKeyDirMapFromLogFiles() error {
	files, err := filepath.Glob(filepath.Join("logs", "*.data"))
	if err != nil {
		return err
	}
	for _, file := range files {
		fi, err := os.Open(file)
		if err != nil {
			return err
		}
		fileN := fi.Name()
		err = w.buildKeyDirFromLogFile(fi, fileN)
		fi.Close()
		if err != nil {
			return err
		}
	}
	//print the whole map after reading the all the log files
	fmt.Println(w.keyDir)
	return nil

}

func (w *WAL) buildKeyDirFromLogFile(fi *os.File, fileName string) error {
	for {
		offset, err := fi.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		var header Header
		if err := binary.Read(fi, binary.BigEndian, &header); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		key := make([]byte, header.KeyLen)
		if _, err := io.ReadFull(fi, key); err != nil {
			return err
		}

		if _, err := fi.Seek(int64(header.ValueLen), io.SeekCurrent); err != nil {
			return err
		}

		w.keyDir[string(key)] = KeyDirEntry{
			FileName:  fileName,
			Offset:    offset,
			ValueSize: header.ValueLen,
		}
	}

	return nil
}

// CreateDataFile creates the next sequential .data file.
func (w *WAL) CreateNextLogFile() error {
	// Find next available segment number
	files, err := filepath.Glob(filepath.Join("logs", "*.data"))
	if err != nil {
		return nil
	}

	maxID := 0
	for _, file := range files {
		var id int
		name := filepath.Base(file)

		// Parse names like 000001.data
		if _, err := fmt.Sscanf(name, "%06d.data", &id); err == nil {
			if id > maxID {
				maxID = id
			}
		}
	}

	filename := fmt.Sprintf("%06d.data", maxID+1)
	path := filepath.Join("logs", filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	// Close previous active file if one exists
	if w.activeFile != nil {
		_ = w.activeFile.Close()
	}

	w.activeFile = file
	w.nextFileID = maxID + 1

	return nil
}
func (w *WAL) WriteRecord(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.activeFile == nil {
		return fmt.Errorf("no active file")
	}
	_, err := w.activeFile.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (w *WAL) ReadRecord(key string) (string, error) {
	entry, ok := w.keyDir[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}

	fi, err := os.Open(entry.FileName)
	if err != nil {
		return "", err
	}
	defer fi.Close()

	// Move to the beginning of the record
	if _, err := fi.Seek(entry.Offset, io.SeekStart); err != nil {
		return "", err
	}

	// Read header
	var header Header
	if err := binary.Read(fi, binary.BigEndian, &header); err != nil {
		return "", err
	}

	// Skip the key
	if _, err := fi.Seek(int64(header.KeyLen), io.SeekCurrent); err != nil {
		return "", err
	}

	// Read value
	data := make([]byte, header.ValueLen)
	if _, err := io.ReadFull(fi, data); err != nil {
		return "", err
	}

	return string(data), nil
}
