package main

import (
	"log"
	"net"
	"os"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/db/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing JSON config file path")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open config file %s: %v", os.Args[1], err)
	}

	var db redis.Redis
	err = db.Load(f)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// TODO: Config
	l, err := net.Listen("tcp", ":50000")
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

	s := server{DB: db}

	// attach the Ping service to the server
	api.RegisterDBServer(grpcServer, &s)

	// start the server
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
