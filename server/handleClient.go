package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		buffHeader := make([]byte, 1024) // len = cap = 1024 initially
		return &buffHeader
	},
}

func handleConnection(conn net.Conn) {

	buffHeaderPtr := bufferPool.Get().(*[]byte)
	buffer := *buffHeaderPtr

	defer func() {
		*buffHeaderPtr = (*buffHeaderPtr)[:cap(*buffHeaderPtr)]
		clear(*buffHeaderPtr)
		bufferPool.Put(buffHeaderPtr)
	}()

	processed := 0

	for {
		_, err := conn.Read(buffer[processed:])
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			slog.Info("Connection Read Error", "err", err)
		}
		for {
			fmt.Println(string(buffer))
			conn.Write([]byte("+Ok\r\n"))
		}
	}
}
