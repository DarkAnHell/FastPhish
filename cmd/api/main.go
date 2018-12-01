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

	opts := []grpc.ServerOption{grpc.Creds(creds)}

	// create a gRPC server object
	grpcServer := grpc.NewServer(opts...)

	s := server{}

	// attach the Ping service to the server
	api.RegisterAPIServer(grpcServer, &s)

	// start the server
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
