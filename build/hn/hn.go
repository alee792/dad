package main

import (
	"context"
	"fmt"
	"net/http"

	dadhttp "github.com/alee792/dad/internal/http"
	"github.com/alee792/dad/internal/resolver"
	"github.com/alee792/dad/pkg/dad"
)

var s *dadhttp.Server

func init() {
	cfg := resolver.Config{
		Source: "hn",
		Dad: dad.Config{
			Order: 1,
		},
	}
	rsv := resolver.NewResolver(cfg)
	log := rsv.ResolveLogger()

	// Prime the Server.
	s = rsv.ResolveHTTP()
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
}

func Handler(w http.ResponseWriter, r *http.Request) {
	s.GetJoke()(w, r)
}

func main() {}
