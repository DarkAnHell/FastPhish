package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/db"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
}

func (s server) Query(srv api.API_QueryServer) error {
	log.Println("Ready!")
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			res := &api.QueryResult{
				Status: &api.Result{
					Status:  api.StatusCode_READ_C_ERR,
					Message: "could not read from clients stream",
				},
			}
			if serr := srv.Send(res); serr != nil {
				return fmt.Errorf("could not send response: %v", err)
			}
		}

		// Check DB
		var resp *api.QueryResult
		slimResp, err := s.QueryDB(req)
		if err != nil && err != ErrNotFound {
			return errors.Wrapf(err, "could not query DB: %v", err)
		}
		if slimResp != nil {
			resp = &api.QueryResult{
				Domain: &api.DomainScore{
					Name:  slimResp.Domain.Name,
					Score: slimResp.Domain.Score,
				},
				Status: &api.Result{
					Message: slimResp.Status.Message,
					Status:  slimResp.Status.Status,
				},
			}

			// TODO:Config for safeness
			if resp.Domain.Score >= 70 {
				resp.Safe = false
			} else {
				resp.Safe = true
			}
		} else {
			// If not, analyze
			var resp *api.QueryResult
			slimResp, err := s.Analyze(req)
			if err != nil {
				return errors.Wrapf(err, "could not query Analyzer: %v", err)
			}
			resp = &api.QueryResult{
				Domain: &api.DomainScore{
					Name:  slimResp.Domain.Name,
					Score: slimResp.Domain.Score,
				},
				Status: &api.Result{
					Message: slimResp.Status.Message,
					Status:  slimResp.Status.Status,
				}}
			// TODO:Config for safeness
			if resp.Domain.Score >= 70 {
				resp.Safe = false
			} else {
				resp.Safe = true
			}
		}

		if err := srv.Send(resp); err != nil {
			return fmt.Errorf("could not send response thorugh stream: %v", err)
		}
	}

	return nil
}

var ErrNotFound = errors.New("query attribute not found")

// Checks if the DB has that already analyzed
func (s server) QueryDB(domain *api.Domain) (*api.SlimQueryResult, error) {
	var conn *grpc.ClientConn

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// TODO: Config
	conn, err = grpc.Dial("localhost:50000", grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect: %v", err)
	}
	defer conn.Close()

	client := api.NewDBClient(conn)
	dscli, err := client.GetDomainsScore(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "could not create DomainsScoreClient: %v", err)
	}

	if err := dscli.Send(domain); err != nil {
		return nil, errors.Wrapf(err, "could not send request: %v", err)
	}

	resp, err := dscli.Recv()
	if err == db.ErrDBNotFound {
		// Key not found, controlled error
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, "could not read response: %v", err)
	}

	if err := dscli.CloseSend(); err != nil {
		return nil, errors.Wrapf(err, "could not close send: %v", err)
	}

	return resp, nil
}

// Checks if the DB has that already analyzed
func (s server) Analyze(domain *api.Domain) (*api.SlimQueryResult, error) {
	var conn *grpc.ClientConn

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("certs/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// TODO: Config
	conn, err = grpc.Dial("localhost:1338", grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect: %v", err)
	}
	defer conn.Close()

	client := api.NewAnalyzerClient(conn)
	dscli, err := client.Analyze(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "could not create Analyzer: %v", err)
	}

	if err := dscli.Send(domain); err != nil {
		return nil, errors.Wrapf(err, "could not send request: %v", err)
	}

	resp, err := dscli.Recv()
	if err != nil {
		return nil, errors.Wrapf(err, "could not read response: %v", err)
	}
	if err := dscli.CloseSend(); err != nil {
		return nil, errors.Wrapf(err, "could not close send: %v", err)
	}

	return resp, nil
}
