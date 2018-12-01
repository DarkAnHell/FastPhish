package main

import (
	"context"
	"fmt"
	"io"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/db"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type server struct {
}

func (s server) Query(srv api.API_QueryServer) error {
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
			// TODO: Analyze proto
			// If not, analyze
			resp = &api.QueryResult{
				Domain: &api.DomainScore{
					Name:  req.GetName(),
					Score: uint32(0),
				},
				Status: &api.Result{
					Message: "OK",
					Status:  api.StatusCode_GENERIC_OK,
				},
				Safe: true,
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
	// TODO: Config
	conn, err := grpc.Dial("localhost:50000", grpc.WithInsecure())
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
