package main

import (
	"fmt"
	"io"
	"log"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/db"
)

type server struct {
	DB db.DB
}

func (s server) GetDomainsScore(srv api.DB_GetDomainsScoreServer) error {
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil || req == nil {
			res := &api.SlimQueryResult{
				Status: &api.Result{
					Status:  api.StatusCode_READ_C_ERR,
					Message: "could not read from clients stream",
				},
			}
			if serr := srv.Send(res); serr != nil {
				return fmt.Errorf("could not send response: %v", err)
			}
		}

		score, err := s.DB.GetScore(req.GetName())
		if err != nil && err != db.ErrDBNotFound {
			return fmt.Errorf("could not get score from DB: %v", err)
		}
		log.Printf("Read result for domain: %s", req.GetName())

		resp := &api.SlimQueryResult{
			Domain: &api.DomainScore{
				Name:  req.GetName(),
				Score: uint32(score),
			},
			Status: &api.Result{
				Message: "OK",
				Status:  api.StatusCode_GET_SCORE_S_OK,
			},
		}
		if err := srv.Send(resp); err != nil {
			return fmt.Errorf("could not send response thorugh stream: %v", err)
		}
	}
}

func (s server) Store(ds api.DB_StoreServer) error {
	for {
		req, err := ds.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			res := &api.SlimQueryResult{
				Status: &api.Result{
					Status:  api.StatusCode_READ_C_ERR,
					Message: "could not read from clients stream",
				},
			}
			if serr := ds.Send(res); serr != nil {
				return fmt.Errorf("could not send response: %v", err)
			}
		}

		if err := s.DB.Store(*req); err != nil {
			return fmt.Errorf("could not store in DB: %v", err)
		}

		resp := &api.SlimQueryResult{
			Domain: &api.DomainScore{
				Name: req.GetName(),
			},
			Status: &api.Result{
				Message: "OK",
				Status:  api.StatusCode_STORE_S_OK,
			},
		}
		if err := ds.Send(resp); err != nil {
			return fmt.Errorf("could not send response thorugh stream: %v", err)
		}
	}
}
