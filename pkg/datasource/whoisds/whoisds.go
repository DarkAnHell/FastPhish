package whoisds

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DarkAnHell/FastPhish/api"
)

// Whoisds uses whoisds.com datasets as data source.
type Whoisds struct {
	date time.Time
}

const (
	baseURL    = "https://whoisds.com//whois-database/newly-registered-domains/"
	baseSuffix = "/nrd"
)

// Request downloads the data from whoisds.com and extracts it.
func (w Whoisds) Request(ctx context.Context, cli *http.Client, domains chan<- api.Domain) error {
	fname := fmt.Sprintf("%v-%d-%v.zip", w.date.Year(), w.date.Month(), w.date.Day())
	encoded := base64.StdEncoding.EncodeToString([]byte(fname))
	fullURL := fmt.Sprintf("%s%s%s", baseURL, encoded, baseSuffix)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}
	req = req.WithContext(ctx)

	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("could not request %s: %v", fullURL, err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return fmt.Errorf("could not read ZIP response: %v", err)
	}

	for _, zf := range zr.File {
		f, err := zf.Open()
		if err != nil {
			return fmt.Errorf("could not read %s file: %v", zf.Name, err)
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			domains <- api.Domain{Name: sc.Text()}
		}
		if err := sc.Err(); err != nil {
			return fmt.Errorf("could not scan: %v", err)
		}
	}

	return nil
}

// New creates a Whoisds data ingestor using the given date's dataset.
func New(date time.Time) *Whoisds {
	return &Whoisds{date: date}
}
