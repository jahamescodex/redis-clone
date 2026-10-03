package server

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
)

const defaultListener = ":6767"

type Server struct {
	ln         net.Listener
	listenAddr string

	isRunning atomic.Bool

	Wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	conMu   sync.RWMutex
	connMap map[net.Conn]struct{}
}

func NewServer(listenAddr string) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	if len(listenAddr) == 0 {
		return &Server{listenAddr: defaultListener, connMap: make(map[net.Conn]struct{}),
			Wg: sync.WaitGroup{}, ctx: ctx, cancel: cancel}
	}
	return &Server{listenAddr: listenAddr, connMap: make(map[net.Conn]struct{}),
		Wg: sync.WaitGroup{}, ctx: ctx, cancel: cancel}
}

func (s *Server) Start() error {
	var err error
	s.ln, err = net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.isRunning.Store(true)
	return nil
}

func (s *Server) AcceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if !s.isRunning.Load() { // non-read method
				return nil
			}
			slog.Info("Accept error", "err", err)
			continue
		}
		s.conMu.Lock()
		if !s.isRunning.Load() {
			s.conMu.Unlock()
			continue
		}

		s.connMap[conn] = struct{}{}
		s.conMu.Unlock()

		s.Wg.Add(1)
		go s.handleConnection(conn)

	}
}

func (s *Server) Shutdown() {
	s.isRunning.Store(false)
	s.ln.Close()
	s.cancel()
}
