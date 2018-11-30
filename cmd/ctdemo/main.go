package main

import (
	"context"
	"log"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/ct"
)

func main() {
	done := make(chan struct{})
	ctlogURLs := []string{
		"https://ct.googleapis.com/rocketeer/",
		"https://ct.googleapis.com/pilot/",
		"https://ct.googleapis.com/logs/argon2018/",
		"https://ct.googleapis.com/icarus/",
		"https://ct.googleapis.com/skydiver/",
		"https://ct.cloudflare.com/logs/nimbus2019/",
		"https://ct.cloudflare.com/logs/nimbus2020/",
	}
	for _, url := range ctlogURLs {
		go func(url string) {
			ct, err := ct.New(url)
			if err != nil {
				log.Fatalf("could not create CT client: %v\n", err)
			}

			domains := make(chan api.Domain)
			go func() {
				err := ct.Handle(context.Background(), 0, 5000, domains)
				if err != nil {
					log.Printf("could not handle CT Log: %v\n", err)
					return
				}
				close(done)
			}()
			go func() {
				var total int
				for {
					select {
					case d := <-domains:
						total++
						if total%1000 == 0 {
							log.Printf("[%s] domain number %d is %s\n", url, total, d.Name)
						}
					}
				}
			}()
			<-done
			ct.Stop()
		}(url)
	}
	<-done
}
