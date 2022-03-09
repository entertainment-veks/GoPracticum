package handler

import (
	"encoding/json"
	"errors"
	"go_practicum/app/config"
	"go_practicum/app/model"
	"go_practicum/app/store"
	"go_practicum/app/util"
	"io/ioutil"
	"net/http"
)

func HandleLinkCreate(s store.Store, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		code, err := util.GenerateCode()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		userID, ok := r.Context().Value(userIDContextKey).(string)
		if !ok {
			respondError(w, http.StatusInternalServerError, errors.New("user id is not provided"))
		}

		l := &model.Link{
			Link:   string(body),
			Code:   code,
			UserID: userID,
		}
		if err := l.Validate(); err != nil {
			respondError(w, http.StatusUnprocessableEntity, err)
			return
		}
		if err := s.Link().Create(l); err != nil {
			if errors.Is(err, store.ErrConflict) {
				respondError(w, http.StatusConflict, err)
				return
			}
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		respond(w, http.StatusCreated, cfg.BaseURL+"/"+code)
	}
}

func HandleLinkCreateJSON(s store.Store, cfg config.Config) http.HandlerFunc {
	type request struct {
		Link string `json:"url"`
	}

	type response struct {
		Output string `json:"result"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		code, err := util.GenerateCode()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		userID, ok := r.Context().Value(userIDContextKey).(string)
		if !ok {
			respondError(w, http.StatusInternalServerError, errors.New("user id is not provided"))
		}

		l := &model.Link{
			Link:   req.Link,
			Code:   code,
			UserID: userID,
		}
		if err := l.Validate(); err != nil {
			respondError(w, http.StatusUnprocessableEntity, err)
			return
		}
		if err := s.Link().Create(l); err != nil {
			if errors.Is(err, store.ErrConflict) {
				respondError(w, http.StatusConflict, err)
				return
			}
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		rawResult := &response{
			cfg.BaseURL + "/" + code,
		}
		respondJSON(w, http.StatusCreated, rawResult)
	}
}

func HandleLinkCreateAll(s store.Store, cfg config.Config) http.HandlerFunc {
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
			respondError(w, http.StatusBadRequest, err)
			return
		}

		userID, ok := r.Context().Value(userIDContextKey).(string)
		if !ok {
			respondError(w, http.StatusInternalServerError, errors.New("user id is not provided"))
		}

		var ls []*model.Link
		for _, current := range req {
			l := &model.Link{
				Link:   current.Link,
				Code:   current.CorrelationID,
				UserID: userID,
			}
			if err := l.Validate(); err != nil {
				respondError(w, http.StatusUnprocessableEntity, err)
				return
			}
			ls = append(ls, l)
		}

		if err := s.Link().CreateAll(ls); err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		resp := []*responseElem{}
		for _, current := range ls {
			resp = append(resp, &responseElem{
				CorrelationID: current.Code,
				Link:          cfg.BaseURL + "/" + current.Code,
			})
		}

		respondJSON(w, http.StatusCreated, resp)
	}
}
