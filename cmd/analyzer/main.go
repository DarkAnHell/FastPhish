package main

import (
	"net"
	"log"
	"os"

	"github.com/DarkAnHell/FastPhish/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TODO: fuzzy logic for deciding upon scores
func main() {
	// Parse config
	if len(os.Args) < 2 {
		log.Fatalf("missing JSON config file path")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open config file %s: %v", os.Args[1], err)
	}
	s, err := New(f, "twitter.com", "google.com", "facebook.com", "paypal.com", "ebay.com", "yahoo.com")
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	// TODO: Config
	l, err := net.Listen("tcp", ":1338")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	// Create the TLS credentials
	// TODO: config.
	creds, err := credentials.NewServerTLSFromFile("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// attach the Ping service to the server
	api.RegisterAnalyzerServer(grpcServer, *s)

	// start the server
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
