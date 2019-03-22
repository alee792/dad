package resolver

import (
	"github.com/alee792/dad/internal/http"
	"github.com/alee792/dad/pkg/dad"
	"github.com/alee792/dad/pkg/dad/storage/json"
	"github.com/alee792/dad/pkg/getzit"
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
func (r *Resolver) ResolveHTTP() *http.Server {
	if r.HTTP == nil {
		s := http.NewServer(
			r.ResolveREST(),
			r.ResolveChain(),
			chi.NewRouter(),
			r.ResolveLogger(),
			r.Config.HTTP,
		)
		r.HTTP = s
	}
	return r.HTTP
}

// ResolveGQL service.
func (r *Resolver) ResolveGQL() *getzit.GraphQLClient {
	if r.GQL == nil {
		r.GQL = getzit.NewGraphQLClient(r.Config.GQL)
	}
	return r.GQL
}

// ResolveREST service.
func (r *Resolver) ResolveREST() *getzit.RESTClient {
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
