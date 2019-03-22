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
		// This is a no bueno and for demonstrative purposes only!
		// tr := &http.Transport{
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// }
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
	req, err := http.NewRequest(http.MethodGet, c.Config.Addr, nil)
	if err != nil {
		return "", errors.Wrap(err, "bad request")
	}
	req.Header.Set("Accept", "text/plain")
	resp, err := c.Do(req)
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
