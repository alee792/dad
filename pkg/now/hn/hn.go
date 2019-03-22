package main

import (
	"net/http"

	"github.com/alee792/dad/pkg/now"
)

var s *now.Server

func init() {
	s = now.Init("hn", 2)
}

// Handler is an artifact for Now.
func Handler(w http.ResponseWriter, r *http.Request) {
	s.GetJoke()(w, r)
}

func main() {}
