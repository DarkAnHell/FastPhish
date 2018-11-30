package ct

import (
	"sync"
	"context"
	"log"
	"net/http"

	"github.com/DarkAnHell/FastPhish/api"
	ct "github.com/google/certificate-transparency-go"
	"github.com/google/certificate-transparency-go/client"
	"github.com/google/certificate-transparency-go/jsonclient"
	"github.com/pkg/errors"
)

// CT is a Certificate Transparency client.
type CT interface {
	// Handle makes listens for new domains in CT logs.
	Handle(ctx context.Context, domains chan<- api.Domain) error
}

// Client is a generic CT Logs consumer.
type Client struct {
	// cfg is the configuration of the CT Client.
	cfg *Config
	// clients stores each log URL with its client.
	// TODO: allow each client to set its start and size values.
	clients map[string]*client.LogClient
}

// Handle handles new domain names received from CT Logs.
func (c *Client) Handle(ctx context.Context, domains chan<- api.Domain) error {
	errs := make(chan error)
	var wg sync.WaitGroup

	for lurl, cli := range c.clients {
		wg.Add(1)
		go func(lurl string, cli *client.LogClient) {
			defer func() {
				wg.Done()
			}()

			re, err := cli.GetRawEntries(ctx, int64(c.cfg.Start), int64(c.cfg.Start+c.cfg.Size))
			if err != nil {
				errs <- errors.Wrap(err, "could not get raw entry")
				return
			}

			var index int
			for _, entry := range re.Entries {
				le, err := ct.LogEntryFromLeaf(int64(c.cfg.Start+index), &entry)
				if err != nil {
					// TODO: switch to a real logger.
					log.Printf("could not convert LeafEntry to LogEntry: %v", err)
					continue
				}
				if le.X509Cert == nil {
					continue
				}
				for _, v := range le.X509Cert.DNSNames {
					domains <- api.Domain{Name: v}
				}
				index++
			}
		}(lurl, cli)
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errs:
		return err
	case <-done:
		return nil
	}
}

// Option type allows customizing HTTP client.
type Option func(*http.Client)

// New sets up a CT Logs client.
func New(cfg *Config, opts ...Option) (*Client, error) {
	hc := &http.DefaultClient
	for _, opt := range opts {
		opt(*hc)
	}

	cli := &Client{
		cfg: cfg,
		clients: make(map[string]*client.LogClient, len(cfg.Logs)),
	}

	for _, lurl := range cfg.Logs {
		lc, err := client.New(lurl, *hc, jsonclient.Options{})
		if err != nil {
			return nil, errors.Wrap(err, "could not create client")
		}
		cli.clients[lurl] = lc
	}

	return cli, nil
}
