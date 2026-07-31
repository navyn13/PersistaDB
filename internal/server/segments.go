package server

import (
	"os"
	"path/filepath"
	"sync"
)

type Disk struct {
	mu sync.Mutex
}

func NewDisk() *Disk {
	// create segment file if not exists
	if err := os.MkdirAll("segments", 0755); err != nil {
		return nil
	}
	return &Disk{mu: sync.Mutex{}}
}

func (d *Disk) CreateSegment() error {
	//only create segment file if not exists
	if _, err := os.Stat(filepath.Join("segments", "segment-000001.data")); os.IsNotExist(err) {
		return os.WriteFile(filepath.Join("segments", "segment-000001.data"), []byte(""), 0644)
	}
	return nil
}

func (d *Disk) Write(data string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return os.WriteFile(filepath.Join("segments", "segment-000001.data"), []byte(data), 0644)
}
