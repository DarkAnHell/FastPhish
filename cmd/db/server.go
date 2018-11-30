package main

import (
	"io"
	"fmt"

	"github.com/DarkAnHell/FastPhish/api"
)

type server struct{}

func (server) GetDomainsScore(srv api.DB_GetDomainsScoreServer) error {
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			res := &api.SlimQueryResult{
				Status: &api.Result{
					Status: api.StatusCode_READ_C_ERR,
					Message: "could not read from clients stream",
				},
			}
			if serr := srv.Send(res); serr != nil {
				return fmt.Errorf("could not send response: %v", err)
			}
		}

		resp := &api.SlimQueryResult{
			Domain: &api.DomainScore{
				Name: req.GetName(),
				Score: uint32(1),
			},
			Status: &api.Result{
				Message: "OK",
				Status: api.StatusCode_GENERIC_OK,
			},
		}
		if err := srv.Send(resp); err != nil {
			return fmt.Errorf("could not send response thorugh stream: %v", err)
		}
	}
}

func (server) Store(ds api.DB_StoreServer) error {
	for {
		req, err := ds.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			res := &api.SlimQueryResult{
				Status: &api.Result{
					Status: api.StatusCode_READ_C_ERR,
					Message: "could not read from clients stream",
				},
			}
			if serr := ds.Send(res); serr != nil {
				return fmt.Errorf("could not send response: %v", err)
			}
		}

		resp := &api.SlimQueryResult{
			Domain: &api.DomainScore{
				Name: req.GetName(),
			},
			Status: &api.Result{
				Message: "OK",
				Status: api.StatusCode_STORE_S_OK,
			},
		}
		if err := ds.Send(resp); err != nil {
			return fmt.Errorf("could not send response thorugh stream: %v", err)
		}
	}
}
