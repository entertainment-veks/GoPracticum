package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go_practicum/app/store/teststore"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleHandleLinkCreate() {
	newLink := "https://fri.11.march.06.35.com"
	cfg := newConfig()
	store := teststore.New()

	resp := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(newLink))
	if err != nil {
		log.Fatal(err)
	}

	AuthMiddleware(HandleLinkCreate(store, cfg)).ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusCreated)
	}

	data := resp.Body.String()
	ds := strings.Split(data, "/")
	code := ds[len(ds)-1]

	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/{key}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"key": code,
	})

	HandleLinkGet(store).ServeHTTP(resp, req)
	if resp.Code != http.StatusTemporaryRedirect {
		log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusTemporaryRedirect)
	}

	fmt.Printf("Short link redirect us to: %s", resp.Header().Get("Location"))

	// Output:
	// Short link redirect us to: https://fri.11.march.06.35.com
}

func ExampleHandleLinkCreateJson() {
	type request struct {
		Link string `json:"url"`
	}

	type response struct {
		Output string `json:"result"`
	}

	newLink := request{
		Link: "https://fri.11.march.10.12.com",
	}
	cfg := newConfig()
	store := teststore.New()

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(newLink)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/api/shorten", &b)
	if err != nil {
		log.Fatal(err)
	}

	resp := httptest.NewRecorder()
	AuthMiddleware(HandleLinkCreateJson(store, cfg)).ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusCreated)
	}

	var d response
	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		log.Fatal(err)
	}
	ds := strings.Split(d.Output, "/")
	code := ds[len(ds)-1]

	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/{key}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"key": code,
	})

	HandleLinkGet(store).ServeHTTP(resp, req)
	if resp.Code != http.StatusTemporaryRedirect {
		log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusTemporaryRedirect)
	}

	fmt.Printf("Short link redirect us to: %s", resp.Header().Get("Location"))

	// Output:
	// Short link redirect us to: https://fri.11.march.10.12.com
}

func ExampleHandleLinkCreateAll() {
	type requestElem struct {
		CorrelationID string `json:"correlation_id"`
		Link          string `json:"original_url"`
	}

	type responseElem struct {
		CorrelationID string `json:"correlation_id"`
		Link          string `json:"short_url"`
	}

	cfg := newConfig()
	store := teststore.New()
	newLinks := []requestElem{
		{
			CorrelationID: "1",
			Link:          "https://sat.12.march.17.16.com",
		},
		{
			CorrelationID: "2",
			Link:          "https://sat.12.march.17.20.com",
		},
		{
			CorrelationID: "3",
			Link:          "https://sat.12.march.18.24.com",
		},
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(newLinks)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/api/shorten/batch", &b)
	if err != nil {
		log.Fatal(err)
	}

	resp := httptest.NewRecorder()
	AuthMiddleware(HandleLinkCreateAll(store, cfg)).ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusCreated)
	}

	var data []responseElem
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range data {
		ds := strings.Split(e.Link, "/")
		code := ds[len(ds)-1]

		resp = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/{key}", nil)
		req = mux.SetURLVars(req, map[string]string{
			"key": code,
		})

		HandleLinkGet(store).ServeHTTP(resp, req)
		if resp.Code != http.StatusTemporaryRedirect {
			log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusTemporaryRedirect)
		}

		fmt.Printf("Short link %s redirect us to: %s\n", e.CorrelationID, resp.Header().Get("Location"))
	}

	// Output:
	// Short link 1 redirect us to: https://sat.12.march.17.16.com
	// Short link 2 redirect us to: https://sat.12.march.17.20.com
	// Short link 3 redirect us to: https://sat.12.march.18.24.com
}
