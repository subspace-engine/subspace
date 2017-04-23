package main

import (
	"flag"
	"fmt"
	"github.com/subspace-engine/subspace/engine/server"
	"os"
	"os/signal"
)

func main() {
	port := flag.Int("port", 11111, "TCP port")
	done := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	srv := server.NewServer(*port, done)
	go srv.Run()

	sig := <-sigChan
	fmt.Printf("Received signal %s", sig)
	close(done)
	fmt.Println("Exiting normally")
}
