// This package is a test client. Don't use it.
package main

import (
	"context"
	"log"

	"github.com/DarkAnHell/FastPhish/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := api.NewDBClient(conn)
	dscli, err := client.GetDomainsScore(context.Background())
	if err != nil {
		log.Fatalf("could not create DomainsScoreClient: %v", err)
	}

	domains := []string{"twitter.com", "fb.com", "hackyhacky.es"}
	for _, v := range domains {
		log.Println("Sending msg from client...")
		domain := &api.Domain{
			Name: v,
		}
		if err := dscli.Send(domain); err != nil {
			log.Printf("could not send request: %v\n", err)
			return
		}

		resp, err := dscli.Recv()
		if err != nil {
			log.Printf("could not read response: %v\n", err)
			return
		}
		log.Printf("Got response with status %v: %s with score: %v\n", resp.GetStatus().Status, resp.GetDomain().Name, resp.GetDomain().Score)
	}
	if err := dscli.CloseSend(); err != nil {
		log.Fatalf("could not close send: %v\n", err)
	}
}
