package joke

import (
	"net/http"

	"github.com/alee792/dad/pkg/now"
)

var s *now.Server

func init() {
	s = now.Init("joke", 1)
}

// Handler is an artifact for Now.
func Handler(w http.ResponseWriter, r *http.Request) {
	s.GetJoke()(w, r)
}
