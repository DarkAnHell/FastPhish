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

/*
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

reader := bufio.NewReader(os.Stdin)
for {

	fmt.Print("Enter key: ")
	key, _ := reader.ReadString('\n')
	fmt.Print("Enter score: ")
	score, _ := reader.ReadString('\n')
	score = strings.TrimSuffix(score, "\n")

	scoreInt, _ := strconv.Atoi(score)

	d := api.DomainScore{Name: key, Score: uint32(scoreInt)}

	fmt.Printf("%v\n", db.Store(d))
}
*/