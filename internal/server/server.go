package server

import (
	"log/slog"
	"net"

	"github.com/navyn13/PersistaDB/internal/config"
)

type Server struct {
	config.Config
	ln           net.Listener
	addPeerCh    chan *Peer
	msgCh        chan Message
	peers        map[*Peer]bool
	wal          *WAL
	segmentGroup *SegmentGroup
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Config:       *cfg,
		addPeerCh:    make(chan *Peer),
		msgCh:        make(chan Message),
		peers:        make(map[*Peer]bool),
		wal:          NewWAL(),
		segmentGroup: NewSegmentGroup(),
	}
}

func (s *Server) Start() error {

	//listen for incoming connections
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		return nil
	}

	//start peer loop
	s.ln = ln
	go s.Peerloop()

	//create log file
	if err := s.wal.CreateNextLogFile(); err != nil {
		slog.Error("Failed to create next log file", "error", err)
		return err
	}
	// create segment file if not exists
	if err := s.segmentGroup.CreateSegment(); err != nil {
		slog.Error("Failed to create segment file", "error", err)
		return err
	}
	slog.Info("PersistaDB Running", "listenAddr", s.ListenAddr)
	return s.acceptClientLoop()
}
func (s *Server) Peerloop() {
	for {
		select {
		case peer := <-s.addPeerCh:
			slog.Info("New peer connected", "peer", peer.conn.RemoteAddr())
			s.peers[peer] = true
		case msg := <-s.msgCh:
			s.handleMessage(msg)
		}
	}
}

func (s *Server) acceptClientLoop() error {

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			continue
		}
		go s.handleClientConn(conn, s.msgCh)
	}
}

func (s *Server) handleClientConn(conn net.Conn, msgCh chan Message) {
	peer := NewPeer(conn, msgCh)
	s.addPeerCh <- peer
	go peer.readLoop()
}

func (s *Server) Shutdown() {
	if s.ln != nil {
		s.ln.Close()
	}
}
