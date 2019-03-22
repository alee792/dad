package getzit

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// RESTClient handles REST calls to ICHDJ.
type RESTClient struct {
	*http.Client
	Config Config
}

// NewRESTClient returns a client with sensible defaults.
func NewRESTClient(httpClient *http.Client, cfg Config) *RESTClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if cfg.Addr == "" {
		cfg.Addr = "https://icanhazdadjoke.com/"
	}
	c := &RESTClient{
		Client: httpClient,
		Config: cfg,
	}
	return c
}

// GetJoke from ICHDJ.
func (c *RESTClient) GetJoke(ctx context.Context) (string, error) {
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
