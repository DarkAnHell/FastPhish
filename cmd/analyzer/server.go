package main

import (
	"log"
	"fmt"
	"io"

	"github.com/DarkAnHell/FastPhish/api"
	lev "github.com/DarkAnHell/FastPhish/pkg/analyzer/levenshtein"
)

// Server exposes the analyzer API.
type Server struct {
	lev *lev.Levenshtein
	// TODO: config.
	against []string
}

// Analyze runs analysis logic against the given domains.
func (s Server) Analyze(d api.Analyzer_AnalyzeServer) error {
	for {
		req, err := d.Recv()
		if err != nil {
			return fmt.Errorf("could not receive domain: %v", err)
		}

		out := s.lev.Process(req.GetName(), s.against)
		var max uint32
		for _, score := range out {
			if score.Score > max {
				max = score.Score
			}
		}

		log.Printf("results for domain %s", req.GetName())
		for index, score := range out {
			log.Printf("got %v score against %s", score, s.against[index])
		}

		resp := &api.SlimQueryResult{
			Domain: &api.DomainScore{
				Name: req.GetName(),
				Score: max,
			},
			Status: &api.Result{
				Message: "ANALYZED",
				Status: api.StatusCode_ANALYZE_OK,
			},
		}
		if err := d.Send(resp); err != nil {
			return fmt.Errorf("could not send analysis result: %v", err)
		}
	}
}

//rpc Analyze(stream Domain) returns (stream SlimQueryResult) {}

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
