package main

import (
	"log"
	"net"

	"github.com/DarkAnHell/FastPhish/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// TODO: Config
	l, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	// Create the TLS credentials
	// TODO: config.
    creds, err := credentials.NewServerTLSFromFile("certs/server.crt", "certs/server.key")
    if err != nil {
        log.Fatalf("could not load TLS keys: %s", err)
    }

	s := grpc.NewServer(grpc.Creds(creds))
	api.RegisterAPIServer(s, server{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("could not serve: %v", err)
	}
}
