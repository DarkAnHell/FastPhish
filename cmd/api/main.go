package main

import (
	"log"
	"net"

	"github.com/DarkAnHell/FastPhish/api"
	"google.golang.org/grpc"
)

func main() {
	// TODO: Config
	l, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterAPIServer(s, server{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("could not serve: %v", err)
	}
}
