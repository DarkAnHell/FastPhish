package db

import (
	"github.com/DarkAnHell/FastPhish/api"
)

// DB interface represents any db
type DB interface {
	Store(api.DomainScore) error
	GetScore(domain string) (score int, err error)
}
