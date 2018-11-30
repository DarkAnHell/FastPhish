package ingestor

import (
	"sync"
	"context"
	"net/http"

	"github.com/DarkAnHell/FastPhish/api/domain"
	"github.com/DarkAnHell/FastPhish/pkg/datasource"
)

// Ingestor collects data from the given source.
type Ingestor interface {
	// Ingest collects domains from datasources and reports them.
	Ingest(context.Context, chan<- domain.Domain, ...datasource.Datasource) error
}

// Default ingestor calls all datasources (one goroutine for each source)
// and reports results through domains channel.
//
// Default ingestor doesn't give priority to any source.
type Default struct {}

// Ingest collects domains from all datasources and reports them. It might report duplicate domains.
func (d Default) Ingest(ctx context.Context, out chan<- domain.Domain, sources ...datasource.Datasource) error {
	errs := make(chan error)
	done := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(len(sources))
	go func() {
		wg.Wait()
		close(done)
	}()

	for _, source := range sources {
		go func(s datasource.Datasource) {
			defer wg.Done()

			cli := &http.DefaultClient
			err := s.Request(ctx, *cli, out)
			if err != nil {
				errs <- err
				return
			}
		}(source)
	}
	select {
	case err := <- errs:
		return err
	case <-done:
		return nil
	}
}
