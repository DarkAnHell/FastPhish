package ct

import (
	"context"
	"log"
	"net/http"

	"github.com/DarkAnHell/FastPhish/api/domain"
	ct "github.com/google/certificate-transparency-go"
	"github.com/google/certificate-transparency-go/client"
	"github.com/google/certificate-transparency-go/jsonclient"
	"github.com/pkg/errors"
)

// CT is a Certificate Transparency client.
type CT interface {
	// Handle makes listens for new domains in CT logs.
	Handle(ctx context.Context, start, size int, domains chan<- domain.Domain) error
	// Stop allows a client to stop listening for new domains.
	Stop()
}

// Client is a generic CT Logs consumer.
type Client struct {
	lc *client.LogClient
	// done channel allows stopping the client.
	done chan struct{}
}

// Handle handles new domain names received from CT Logs.
func (c *Client) Handle(ctx context.Context, start, size int, domains chan<- domain.Domain) error {
	if start < 0 {
		return errors.New("start value must be higher or equal to 0")
	}
	if size <= 0 {
		return errors.New("size must be higher or equal to 0")
	}

	for {
		select {
		case <-c.done:
			return nil
		default:
		}

		re, err := c.lc.GetRawEntries(ctx, int64(start), int64(start+size))
		if err != nil {
			return errors.Wrap(err, "could not get raw entry")
		}

		var index int
		for _, entry := range re.Entries {
			le, err := ct.LogEntryFromLeaf(int64(index), &entry)
			if err != nil {
				// TODO: switch to a real logger.
				log.Printf("could not convert LeafEntry to LogEntry: %v", err)
				continue
			}

			if le.X509Cert == nil {
				continue
			}
			for _, v := range le.X509Cert.DNSNames {
				domains <- domain.Domain{Name: v}
			}

			index++
		}
	}
}

// Stop finishes the client.
func (c *Client) Stop() {
	close(c.done)
}

// Option type allows customizing HTTP client.
type Option func(*http.Client)

// New sets up a CT Logs client.
func New(ctlog string, opts ...Option) (*Client, error) {
	hc := &http.DefaultClient
	for _, opt := range opts {
		opt(*hc)
	}

	lc, err := client.New(ctlog, *hc, jsonclient.Options{})
	if err != nil {
		return nil, errors.Wrap(err, "could not create client")
	}

	done := make(chan struct{})
	return &Client{
		lc:   lc,
		done: done,
	}, nil
}
