// Package now accesses internal packages
// to enable Zeit Now deployments.
package now

import (
	"context"
	"fmt"

	dadhttp "github.com/alee792/dad/internal/http"
	"github.com/alee792/dad/internal/resolver"
	"github.com/alee792/dad/pkg/dad"
)

// Server wraps up internals.
type Server struct {
	*dadhttp.Server
}

// Init a Dad in Now fashion.
func Init(src string, order int) *Server {
	cfg := resolver.Config{
		Source: src,
		Dad: dad.Config{
			Order: order,
		},
	}
	rsv := resolver.NewResolver(cfg)
	log := rsv.ResolveLogger()

	// Prime the Server.
	s := rsv.ResolveHTTP()
	path := fmt.Sprintf("./bin/%dgrams-%s.json", cfg.Dad.Order, cfg.Source)
	go func() {
		if err := s.WarmUp(context.Background(), 200); err != nil {
			log.Warnw("could not seed grams from API", "err", err)
		}
		if err := s.Chainer.Save(path); err != nil {
			log.Warnw("could not save grams", "path", path, "err", err)
		}
	}()
	loaded, err := s.Chainer.Load(path)
	if err != nil {
		log.Warnw("could not load grams", "path", path, "err", err)
	}
	log.Infof("%d grams loaded from cache", loaded)
	return &Server{s}
}
