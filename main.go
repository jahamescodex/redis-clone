package main

import (
	"log/slog"
	"os"
	"os/signal"
	"redis-clone/server"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	ListenAndServe := server.NewServer("")
	err := ListenAndServe.Start()
	if err != nil {
		slog.Info("error", "error", err)
		os.Exit(1)
	}

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownSig
	ListenAndServe.Shutdown()
	signal.Stop(shutdownSig)

}
