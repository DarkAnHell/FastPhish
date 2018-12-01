package ct

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"io"
)

// Config sets CT Client behavior.
type Config struct {
	// Start is the CT Log first node to read.
	Start int `json:"start"`
	// Size is the number of logs to read.
	Size int `json:"size"`
	// Logs stores the list of CT Logs to use.
	Logs []string `json:"logs"`
	// Concurrency is the number of concurrent goroutines to use while
	// reading the CT Log.
	Concurrency int `json:"concurrency"`
}

// Load allows reading configuration from any type that implements the
// io.Reader interface, like files, network shares, etc.
func Load(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	if cfg.Start < 0 {
		return nil, errors.New("start index must be higher or equal to 0")
	}
	if cfg.Size < 1 {
		return nil, errors.New("size must be higher or equal to 0")
	}
	if cfg.Concurrency < 1 {
		return nil, errors.New("at least one goroutine is needed")
	}
	if len(cfg.Logs) < 1 {
		return nil, errors.New("at least one CT Log URL is needed")
	}
	if cfg.Concurrency > len(cfg.Logs) {
		cfg.Concurrency = len(cfg.Logs)
	}

	return &cfg, nil
}
