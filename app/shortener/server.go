package shortener

import (
	"github.com/gorilla/mux"
	"go_practicum/app/config"
	"go_practicum/app/handler"
	"go_practicum/app/store"
	"net/http"
)

type server struct {
	router  *mux.Router
	store   store.Store
	baseURL string
}

func newServer(store store.Store, cfg config.Config) *server {
	s := &server{
		router:  mux.NewRouter(),
		store:   store,
		baseURL: cfg.BaseURL,
	}

	s.configureRouter(store, cfg)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter(store store.Store, cfg config.Config) {
	s.router.Handle("/", handler.HandleLinkCreate(store, cfg)).Methods(http.MethodPost)
	s.router.Handle("/api/shorten", handler.HandleLinkCreateJSON(store, cfg)).Methods(http.MethodPost)
	s.router.Handle("/{key}", handler.HandleLinkGet(store)).Methods(http.MethodGet)
	s.router.Handle("/ping", handler.HandlePing(store)).Methods(http.MethodGet)
	s.router.Handle("/user/urls", handler.HandleGetUserLinks(store, cfg)).Methods(http.MethodGet)
	s.router.Handle("/api/shorten/batch", handler.HandleLinkCreateAll(store, cfg)).Methods(http.MethodPost)
	s.router.Handle("/api/user/urls", handler.HandleLinkDelete(store)).Methods(http.MethodDelete)

	s.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	s.router.Use(handler.AuthMiddleware)
	s.router.Use(handler.GzipMiddleware)
}
