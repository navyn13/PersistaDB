package server

import (
	"log/slog"
	"net"

	"github.com/navyn13/PersistaDB/internal/config"
)

type Server struct {
	config.Config
	ln        net.Listener
	addPeerCh chan *Peer
	msgCh     chan Message
	peers     map[*Peer]bool
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Config:    *cfg,
		addPeerCh: make(chan *Peer),
		msgCh:     make(chan Message),
		peers:     make(map[*Peer]bool),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		return nil
	}
	s.ln = ln
	go s.Peerloop()
	slog.Info("PersistaDB Running", "listenAddr", s.ListenAddr)
	return s.acceptClientLoop()
}
func (s *Server) Peerloop() {
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		case msg := <-s.msgCh:
			slog.Info("Received message", "message", msg)
		}
	}
}

func (s *Server) acceptClientLoop() error {

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			continue
		}
		go s.handleClientConn(conn)
	}
}

func (s *Server) handleClientConn(conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerCh <- peer
	// write go read loop here - to read messages
}

func (s *Server) Shutdown() {
	if s.ln != nil {
		s.ln.Close()
	}
}
