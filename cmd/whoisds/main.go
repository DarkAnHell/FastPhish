package main

import (
	"github.com/DarkAnHell/FastPhish/api/domain"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/DarkAnHell/FastPhish/pkg/datasource/whoisds"
)

func main() {
	date := time.Now()
	date = date.AddDate(0, 0, -3)
	ws := whoisds.New(date)
	domains := make(chan domain.Domain)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case v := <-domains:
				fmt.Println(v.Name)
			}
		}
	}()
	go func() {
		err := ws.Request(context.Background(), http.DefaultClient, domains)
		if err != nil {
			log.Fatalf("failed to request: %v\n", err)
		}
		close(done)
	}()
	<-done
}
