package getzit

import (
	"context"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

// GraphQLClient handles  Graphql calls to ICHDJ.
type GraphQLClient struct {
	*graphql.Client
	Config Config
}

// NewGraphQLClient returns a client with sensible defaults.
func NewGraphQLClient(cfg Config) *GraphQLClient {
	if cfg.Addr == "" {
		cfg.Addr = "https://icanhazdadjoke.com/graphql"
	}
	gql := graphql.NewClient(cfg.Addr)
	c := &GraphQLClient{
		Client: gql,
		Config: cfg,
	}
	return c
}

// GetJoke from this wonderful service.
func (c *GraphQLClient) GetJoke(ctx context.Context) (string, error) {
	type Data struct {
		Joke struct {
			ID        string
			Joke      string
			Permalink string
		}
	}
	req := graphql.NewRequest(`
	query {
		joke {joke}
	}
	`)
	var resp Data
	if err := c.Run(ctx, req, &resp); err != nil {
		return "", errors.Wrap(err, "GraphQL request failed")
	}
	return resp.Joke.Joke, nil
}
