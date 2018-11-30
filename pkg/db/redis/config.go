package redis

import (
	"time"
)

// Config for redis DB
type Config struct {
	// Listen location (i.e localhost:7070)
	Listen     string        `json:"Listen"`
	Password   string        `json:"Pass"`
	DB         int           `json:"DB"`
	Expiration time.Duration `json:"Expiration"`
}
