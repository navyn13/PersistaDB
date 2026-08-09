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
		if err := s.storage.Set(parts[1], parts[2]); err != nil {
			msg.peer.Send([]byte("ERROR WRITING VALUE\r\n"))
			return err
		}
		msg.peer.Send([]byte("OK\r\n"))
		return nil

	case "get":
		if len(parts) != 2 {
			msg.peer.Send([]byte("INVALID GET COMMAND\r\n"))
			return nil
		}
		value, err := s.storage.Get(parts[1])
		if err != nil {
			msg.peer.Send([]byte("ERROR READING VALUE\r\n"))
			return err
		}
		msg.peer.Send([]byte(value + "\r\n"))
		return nil

	case "delete":
		if len(parts) != 2 {
			msg.peer.Send([]byte("INVALID DELETE COMMAND\r\n"))
			return nil
		}
		if err := s.storage.Delete(parts[1]); err != nil {
			msg.peer.Send([]byte("ERROR DELETING VALUE\r\n"))
			return err
		}
		msg.peer.Send([]byte("OK\r\n"))
		return nil
	default:
		msg.peer.Send([]byte("INVALID COMMAND\r\n"))
		return nil
	}
}
