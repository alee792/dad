package resolver

import (
	"github.com/alee792/dad/internal/http"
	"github.com/alee792/dad/pkg/dad"
	"github.com/alee792/dad/pkg/getzit"
)

// Resolver for a dad service.
type Resolver struct {
	HTTP   *http.Server
	Dad    *dad.Chain
	GQL    *getzit.GraphQLClient
	REST   *getzit.RESTClient
	Config Config
}

// Config for Resolver.
type Config struct {
	HTTP http.Config
	Dad  dad.Config
	GQL  getzit.Config
	REST getzit.Config
}

// NewResolver is created from a Config.
func NewResolver(cfg Config) *Resolver {
	r := &Resolver{
		Config: cfg,
	}
	return r
}


// ResolveHTTP service.
func ResolveHTTTP() *http.Server P
