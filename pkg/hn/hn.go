// Package hn is where Dad tries to read HackerNews
// and tries to explain an article he read back to you.
package hn

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client handles  calls to ICHDJ.
type Client struct {
	Stories *Stories
	*http.Client
	Config Config
}

// Config for any client.
type Config struct {
	ItemAddr string
	TopAddr  string
}

// NewClient returns a client with sensible defaults.
func NewClient(httpClient *http.Client, cfg Config) *Client {
	if httpClient == nil {
		// This is a no bueno and for demonstrative purposes only!
		// tr := &http.Transport{
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// }
		httpClient = http.DefaultClient
	}
	if cfg.TopAddr == "" {
		cfg.TopAddr = "https://hacker-news.firebaseio.com/v0/topstories.json"
	}
	if cfg.ItemAddr == "" {
		cfg.ItemAddr = "https://hacker-news.firebaseio.com/v0/item/"
	}
	c := &Client{
		Stories: &Stories{},
		Client:  httpClient,
		Config:  cfg,
	}
	rand.Seed(time.Now().Unix())
	return c
}

// Stories provided by HN API.
type Stories struct {
	IDs []int
	eat time.Time
}

// Item provided by HN API.
type Item struct {
	Title string `json:"title"`
}

// Get titles from HackerNews.
func (c *Client) Get(ctx context.Context) (string, error) {
	ss, err := c.getStories(ctx)
	if err != nil {
		return "", err
	}
	id := ss.IDs[rand.Intn(len(ss.IDs))]
	u := fmt.Sprintf("%s/%d.json", c.Config.ItemAddr, id)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	var itm *Item
	if err := json.NewDecoder(resp.Body).Decode(&itm); err != nil {
		return "", errors.Wrap(err, "JSON decode failed")
	}
	defer resp.Body.Close()
	return itm.Title, nil
}

func (c *Client) getStories(ctx context.Context) (*Stories, error) {
	if c.Stories.eat.After(time.Now()) {
		return c.Stories, nil
	}
	resp, err := http.Get(c.Config.TopAddr)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP GET failed")
	}
	var itms []int
	if err := json.NewDecoder(resp.Body).Decode(&itms); err != nil {
		return nil, errors.Wrap(err, "JSON decode failed")
	}
	defer resp.Body.Close()
	if len(itms) < 1 {
		return nil, errors.New("could not retrieve items")
	}
	c.Stories = &Stories{
		IDs: itms,
		eat: time.Now().Add(time.Hour),
	}
	return c.Stories, nil
}
