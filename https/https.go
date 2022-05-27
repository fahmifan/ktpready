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

type Server struct {
	Port        string
	mux         *chi.Mux
	server      *http.Server
	render      *renderer.Render
	NameChecker *ktpready.NameChecker
}

func NewServer(port string, nameChecker *ktpready.NameChecker) *Server {
	return &Server{
		Port: port,
		render: renderer.New(renderer.Options{
			ParseGlobPattern: "https/view/*.html",
			LeftDelim:        "[[",
			RightDelim:       "]]",
		}),
		NameChecker: nameChecker,
	}
}

func (s *Server) Run() error {
	s.routes()
	s.server = &http.Server{Addr: ":" + s.Port, Handler: s.mux}

	log.Info().Msgf("run server at localhost:%s", s.Port)
	return s.server.ListenAndServe()
}

func (s *Server) routes() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	ktp := KTP{s}

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		s.render.JSON(w, http.StatusOK, Map{"ping": "pong"})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		s.render.HTML(w, http.StatusOK, "index.html", Map{})
	})
	r.Post("/ktp", ktp.create())

	s.mux = r
}
