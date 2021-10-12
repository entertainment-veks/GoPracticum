package shortener

import (
	"encoding/json"
	"go_practicum/internal/app/model"
	"go_practicum/internal/app/store"
	"go_practicum/internal/app/util"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	USERID_COOKIE_KEY = "userid"
)

type server struct {
	router  *mux.Router
	store   store.Store
	baseURL string
}

func newServer(store store.Store, baseURL string) *server {
	s := &server{
		router:  mux.NewRouter(),
		store:   store,
		baseURL: baseURL,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Handle("/", s.handleLinkCreate()).Methods(http.MethodPost)
	s.router.Handle("/api/shorten", s.handleLinkCreateJSON()).Methods(http.MethodPost)
	s.router.Handle("/{key}", s.handleLinkGet()).Methods(http.MethodGet)

	//s.router = s.authMiddleware(s.router)

	// s.router.Handle("/{key}", AuthMiddleware(GzipMiddleware(GetHandler(&service)))).Methods(http.MethodGet)
	// s.router.Handle("/", AuthMiddleware(GunzipMiddleware(PostHandler(&service)))).Methods(http.MethodPost)
	// s.router.Handle("/api/shorten", AuthMiddleware(GunzipMiddleware(PostJSONHandler(&service)))).Methods(http.MethodPost)
	// s.router.Handle("/user/urls", UserURLsHandler(&service))
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, err.Error())
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data string) {
	w.WriteHeader(code)
	if len(data) != 0 {
		w.Write([]byte(data))
	}
}

func (s *server) handleLinkCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		code, err := util.GenerateCode()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		userIdCookie, err := r.Cookie(USERID_COOKIE_KEY)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.store.Link().Create(&model.Link{
			Link:   string(body),
			Code:   code,
			UserID: userIdCookie.Value,
		})

		s.respond(w, r, http.StatusCreated, s.baseURL+"/"+code)
	}
}

func (s *server) handleLinkCreateJSON() http.HandlerFunc {
	type request struct {
		Link string `json:"url"`
	}

	type response struct {
		Output string `json:"result"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		code, err := util.GenerateCode()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		userIdCookie, err := r.Cookie(USERID_COOKIE_KEY)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		l := &model.Link{
			Link:   req.Link,
			Code:   code,
			UserID: userIdCookie.Value,
		}
		if err := s.store.Link().Create(l); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		rawResult := &response{
			s.baseURL + "/" + code,
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rawResult); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusCreated, "")
	}
}

func (s *server) handleLinkGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, err := s.store.Link().GetByCode(mux.Vars(r)["key"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Location", l.Link)
		s.respond(w, r, http.StatusTemporaryRedirect, "")
	}
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(USERID_COOKIE_KEY)

		if err == http.ErrNoCookie {
			newUserId, err := util.GenerateCode()
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			cookie = &http.Cookie{
				Name:  USERID_COOKIE_KEY,
				Value: newUserId,
			}
			http.SetCookie(w, cookie)
		}

		next.ServeHTTP(w, r)
	})
}
