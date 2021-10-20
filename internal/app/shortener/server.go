package shortener

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"go_practicum/internal/app/model"
	"go_practicum/internal/app/store"
	"go_practicum/internal/app/util"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	userIDCookieKey         = "userid"
	userIDContextKey ctxKey = iota
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
	s.router.Handle("/ping", s.handlePing()).Methods(http.MethodGet)
	s.router.Handle("/user/urls", s.handleUserLinks()).Methods(http.MethodGet)
	s.router.Handle("/api/shorten/batch", s.handleLinkCreateAll()).Methods(http.MethodPost)

	s.router.Use(s.authMiddleware)
	s.router.Use(s.gzipMiddleware)
}

func (s *server) error(w http.ResponseWriter, code int, err error) {
	s.respond(w, code, err.Error())
}

func (s *server) respond(w http.ResponseWriter, code int, data string) {
	w.WriteHeader(code)
	if len(data) != 0 {
		w.Write([]byte(data))
	}
}


func (s *server) respondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	if err := json.NewEncoder(w).Encode(body); err != nil {
		s.error(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	s.respond(w, statusCode, "")
}

func (s *server) handleLinkCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		code, err := util.GenerateCode()
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		userID, ok := r.Context().Value(userIDContextKey).(string)
		if !ok {
			s.error(w, http.StatusInternalServerError, errors.New("user id is not provided"))
		}

		l := &model.Link{
			Link:   string(body),
			Code:   code,
			UserID: userID,
		}
		if err := l.Validate(); err != nil {
			s.error(w, http.StatusUnprocessableEntity, err)
			return
		}
		if err := s.store.Link().Create(l); err != nil {
			if errors.Is(err, store.ErrConflict) {
				s.error(w, http.StatusConflict, err)
				return
			}
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusCreated, s.baseURL+"/"+code)
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
			s.error(w, http.StatusBadRequest, err)
			return
		}

		code, err := util.GenerateCode()
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		userID, ok := r.Context().Value(userIDContextKey).(string)
		if !ok {
			s.error(w, http.StatusInternalServerError, errors.New("user id is not provided"))
		}

		l := &model.Link{
			Link:   req.Link,
			Code:   code,
			UserID: userID,
		}
		if err := l.Validate(); err != nil {
			s.error(w, http.StatusUnprocessableEntity, err)
			return
		}
		if err := s.store.Link().Create(l); err != nil {
			if errors.Is(err, store.ErrConflict) {
				s.error(w, http.StatusConflict, err)
				return
			}
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		rawResult := &response{
			s.baseURL + "/" + code,
		}
		s.respondJSON(w, http.StatusCreated, rawResult)
	}
}

func (s *server) handleLinkCreateAll() http.HandlerFunc {
	type requestElem struct {
		CorrelationID string `json:"correlation_id"`
		Link          string `json:"original_url"`
	}

	type responseElem struct {
		CorrelationID string `json:"correlation_id"`
		Link          string `json:"short_url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := []*requestElem{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		userID, ok := r.Context().Value(userIDContextKey).(string)
		if !ok {
			s.error(w, http.StatusInternalServerError, errors.New("user id is not provided"))
		}

		var ls []*model.Link
		for _, current := range req {
			l := &model.Link{
				Link:   current.Link,
				Code:   current.CorrelationID,
				UserID: userID,
			}
			if err := l.Validate(); err != nil {
				s.error(w, http.StatusUnprocessableEntity, err)
				return
			}
			ls = append(ls, l)
		}

		if err := s.store.Link().CreateAll(ls); err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		resp := []*responseElem{}
		for _, current := range ls {
			resp = append(resp, &responseElem{
				CorrelationID: current.Code,
				Link:          s.baseURL + "/" + current.Code,
			})
		}

		s.respondJSON(w, http.StatusCreated, resp)
	}
}

func (s *server) handleLinkGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := mux.Vars(r)["key"]
		l, err := s.store.Link().GetByCode(code)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		w.Header().Set("Location", l.Link)
		s.respond(w, http.StatusTemporaryRedirect, "")
	}
}

func (s *server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.store != nil {
			s.respond(w, http.StatusOK, "")
		} else {
			s.respond(w, http.StatusInternalServerError, "")
		}
	}
}

func (s *server) handleUserLinks() http.HandlerFunc {
	type userLink struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		links, err := s.store.Link().GetAllByUserID(r.Context().Value(userIDContextKey).(string))
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, http.StatusNoContent, err)
				return
			}
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		userLinks := []userLink{}
		for _, current := range links {
			userLinks = append(userLinks, userLink{
				ShortURL:    s.baseURL + "/" + current.Code,
				OriginalURL: current.Link,
			})
		}

		s.respondJSON(w, http.StatusCreated, userLinks)
	}
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(userIDCookieKey)

		var newUserID string
		if err == http.ErrNoCookie {
			newUserID = uuid.New().String()
			if err != nil {
				s.error(w, http.StatusInternalServerError, err)
				return
			}

			cookie := &http.Cookie{
				Name:  userIDCookieKey,
				Value: newUserID,
			}
			http.SetCookie(w, cookie)
		} else {
			newUserID = cookie.Value
		}

		next.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(),
				userIDContextKey,
				newUserID,
			),
		))
	})
}

func (s *server) gzipMiddleware(next http.Handler) http.Handler {
	type gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "application/gzip") {
			result, err := gzip.NewReader(r.Body)
			if err != nil {
				s.error(w, http.StatusInternalServerError, err)
				return
			}
			r.Body = result
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "application/gzip") {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				s.error(w, http.StatusInternalServerError, err)
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "application/gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
