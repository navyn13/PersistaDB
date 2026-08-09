package server

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type WAL struct {
	activeFile *os.File
	mu         sync.Mutex
}

func NewWAL() *WAL {
	w := &WAL{}
	// Ensure data directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil
	}
	return w
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

	filename := fmt.Sprintf("log-%06d.data", maxID+1)
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

	return nil
}
func (w *WAL) Write(data []byte) error {
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
