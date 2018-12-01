// This package is a test client. Don't use it.
package main

import (
	"context"
	"log"

	"github.com/DarkAnHell/FastPhish/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// TODO: config.

	var conn *grpc.ClientConn

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// Initiate a connection with the server
	conn, err = grpc.Dial("localhost:1337", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := api.NewAPIClient(conn)

	dscli, err := c.Query(context.Background())
	if err != nil {
		log.Fatalf("error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %v", dscli)

	domains := []string{"twitter.com", "fb.com", "hackyhacky.es", "thisisafakedomain.es"}

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
