package shortener

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"go_practicum/internal/app/model"
	"go_practicum/internal/app/store"
	"go_practicum/internal/app/util"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	USERID_COOKIE_KEY         = "userid"
	USERID_CONTEXT_KEY ctxKey = iota
)

type ctxKey int8

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
	s.router.Handle("/user/urls", s.handleUserLinks()).Methods(http.MethodGet)

	s.router.Use(s.authMiddleware)
	//s.router.Use(s.gzipMiddleware)
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

		if err := s.store.Link().Create(&model.Link{
			Link:   string(body),
			Code:   code,
			UserID: r.Context().Value(USERID_CONTEXT_KEY).(string),
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

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

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		l := &model.Link{
			Link:   req.Link,
			Code:   code,
			UserID: r.Context().Value(USERID_CONTEXT_KEY).(string),
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
			return
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

func (s *server) handleUserLinks() http.HandlerFunc {
	type userLink struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		links, err := s.store.Link().GetAllByUserID(r.Context().Value(USERID_CONTEXT_KEY).(string))
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNoContent, err)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		userLinks := []userLink{}
		for _, current := range links {
			userLinks = append(userLinks, userLink{
				ShortURL:    s.baseURL + "/" + current.Code,
				OriginalURL: current.Link,
			})
		}

		if err := json.NewEncoder(w).Encode(userLinks); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, "")
	}
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(USERID_COOKIE_KEY)

		var newUserId string
		if err == http.ErrNoCookie {
			newUserId, err = util.GenerateCode()
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			cookie := &http.Cookie{
				Name:  USERID_COOKIE_KEY,
				Value: newUserId,
			}
			http.SetCookie(w, cookie)
		} else {
			newUserId = cookie.Value
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), USERID_CONTEXT_KEY, newUserId)))
	})
}

func (s *server) gzipMiddleware(next http.Handler) http.Handler {
	type gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			result, err := gzip.NewReader(r.Body)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			r.Body = result
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
