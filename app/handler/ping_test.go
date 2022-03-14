package handler

import (
	"fmt"
	"go_practicum/app/store/teststore"
	"log"
	"net/http"
	"net/http/httptest"
)

func ExampleHandlePing() {
	store := teststore.New()

	resp := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		log.Fatal(err)
	}

	HandlePing(store).ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		log.Fatalf("Response code is %d, but expected %d", resp.Code, http.StatusCreated)
	}

	fmt.Println("Pong!")

	// Output:
	// Pong!
}
