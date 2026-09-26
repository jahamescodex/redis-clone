package main

import (
	"errors"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"redis-clone/server"
	"syscall"
)

func main() {
	ListenAndServe := server.NewServer(":6379")

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownSig
		log.Println("Server Shutdown")
		ListenAndServe.Shutdown()
		signal.Stop(shutdownSig)
	}()

	err := ListenAndServe.Start()
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			return
		}
		slog.Info("Listen Error", "err:", err)
		os.Exit(1)
	}
}
