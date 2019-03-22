package resolver

import (
	"github.com/alee792/dad/internal/http"
	"github.com/alee792/dad/pkg/dad"
	"github.com/alee792/dad/pkg/dad/storage/json"
	"github.com/alee792/dad/pkg/getzit"
	"github.com/alee792/dad/pkg/hn"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// Resolver is a lazy loading dependency resolver.
type Resolver struct {
	HTTP   *http.Server
	Chain  *dad.Chain
	GQL    *getzit.GraphQLClient
	REST   *getzit.RESTClient
	Logger *zap.SugaredLogger
	HN     *hn.Client
	Config Config
}

// Config for Resolver.
type Config struct {
	Source string
	HTTP   http.Config
	Dad    dad.Config
	HN     hn.Config
	GQL    getzit.Config
	REST   getzit.Config
}

// NewResolver is created from a Config.
func NewResolver(cfg Config) *Resolver {
	r := &Resolver{
		Config: cfg,
	}
	return r
}

// ResolveHTTP service.
func (r *Resolver) ResolveHTTP() *http.Server {
	if r.HTTP == nil {
		s := http.NewServer(
			r.ResolveJoker(),
			r.ResolveChain(),
			chi.NewRouter(),
			r.ResolveLogger(),
			r.Config.HTTP,
		)
		r.HTTP = s
	}
	return r.HTTP
}

// ResolveJoker service.
func (r *Resolver) ResolveJoker() http.Joker {
	switch r.Config.Source {
	case "hn":
		return r.ResolveHN()
	case "joke":
		return r.ResolveJokeREST()
	default:
		r.ResolveLogger().Fatal("%s is not a valid source", r.Config.Source)
	}
	return nil
}

// ResolveHN service.
func (r *Resolver) ResolveHN() *hn.Client {
	if r.HN == nil {
		r.HN = hn.NewClient(nil, r.Config.HN)
	}
	return r.HN
}

// ResolveJokeGQL service.
func (r *Resolver) ResolveJokeGQL() *getzit.GraphQLClient {
	if r.GQL == nil {
		r.GQL = getzit.NewGraphQLClient(r.Config.GQL)
	}
	return r.GQL
}

// ResolveJokeREST service.
func (r *Resolver) ResolveJokeREST() *getzit.RESTClient {
	if r.REST == nil {
		r.REST = getzit.NewRESTClient(nil, r.Config.REST)
	}
	return r.REST
}

// ResolveChain service.
func (r *Resolver) ResolveChain() *dad.Chain {
	if r.Chain == nil {
		r.Chain = dad.NewChain(
			&json.Store{},
			r.Config.Dad,
		)
	}
	return r.Chain
}

// ResolveLogger service.
func (r *Resolver) ResolveLogger() *zap.SugaredLogger {
	if r.Logger == nil {
		r.Logger = zap.NewExample().Sugar()
	}
	return r.Logger
}
