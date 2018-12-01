package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/DarkAnHell/FastPhish/api"
	lev "github.com/DarkAnHell/FastPhish/pkg/analyzer/levenshtein"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server exposes the analyzer API.
type Server struct {
	lev *lev.Levenshtein
	// TODO: config.
	against []string
}

// Analyze runs analysis logic against the given domains.
func (s Server) Analyze(d api.Analyzer_AnalyzeServer) error {
	log.Println("Ready!")
	for {
		req, err := d.Recv()
		if err != nil {
			return fmt.Errorf("could not receive domain: %v", err)
		}
		log.Printf("Recieved: %s\n", req.GetName())

		out := s.lev.Process(req.GetName(), s.against)
		var max uint32
		for _, score := range out {
			if score.Score > max {
				max = score.Score
			}
		}

		log.Printf("Results for domain %s", req.GetName())
		for index, score := range out {
			log.Printf("got %d score against %s", score.Score, s.against[index])
		}
		log.Printf("Score: %d", max)

		resp := &api.SlimQueryResult{
			Domain: &api.DomainScore{
				Name:  req.GetName(),
				Score: max,
			},
			Status: &api.Result{
				Message: "ANALYZED",
				Status:  api.StatusCode_ANALYZE_OK,
			},
		}

		log.Printf("Storing...")
		err = s.StoreOnDB(resp.Domain)
		if err != nil {
			resp.Status.Message = "ERROR on store"
			resp.Status.Status = api.StatusCode_STORE_S_ERROR
			log.Printf("Error while storing...")
		}

		if err := d.Send(resp); err != nil {
			return fmt.Errorf("could not send analysis result: %v", err)
		}
	}
}

// New creates a server and configures it.
func New(r io.Reader, doms ...string) (*Server, error) {
	l := &lev.Levenshtein{}
	if err := l.Load(r); err != nil {
		return nil, fmt.Errorf("could not create server: %v", err)
	}
	against := make([]string, 0, len(doms))
	for _, v := range doms {
		against = append(against, v)
	}
	return &Server{lev: l, against: against}, nil
}

func (s Server) StoreOnDB(d *api.DomainScore) error {
	var conn *grpc.ClientConn

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// TODO: Config
	conn, err = grpc.Dial("localhost:50000", grpc.WithTransportCredentials(creds))
	if err != nil {
		return errors.Wrapf(err, "failed to connect: %v", err)
	}
	defer conn.Close()

	client := api.NewDBClient(conn)
	dscli, err := client.Store(context.Background())
	if err != nil {
		return errors.Wrapf(err, "could not create DomainsScoreClient: %v", err)
	}

	if err := dscli.Send(d); err != nil {
		return errors.Wrapf(err, "could not send request: %v", err)
	}

	_, err = dscli.Recv()
	if err != nil {
		return errors.Wrapf(err, "could not read response: %v", err)
	}

	if err := dscli.CloseSend(); err != nil {
		return errors.Wrapf(err, "could not close send: %v", err)
	}

	return nil
}
