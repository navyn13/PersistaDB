package server

import "errors"

var ErrKeyNotFound = errors.New("key not found")

type MemIndex map[string]LogEntry

type Storage struct {
	wal      *WAL
	memIndex MemIndex
}

func NewStorage() *Storage {
	return &Storage{
		wal: NewWAL(),
	}
}

func (s *Storage) Open() error {
	s.memIndex = make(MemIndex)
	if err := s.wal.Replay(s.applyRecord); err != nil {
		return err
	}
	return s.wal.CreateNextLogFile()
}

func (s *Storage) Set(key, value string) error {
	encoder := &SetRecordEncoder{
		key:   key,
		value: value,
	}
	entry, err := s.wal.WriteRecord(encoder.Encode())
	if err != nil {
		return err
	}
	s.memIndex[key] = entry
	return nil
}

func (s *Storage) Get(key string) (string, error) {
	entry, ok := s.memIndex[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return s.wal.ReadValue(entry)
}

func (s *Storage) Delete(key string) error {
	encoder := &DeleteRecordEncoder{
		key: key,
	}
	if _, err := s.wal.WriteRecord(encoder.Encode()); err != nil {
		return err
	}
	delete(s.memIndex, key)
	return nil
}

func (s *Storage) Close() error {
	return s.wal.Close()
}

func (s *Storage) applyRecord(record LogRecord) error {
	if record.Deleted {
		delete(s.memIndex, record.Key)
		return nil
	}

	s.memIndex[record.Key] = LogEntry{
		FileName: record.FileName,
		Offset:   record.Offset,
	}
	return nil
}
