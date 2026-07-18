package server

import "log/slog"

func (s *Server) handleMessage(msg Message) error {
	slog.Info("Received msg", "msg", msg.Msg)
	return nil
}
