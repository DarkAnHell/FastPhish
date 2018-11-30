package redis

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/go-redis/redis"
)

// Redis is the DB implementation for redis
type Redis struct {
	client *redis.Client

	expirationTime time.Duration
}

func (R *Redis) Load(r io.Reader) error {

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return err
	}

	R.expirationTime = cfg.Expiration

	R.client = redis.NewClient(&redis.Options{
		Addr:     cfg.Listen,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err = R.client.Ping().Result()

	return err
}

func (r *Redis) Store(d api.DomainScore) error {

	return r.client.Set(d.Name, d.Score, r.expirationTime).Err()
}

func (r *Redis) GetScore(domain string) (score int, err error) {

	val, err := r.client.Get(domain).Result()
	if err != nil {
		return -1, err
	}

	score, err = strconv.Atoi(val)
	if err != nil {
		return -1, err
	}

	return
}
