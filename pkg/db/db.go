package db

import (
	"errors"

	"github.com/DarkAnHell/FastPhish/api"
)

// DB interface represents any db
type DB interface {
	Store(api.DomainScore) error
	GetScore(domain string) (score int, err error)
}

var ErrDBNotFound = errors.New("key wasn't found on DB")
