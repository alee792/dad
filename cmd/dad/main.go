package main

import (
	"context"
	"net/http"

	"github.com/alee792/dad/internal/resolver"
)

func main() {
	cfg := resolver.Config{}
	r := resolver.NewResolver(cfg)
	s := r.ResolveHTTP()
	go s.WarmUp(context.Background(), 200)
	r.ResolveLogger().Infof("Serving on %s", s.Config.Addr)
	http.ListenAndServe(s.Config.Addr, s.Router)
}
