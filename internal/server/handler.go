package server

import (
	"strings"
)

func (s *Server) handleMessage(msg Message) error {
	parts := strings.SplitN(msg.data, " ", 3)

	switch parts[0] {
	case "set":
		if len(parts) != 3 {
			msg.peer.Send([]byte("INVALID SET COMMAND\r\n"))
			return nil
		}
		encoder := &SetRecordEncoder{
			key:   parts[1],
			value: parts[2],
		}
		encoded := encoder.Encode()
		err := s.wal.WriteRecord(encoded)
		if err != nil {
			msg.peer.Send([]byte("ERROR WRITING TO WAL\r\n"))
			return err
		}
		msg.peer.Send([]byte("OK\r\n"))
		return nil

	case "get":
		if len(parts) != 2 {
			msg.peer.Send([]byte("INVALID GET COMMAND\r\n"))
			return nil
		}
		value, err := s.wal.ReadRecord(parts[1])
		if err != nil {
			msg.peer.Send([]byte("ERROR READING FROM WAL\r\n"))
			return err
		}
		msg.peer.Send([]byte(value + "\r\n"))
		return nil

	case "delete":
		if len(parts) != 2 {
			msg.peer.Send([]byte("INVALID DELETE COMMAND\r\n"))
			return nil
		}
		encoder := &DeleteRecordEncoder{
			key: parts[1],
		}
		encoded := encoder.Encode()
		err := s.wal.WriteRecord(encoded)
		if err != nil {
			msg.peer.Send([]byte("ERROR DELETING TO WAL\r\n"))
			return err
		}
		msg.peer.Send([]byte("OK\r\n"))
		return nil
	default:
		msg.peer.Send([]byte("INVALID COMMAND\r\n"))
		return nil
	}
}
