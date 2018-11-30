package main

import (
	"os"
	"context"
	"log"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/ct"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing JSON config file path")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open config file %s: %v", os.Args[1], err)
	}

	cfg, err := ct.Load(f)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	ct, err := ct.New(cfg)
	if err != nil {
		log.Fatalf("could not create CT client: %v", err)
	}

	domains := make(chan api.Domain)
	done := make(chan struct{})
	go func() {
		if err := ct.Handle(context.Background(), domains); err != nil {
			log.Fatalf("failed to handle CT Logs: %v", err)
		}
		close(done)
	}()
	go func () {
		var total int
		for {
			select {
			case d := <-domains:
				total++
				if total%1000 == 0 {
					log.Printf("domain number %d is %s\n", total, d.Name)
				}
			}
		}
	}()
	<-done
}
