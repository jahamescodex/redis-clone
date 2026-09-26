package server

import (
	"errors"
	"log/slog"
	"net"
)

const defaultListener = ""

type Server struct {
	ln         net.Listener
	listenAddr string
	isRunning  bool
}

func NewServer(listenAddr string) *Server {
	if len(listenAddr) == 0 {
		return &Server{listenAddr: defaultListener, isRunning: true}
	}
	return &Server{listenAddr: listenAddr, isRunning: true}
}

func (s *Server) Start() error {
	var err error
	s.ln, err = net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	return s.AcceptLoop()
}

func (s *Server) AcceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return err
			}
			slog.Info("accept error", "err", err)
			continue
		}
		go handleConnection(conn)
	}
}

func (s *Server) Shutdown() {
	s.ln.Close()
}
