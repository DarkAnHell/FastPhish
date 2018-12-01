package main

import (
	"log"
	"net"
	"os"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/db/redis"
	"google.golang.org/grpc"
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

	s := grpc.NewServer()
	api.RegisterDBServer(s, server{DB: db})
	if err := s.Serve(l); err != nil {
		log.Fatalf("could not serve: %v", err)
	}
}
