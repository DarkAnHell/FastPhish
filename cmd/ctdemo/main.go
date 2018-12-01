package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/ct"
	"github.com/DarkAnHell/FastPhish/pkg/db"
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

	cfg, err := ct.Load(f)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	ct, err := ct.New(cfg)
	if err != nil {
		log.Fatalf("could not create CT client: %v", err)
	}

	var dbconn *grpc.ClientConn

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// TODO: config.
	analconn, err := grpc.Dial("localhost:1338", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect to analyzer service: %v", err)
	}
	defer analconn.Close()

	// TODO: config.
	dbconn, err = grpc.Dial("localhost:50000", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer dbconn.Close()

	aclient := api.NewAnalyzerClient(analconn)
	anal, err := aclient.Analyze(context.Background())
	if err != nil {
		log.Fatalf("could not create analyzer: %v", err)
	}

	dbclient := api.NewDBClient(dbconn)
	db, err := dbclient.GetDomainsScore(context.Background())
	if err != nil {
		log.Fatalf("could not create domain score receiver: %v", err)
	}

	domains := make(chan api.Domain)
	done := make(chan struct{})
	go func() {
		if err := ct.Handle(context.Background(), domains); err != nil {
			log.Fatalf("failed to handle CT Logs: %v", err)
		}
		close(done)
	}()
	go func() {
		for {
			select {
			case d := <-domains:
				ok, err := isNewDomain(db, d)
				if err != nil {
					log.Printf("failed to check in DB: %v", err)
					continue
				}
				if !ok {
					log.Printf("already had domain %s", d.Name)
					continue
				}
				log.Printf("new domain %s, analyzing", d.Name)
				if err := anal.Send(&d); err != nil {
					log.Printf("failed to send domain to analyzer: %v", err)
					return
				}

				resp, err := anal.Recv()
				if err != nil {
					log.Printf("could not read response: %v", err)
					return
				}
				log.Printf("Got response with status %v: %s with score: %v\n", resp.GetStatus().Status, resp.GetDomain().Name, resp.GetDomain().Score)
			}
		}
	}()
	<-done
	if err := db.CloseSend(); err != nil {
		log.Fatalf("could not close domains client: %v", err)
	}
}

func isNewDomain(dbcli api.DB_GetDomainsScoreClient, d api.Domain) (bool, error) {
	err := dbcli.Send(&d)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("could not query DB with domain %s: %v", d.Name, err)
	}

	recv, err := dbcli.Recv()
	if err == db.ErrDBNotFound || recv.Status.Status == api.StatusCode_DOMAIN_NOT_FOUND_ON_DB {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not read response from DB: %v", err)
	}

	return false, nil
}
