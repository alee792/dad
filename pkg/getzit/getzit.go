// Package getzit is an API for https://icanhazdadjoke.com/ (ICHDJ).
package getzit

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Client handles API calls to https://icanhazdadjoke.com/.
type Client struct {
	*http.Client
	Config Config
}

// Config for Client.
type Config struct {
	Addr string
}

// NewClient returns a Client with sensible defaults.
func NewClient(httpClient *http.Client, cfg Config) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if cfg.Addr == "" {
		cfg.Addr = "https://icanhadadjoke.com/"
	}
	c := &Client{
		Client: httpClient,
		Config: cfg,
	}
	return c
}

// GetJoke from ICHDJ.
func (c *Client) GetJoke(ctx context.Context) (string, error) {
	resp, err := c.Get(c.Config.Addr)
	if err != nil {
		return "", errors.Wrap(err, "HTTP GET failed")
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "could not read response")
	}
	defer resp.Body.Close()
	return string(bb), nil
}
