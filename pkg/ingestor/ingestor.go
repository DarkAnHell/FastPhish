package ingestor

import (
	"github.com/DarkAnHell/FastPhish/api/domain"
	"github.com/DarkAnHell/FastPhish/pkg/datasource"
)

// Ingestor collects data from the given source.
type Ingestor interface {
	Ingest(datasource.Datasource, chan<- domain.Domain) error
}
