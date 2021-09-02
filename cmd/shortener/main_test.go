package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type want struct {
	code        int
	data        string
	contentType string
}

func Test_Endpoints(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		data   string
		want   want
	}{
		{
			name:   "POST-endpoint test",
			method: http.MethodPost,
			url:    "https://localhost:8080/",
			data:   "https://www.google.com/",
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name:   "GET-endpoint test",
			method: http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "Location",
				data:        "https://www.google.com/",
			},
		},
	}

	linkForGetTest := ""

	t.Run(tests[0].name, func(t *testing.T) { //POST method
		request := httptest.NewRequest(http.MethodPost, tests[0].url, bytes.NewBufferString(tests[0].data))
		defer request.Body.Close()

		w := httptest.NewRecorder()
		router := SetupServer()
		router.ServeHTTP(w, request)

		response := w.Result()

		if tests[0].want.code != response.StatusCode {
			t.Errorf("Expected code: %v, Actual: %v", tests[0].want.code, response.StatusCode)
		}

		defer response.Body.Close()
		resBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		linkForGetTest = string(resBody)
	})

	t.Run(tests[1].name, func(t *testing.T) { //GET method
		request := httptest.NewRequest(http.MethodGet, linkForGetTest, nil)

		w := httptest.NewRecorder()
		router := SetupServer()
		router.ServeHTTP(w, request)

		response := w.Result()

		if tests[1].want.code != response.StatusCode {
			t.Errorf("Expected code: %v, Actual: %v", tests[0].want.code, response.StatusCode)
		}

		defer response.Body.Close()

		if tests[1].want.data == response.Header.Get(tests[1].want.contentType) {
			t.Errorf("!Expected body: %v, Actual: %v", tests[1].want.data, response.Header.Get(tests[1].want.contentType))
		}
	})
}
