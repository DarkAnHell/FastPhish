package datasource

import (
	"context"
	"io"
	"net/http"
)

// Datasource interface represents any source from which domains are taken.
type Datasource interface {
	// Request receives the data from the datasource.
	Request(context.Context, *http.Client) (io.Reader, error)
}
