// Package http is an http implementation of Dad.
package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"unicode"

	"github.com/pkg/errors"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// Joker is where `dad` gets his jokes from.
// Could be his dad, your Uncle, or a website called icanhazdadjoke.
// Who knows?
type Joker interface {
	Get(context.Context) (string, error)
}

// Chainer is the API for a Markov Chain.
type Chainer interface {
	Read(context.Context, io.Reader)
	Generate(context.Context) string
	Save(path string) error
	Load(path string) (int, error)
}

// Server manages the dependencies and operations of a Dad service.
type Server struct {
	Joker   Joker
	Chainer Chainer
	Router  chi.Router
	Logger  *zap.SugaredLogger
	Config  Config
}

// Config for a Dad Server.
type Config struct {
	Addr           string
	MaxConcurrency int64
}

// NewServer with appropriate validation and defaults.
func NewServer(j Joker, c Chainer, r chi.Router, l *zap.SugaredLogger, cfg Config) *Server {
	if l == nil {
		l = zap.NewExample().Sugar()
	}
	if r == nil {
		r = chi.NewRouter()
	}
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.MaxConcurrency < 1 {
		cfg.MaxConcurrency = 5
	}
	s := &Server{
		Joker:   j,
		Chainer: c,
		Router:  r,
		Logger:  l,
		Config:  cfg,
	}
	s.Routes()
	return s
}

// Routes registers paths with handlers.
func (s *Server) Routes() {
	s.Router.Use(middleware.Heartbeat("/health"))
	s.Router.Use(middleware.Logger)
	s.Router.Get("/joke", s.GetJoke())
	s.Router.Get("/realjoke", s.GetRealJoke())
}

// GetJoke generates a joke from the Chainer.
func (s *Server) GetJoke() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var j string
		for j == "" {
			j = s.Chainer.Generate(r.Context())
		}
		bb := []byte(j)
		bb[0] = byte(unicode.ToUpper(rune(bb[0])))
		fmt.Fprintf(w, gnarlyTemplate, string(bb))
	}
}

// GetRealJoke directly from a Joker.
func (s *Server) GetRealJoke() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		joke, err := s.Joker.Get(r.Context())
		if err != nil {
			http.Error(
				w,
				errors.Wrap(err, "failed to retrieve a real joke").Error(),
				http.StatusInternalServerError,
			)
		}
		fmt.Fprintf(w, joke)
	}
}

// gnarlyTemplate for Now.
var gnarlyTemplate = `
<html>
	<head>
        
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Have you heard this one?</title>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.7.4/css/bulma.min.css">
		<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.8.1/css/all.css" integrity="sha384-50oBUHEmvpQ+1lW4y57PTFmhCaXp0ML5d60M1M7uH2+nqUivzIebhndOJK28anvf" crossorigin="anonymous">
		<link href="https://fonts.googleapis.com/css?family=Slabo+27px" rel="stylesheet">
		<style>
			p {
				font-family: 'Slabo 27px', serif;
			}
			.valign{
				position: relative;
				top: 33%%;
				transform: translateY(-67%%);
			}
		</style>
	</head>

		<body>
		<div class="container valign">
		<section class="section has-background-light">
			<div class="column is-centered is-vcentered">
				<div class="container">
					<p>%s</p>
				</div>
				<div class="container">
					<span class="icon">
						<a href="https://github.com/alee792/dad" title="Source" class="fab fa-github"></a>
					</span>
					<span class="icon has-text-info">
						<a href="https://dad-git-master.alee792.now.sh" title="Zeit Now" class="fas fa-bolt"></a>
					</span>
				</div>
			</div>    
		</section>
		</div>
	</body>
</html>
`
