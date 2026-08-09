package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type WAL struct {
	activeFile *os.File
	nextFileID int
	mu         sync.Mutex
}

const HEADER_SIZE = 20
const tombstoneValueLen = ^uint32(0)

type LogEntry struct {
	FileName string
	Offset   int64
}

type LogRecord struct {
	Key      string
	FileName string
	Offset   int64
	Deleted  bool
}

func NewWAL() *WAL {
	return &WAL{}
}

func (w *WAL) Replay(apply func(LogRecord) error) error {
	files, err := filepath.Glob(filepath.Join("logs", "*.data"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		fi, err := os.Open(file)
		if err != nil {
			return err
		}

		err = w.replayLogFile(fi, fi.Name(), apply)
		fi.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *WAL) replayLogFile(fi *os.File, fileName string, apply func(LogRecord) error) error {
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

		deleted := header.ValueLen == tombstoneValueLen
		if !deleted {
			if _, err := fi.Seek(int64(header.ValueLen), io.SeekCurrent); err != nil {
				return err
			}
		}

		if err := apply(LogRecord{
			Key:      string(key),
			FileName: fileName,
			Offset:   offset,
			Deleted:  deleted,
		}); err != nil {
			return err
		}
	}

	return nil
}

// CreateDataFile creates the next sequential .data file.
func (w *WAL) CreateNextLogFile() error {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// Find next available segment number
	files, err := filepath.Glob(filepath.Join("logs", "*.data"))
	if err != nil {
		return err
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

func (w *WAL) WriteRecord(data []byte) (LogEntry, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.activeFile == nil {
		return LogEntry{}, fmt.Errorf("no active file")
	}

	offset, err := w.activeFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return LogEntry{}, err
	}

	_, err = w.activeFile.Write(data)
	if err != nil {
		return LogEntry{}, err
	}
	return LogEntry{
		FileName: w.activeFile.Name(),
		Offset:   offset,
	}, nil
}

func (w *WAL) ReadValue(entry LogEntry) (string, error) {
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

func (w *WAL) Close() error {
	if w.activeFile == nil {
		return nil
	}
	return w.activeFile.Close()
}
