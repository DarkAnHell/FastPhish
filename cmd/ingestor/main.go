package main

import (
	"context"
	"log"
	"time"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/datasource/whoisds"
	"github.com/DarkAnHell/FastPhish/pkg/ingestor"
)

func main() {
	i := &ingestor.Default{}
	d := make(chan api.Domain)
	source := whoisds.New(time.Now().AddDate(0, 0, -3))
	done := make(chan struct{})

	go func() {
		for v := range d {
			log.Println(v.Name)
		}
		close(done)
	}()

	err := i.Ingest(context.Background(), d, source)
	if err != nil {
		log.Printf("failed to ingest: %v", err)
	}
	log.Println("ingest method finished...")
	close(d)
	log.Println("Exiting...")
	<-done
}
