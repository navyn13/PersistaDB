package server

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type WAL struct {
	activeFile *os.File
	nextID     int
	mu         sync.Mutex
}

func NewWAL() *WAL {
	w := &WAL{}

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil
	}

	// Find next available segment number
	files, err := filepath.Glob(filepath.Join("data", "*.data"))
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

	w.nextID = maxID + 1

	return w
}

// CreateDataFile creates the next sequential .data file.
func (w *WAL) CreateDataFile() error {
	filename := fmt.Sprintf("%06d.data", w.nextID)
	path := filepath.Join("data", filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	// Close previous active file if one exists
	if w.activeFile != nil {
		_ = w.activeFile.Close()
	}

	w.activeFile = file
	w.nextID++

	return nil
}
func (w *WAL) Write(data string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.activeFile == nil {
		return fmt.Errorf("no active file")
	}
	_, err := w.activeFile.WriteString(data + "\n")
	if err != nil {
		return err
	}
	return nil
}
