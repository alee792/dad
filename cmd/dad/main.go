package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/alee792/dad/pkg/dad"

	"github.com/alee792/dad/internal/resolver"
	"github.com/peterbourgon/ff"
)

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		order = fs.Int("n", 2, "the order (size) of an n-gram")
	)
	ff.Parse(fs, os.Args[1:])

	cfg := resolver.Config{
		Dad: dad.Config{
			Order: *order,
		},
	}
	r := resolver.NewResolver(cfg)
	log := r.ResolveLogger()

	// Prime the Server.
	s := r.ResolveHTTP()
	path := fmt.Sprintf("./bin/%dgrams.json", *order)
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

	log.Infof("Serving on %s", s.Config.Addr)
	http.ListenAndServe(s.Config.Addr, s.Router)
}
