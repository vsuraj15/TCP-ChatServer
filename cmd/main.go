package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jessevdk/go-flags"
	"loconav.com/projects/chat/internal/config"
	"loconav.com/projects/chat/version"
)

type Options struct {
	ServerPort int  `short:"p" long:"port" description:"Flag to set server port number"`
	Version    bool `short:"v" long:"version" description:"Flag to get the app version"`
}

var (
	serverHost = "0.0.0.0"
	protocol   = "tcp"
)

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalf("Error while reading flag values.Err: %+v", err)
	}
	if opts.Version {
		version.DisplayVersion("Real-Time Chat App")
		return
	}
	if config.HaltIfEmpty(opts.ServerPort) {
		log.Fatalf("Failed to get server port number. Expected > 0, Found: %+v", opts.ServerPort)
	}

	server := config.NewServer()
	go server.Run()

	serverAddr := fmt.Sprintf("%s:%d", serverHost, opts.ServerPort)
	listener, listenErr := net.Listen(protocol, serverAddr)
	if listenErr != nil {
		log.Fatalf("Error while listening on expected server host and server port.Err: %+v", listenErr)
	}
	defer listener.Close()
	log.Printf("Launch server on address: %s", serverAddr)
	for {
		conn, connErr := listener.Accept()
		if connErr != nil {
			log.Printf("Failed to accept connection from listener.Err: %+v", connErr)
		}
		client := server.NewClient(conn)
		go client.ReadInput()
	}
}
