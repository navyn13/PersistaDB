package server

import (
	"log/slog"
	"strings"
)

func (s *Server) handleMessage(msg Message) error {
	parts := strings.Split(msg.data, " ")

	switch parts[0] {
	case "set":
		err := s.wal.Write([]byte(msg.data))
		if err != nil {
			msg.peer.Send([]byte("ERROR WRITING TO WAL\r\n"))
			return err
		}
		msg.peer.Send([]byte("OK\r\n"))
		return nil

	case "get":
		slog.Info("Getting value", "key", parts[1])
		return nil
	default:
		msg.peer.Send([]byte("INVALID COMMAND\r\n"))
		return nil
	}
}
