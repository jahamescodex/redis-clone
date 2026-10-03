package server

import (
	"net"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		buffHeader := make([]byte, 1024) // len = cap = 1024 initially
		return &buffHeader
	},
}

func (s *Server) handleConnection(conn net.Conn) {
	done := make(chan struct{})
	defer func() {
		close(done)
		s.conMu.Lock()
		delete(s.connMap, conn)
		s.conMu.Unlock()
		s.Wg.Done()
		conn.Close()
	}()

	go func() {
		select {
		case <-s.ctx.Done():
			conn.Close()
		case <-done:
			return
		}
	}()

	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			return
		}
		conn.Write([]byte("+OK\r\n"))
	}

}
