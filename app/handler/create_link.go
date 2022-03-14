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

// HandleLinkCreate uses for creating one single link by raw data.
//
// Link which needs to be shortened passes as raw body.
// Returns status 201 and shorted link in response or an error.
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
			return
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

// HandleLinkCreateJson uses for creating one single link by json object.
//
// Link which needs to short passes as json object, ex: { "url":"https://google.com" }.
// Returns status 201 and shorted link in response json object,
// ex: { "result":"http://127.0.0.1:8080/code1234" } or an error.
func HandleLinkCreateJson(s store.Store, cfg config.Config) http.HandlerFunc {
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

// HandleLinkCreateAll uses for creating many links by json array.
//
// Links which need to short passes as json array, ex:
// [
// 		{
//	 		"correlation_id":"1",
//	 		"original_url":"https://google.com"
//	 	},
//	 	{
//	 		"correlation_id":"2",
//	 		"original_url":"https://ya.ru"
//		}
// ].
//
// Returns status 201 and shorted link in response json array, ex:
// [
// 		{
//	 		"correlation_id":"1",
//	 		"short_url":"https://127.0.0.1:8080/code1234"
//	 	},
//	 	{
//	 		"correlation_id":"2",
//	 		"short_url":"http://127.0.0.1:8080/1234code"
//		}
// ].
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
		var req []*requestElem
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

		var resp []*responseElem
		for _, current := range ls {
			resp = append(resp, &responseElem{
				CorrelationID: current.Code,
				Link:          cfg.BaseURL + "/" + current.Code,
			})
		}

		respondJSON(w, http.StatusCreated, resp)
	}
}
