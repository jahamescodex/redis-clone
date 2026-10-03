package main

import (
	"log"
	"os"
	"os/signal"
	"redis-clone/server"
	"syscall"
)

func main() {
	app := server.NewServer(":6379")

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdownSig
		log.Println("Server Shutdown")
		app.Shutdown()
		signal.Stop(shutdownSig)
	}()

	err := app.Start()
	if err != nil {
		log.Fatal("Listen Error:", err)
	}

	if err := app.AcceptLoop(); err != nil {
		log.Fatal(err)
	}
	app.Wg.Wait()
}
