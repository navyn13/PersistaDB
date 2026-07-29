package server

import (
	"log/slog"
	"strings"
)

func (s *Server) handleMessage(msg Message) error {
	parts := strings.Split(msg.data, " ")
	if len(parts) != 3 {
		msg.peer.Send([]byte("INVALID MESSAGE FORMAT\r\n"))
		return nil
	}

	switch parts[0] {
	case "set":
		slog.Info("Setting value", "key", parts[1], "value", parts[2])
		return nil
	case "get":
		slog.Info("Getting value", "key", parts[1])
		return nil
	default:
		msg.peer.Send([]byte("INVALID COMMAND\r\n"))
		return nil
	}
}
