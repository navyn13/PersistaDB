package server

import (
	"log/slog"
	"strings"
)

func (s *Server) handleMessage(msg Message) error {
	parts := strings.Split(msg.data, " ")

	switch parts[0] {
	case "set":
		if len(parts) != 3 {
			msg.peer.Send([]byte("INVALID SET COMMAND\r\n"))
			return nil
		}
		err := s.wal.Write(msg.data)
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
		slog.Info("Getting value", "key", parts[1])
		return nil

	case "delete":
		if len(parts) != 2 {
			msg.peer.Send([]byte("INVALID DELETE COMMAND\r\n"))
			return nil
		}
		err := s.wal.Write(msg.data)
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
