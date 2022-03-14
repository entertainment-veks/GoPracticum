package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"go_practicum/app/store/teststore"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func ExampleHandleLinkDelete() {
	newLink := "https://sun.12.march.20.18.com"
	cfg := newConfig()
	store := teststore.New()

	resp := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(newLink))
	if err != nil {
		log.Fatal(err)
	}

	AuthMiddleware(HandleLinkCreate(store, cfg)).ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		log.Fatalf("Response c is %d, but expected %d", resp.Code, http.StatusCreated)
	}

	data := resp.Body.String()
	ds := strings.Split(data, "/")
	c := fmt.Sprintf("[%s]", ds[len(ds)-1])

	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(c))

	HandleLinkDelete(store).ServeHTTP(resp, req)
	if resp.Code != http.StatusAccepted {
		log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusAccepted)
	}

	time.Sleep(2 * time.Second) //deletion is async operation. wait till operation will invoke

	resp = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/{key}", nil)
	req = mux.SetURLVars(req, map[string]string{
		"key": ds[len(ds)-1],
	})

	HandleLinkGet(store).ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusBadRequest)
	}

	fmt.Println("So link was deleted")

	// Output:
	// So link was deleted
}
