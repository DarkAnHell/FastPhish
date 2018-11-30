package main

import (
	"github.com/DarkAnHell/FastPhish/api"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterDBServer(s, server{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("could not serve: %v", err)
	}
}
