package shortener

import (
	"github.com/gorilla/mux"
	"go_practicum/app/config"
	"go_practicum/app/handler"
	"go_practicum/app/store"
	"net/http"
)

type router struct {
	router  *mux.Router
	store   store.Store
	baseURL string
}

func newRouter(store store.Store, cfg config.Config) *router {
	r := &router{
		router:  mux.NewRouter(),
		store:   store,
		baseURL: cfg.BaseURL,
	}

	r.configureRouter(store, cfg)

	return r
}

func (s *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *router) configureRouter(store store.Store, cfg config.Config) {
	s.router.Handle("/", handler.HandleLinkCreate(store, cfg)).Methods(http.MethodPost)
	s.router.Handle("/api/shorten", handler.HandleLinkCreateJson(store, cfg)).Methods(http.MethodPost)
	s.router.Handle("/{key}", handler.HandleLinkGet(store)).Methods(http.MethodGet)
	s.router.Handle("/ping", handler.HandlePing(store)).Methods(http.MethodGet)
	s.router.Handle("/user/urls", handler.HandleGetUserLinks(store, cfg)).Methods(http.MethodGet)
	s.router.Handle("/api/shorten/batch", handler.HandleLinkCreateAll(store, cfg)).Methods(http.MethodPost)
	s.router.Handle("/api/user/urls", handler.HandleLinkDelete(store)).Methods(http.MethodDelete)

	s.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	s.router.Use(handler.AuthMiddleware)
	s.router.Use(handler.GzipMiddleware)
}
