package datasource

import (
	"context"
	"net/http"

	"github.com/DarkAnHell/FastPhish/api"
)

// Datasource interface represents any source from which domains are taken.
type Datasource interface {
	// Request receives the data from the datasource.
	Request(context.Context, *http.Client, chan<- api.Domain) error
}
