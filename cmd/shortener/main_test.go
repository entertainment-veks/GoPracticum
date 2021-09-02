package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
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

	StartServer()
	linkForGetTest := ""

	t.Run(tests[0].name, func(t *testing.T) { //POST method

		client := &http.Client{}

		request, _ := http.NewRequest(http.MethodPost, tests[0].url, bytes.NewBufferString(tests[0].data))
		response, err := client.Do(request)
		if err != nil {
			t.Fail()
		}

		body, _ := ioutil.ReadAll(response.Body)
		linkForGetTest = string(body)

		if tests[0].want.code != response.StatusCode {
			t.Errorf("Expected: %v, Actual: %v", tests[0].want.code, response.StatusCode)
		}
	})

	t.Run(tests[1].name, func(t *testing.T) { //GET method
		client := &http.Client{}

		request, _ := http.NewRequest(http.MethodGet, linkForGetTest, nil)
		response, err := client.Do(request)
		if err != nil {
			t.Fail()
		}

		if tests[1].want.code != response.StatusCode {
			t.Errorf("Expected: %v, Actual: %v", tests[0].want.code, response.StatusCode)
		}

		if tests[1].want.data == response.Header.Get(tests[1].want.contentType) {
			t.Errorf("Expected: %v, Actual: %v", tests[1].want.data, response.Header.Get(tests[1].want.contentType))
		}
	})
}
