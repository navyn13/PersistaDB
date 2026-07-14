package server

import (
	"log/slog"
	"net"

	"github.com/navyn13/PersistaDB/internal/config"
)

type Server struct {
	config.Config
	ln net.Listener
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Config: *cfg,
	}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		return
	}
	s.ln = ln
	slog.Info("PersistaDB Running", "listenAddr", s.ListenAddr)

}

func (s *Server) Shutdown() {
	if s.ln != nil {
		s.ln.Close()
	}
}
