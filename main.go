package main

import (
	"github.com/progimage/apiserver"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	server := apiserver.NewAPIServer(apiserver.NewAPIServerOpts{})
	server.Run()
}
