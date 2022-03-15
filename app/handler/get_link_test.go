package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go_practicum/app/store/teststore"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
)

func ExampleHandleLinkGet() {
	newLink := "https://sat.13.march.09.11.com"
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
	// Short link redirect us to: https://sat.13.march.09.11.com
}

func ExampleHandleGetUserLinks() {
	type userLink struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	newLinks := []string{
		"https://sat.13.march.09.11.com",
		"https://sat.13.march.09.12.com",
		"https://sat.13.march.09.13.com",
		"https://sat.13.march.09.14.com",
		"https://sat.13.march.09.15.com",
	}
	cfg := newConfig()
	store := teststore.New()

	resp := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, newLink := range newLinks {
		req.Body = ioutil.NopCloser(strings.NewReader(newLink))

		AuthMiddleware(HandleLinkCreate(store, cfg)).ServeHTTP(resp, req)
		if resp.Code != http.StatusCreated {
			log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusCreated)
		}

		authCookie := strings.TrimPrefix(resp.Result().Header["Set-Cookie"][0], "shortener-userid=")
		cookie := &http.Cookie{
			Name:  "shortener-userid",
			Value: authCookie,
		}
		req.AddCookie(cookie)
	}

	authCookie := strings.TrimPrefix(resp.Result().Header["Set-Cookie"][0], "shortener-userid=")
	cookie := &http.Cookie{
		Name:  "shortener-userid",
		Value: authCookie,
	}

	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/user/urls", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.AddCookie(cookie)

	AuthMiddleware(HandleGetUserLinks(store, cfg)).ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusOK)
	}

	var result []userLink
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].OriginalURL < result[j].OriginalURL
	})

	fmt.Println("Link which was created by this user:")
	for _, e := range result {
		fmt.Println(e.OriginalURL)
	}

	// Output:
	// Link which was created by this user:
	// https://sat.13.march.09.11.com
	// https://sat.13.march.09.12.com
	// https://sat.13.march.09.13.com
	// https://sat.13.march.09.14.com
	// https://sat.13.march.09.15.com
}
