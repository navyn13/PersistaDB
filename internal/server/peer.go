package server

import (
	"net"
	"sync"
)

type Message struct {
	Peer *Peer
	Msg  string
}

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	mu    sync.Mutex // serializes writes (pipelined msgs may run on many workers)
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: make(chan Message),
	}
}
