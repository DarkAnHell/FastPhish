package db

import (
	"io"

	"github.com/DarkAnHell/FastPhish/api"
)

// DB interface represents any db
type DB interface {
	Store(api.Domain) error

	GetScore(domain string) (score int, err error)

	Load(r io.Reader) error
}
