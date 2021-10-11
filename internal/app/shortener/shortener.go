package shortener

import (
	"fmt"
	"go_practicum/internal/app/store"
	"net/http"

	"github.com/gorilla/mux"
)

type Shortener struct {
	config *Config
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *Shortener {
	return &Shortener{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *Shortener) Start() error {
	s.configureRouter()
	if err := s.configureStore(); err != nil {
		return err
	}

	fmt.Println("Starting server")
	return http.ListenAndServe(s.config.ServerAddress, s.router)
}

func (s *Shortener) configureRouter() {
	// s.router.Handle("/{key}", AuthMiddleware(GzipMiddleware(GetHandler(&service)))).Methods(http.MethodGet)
	// s.router.Handle("/", AuthMiddleware(GunzipMiddleware(PostHandler(&service)))).Methods(http.MethodPost)
	// s.router.Handle("/api/shorten", AuthMiddleware(GunzipMiddleware(PostJSONHandler(&service)))).Methods(http.MethodPost)
	// s.router.Handle("/user/urls", UserURLsHandler(&service))
}

func (s *Shortener) configureStore() error {
	st := store.New(s.config.StoreConfig)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st
	return nil
}
