package https

import (
	"net/http"

	"github.com/fahmifan/ktpready"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/thedevsaddam/renderer"
)

type Map map[string]interface{}

const (
	DefaultConcurrent           = 10
	DefaultPort                 = "8080"
	DefaultEnableFathomAnalytic = false
)

type Server struct {
	port                 string
	mux                  *chi.Mux
	server               *http.Server
	render               *renderer.Render
	nameChecker          *ktpready.NameChecker
	concurrent           int
	enableFathomAnalytic bool
}

type ServerOpt func(s *Server)

// Groups of server options
var ServerOpts = struct {
	WithConcurrent      func(c int) ServerOpt
	WithPort            func(port string) ServerOpt
	WithFathomAnalytics func(enable bool) ServerOpt
}{
	WithConcurrent: func(c int) ServerOpt {
		return func(s *Server) {
			if c < 1 {
				return
			}
			s.concurrent = c
		}
	},
	WithPort: func(port string) ServerOpt {
		return func(s *Server) {
			if port == "" {
				return
			}
			s.port = port
		}
	},
	WithFathomAnalytics: func(enable bool) ServerOpt {
		return func(s *Server) {
			s.enableFathomAnalytic = enable
		}
	},
}

func NewServer(nameChecker *ktpready.NameChecker, opts ...ServerOpt) *Server {
	server := &Server{
		port: DefaultPort,
		render: renderer.New(renderer.Options{
			ParseGlobPattern: "https/view/*.html",
			LeftDelim:        "[[",
			RightDelim:       "]]",
		}),
		nameChecker:          nameChecker,
		concurrent:           DefaultConcurrent,
		enableFathomAnalytic: DefaultEnableFathomAnalytic,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Run() error {
	s.routes()
	s.server = &http.Server{Addr: ":" + s.port, Handler: s.mux}

	log.Info().Msgf("run server at localhost:%s", s.port)
	return s.server.ListenAndServe()
}

func (s *Server) routes() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Throttle(s.concurrent))

	ktp := KTP{s}

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		s.render.JSON(w, http.StatusOK, Map{"ping": "pong"})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		s.render.HTML(w, http.StatusOK, "index.html", Map{
			"EnableFathom": s.enableFathomAnalytic,
		})
	})
	r.Post("/ktp", ktp.create())

	s.mux = r
}
