package ingestor

import (
	"github.com/darkanhell/Fastphish/api/domain"
	"github.com/darkanhell/Fastphish/pkg/datasource"
)

// Ingestor collects data from the given source.
type Ingestor interface {
	Ingest(datasource.Datasource, chan<- domain.Domain) error
}
