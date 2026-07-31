package server

import (
	"os"
	"path/filepath"
	"sync"
)

type SegmentGroup struct {
	mu sync.Mutex
}

func NewSegmentGroup() *SegmentGroup {
	// create segment file if not exists
	if err := os.MkdirAll("segments", 0755); err != nil {
		return nil
	}
	return &SegmentGroup{mu: sync.Mutex{}}
}

func (s *SegmentGroup) CreateSegment() error {
	//only create segment file if not exists
	if _, err := os.Stat(filepath.Join("segments", "segment-000001.seg")); os.IsNotExist(err) {
		return os.WriteFile(filepath.Join("segments", "segment-000001.seg"), []byte(""), 0644)
	}
	return nil
}

func (s *SegmentGroup) Write(data string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := os.OpenFile(filepath.Join("segments", "segment-000001.seg"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(data + "\n")
	if err != nil {
		return err
	}
	return file.Close()
}
