package redis

import (
	"io"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/go-redis/redis"
)

// Redis is the DB implementation for redis
type Redis struct {
}

// TODO: Replace by load
func (r Redis) Load(r io.Reader) error {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (r Redis) Store(api.Domain) error {

}

func (r Redis) GetScore(domain string) (score int, err error) {

}
