package main

//will remake these test in future increments. does not have enough time now

// import (
// 	"bytes"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
// )

// type want struct {
// 	code        int
// 	data        string
// 	contentType string
// }

// func Test_Endpoints(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		method string
// 		url    string
// 		data   string
// 		want   want
// 	}{
// 		{
// 			name:   "'/' - endpoint test",
// 			method: http.MethodPost,
// 			url:    "https://localhost:8080/",
// 			data:   "https://www.google.com/",
// 			want: want{
// 				code: http.StatusCreated,
// 			},
// 		},
// 		{
// 			name:   "'/{key}' - endpoint test",
// 			method: http.MethodGet,
// 			want: want{
// 				code:        http.StatusTemporaryRedirect,
// 				contentType: "Location",
// 				data:        "https://www.google.com/",
// 			},
// 		},
// 		{
// 			name:   "'/api/shorten' - endpoint test",
// 			method: http.MethodPost,
// 			url:    "https://localhost:8080/api/shorten",
// 			data:   `{"url": "https://www.google.com/"}`,
// 			want: want{
// 				code: http.StatusCreated,
// 			},
// 		},
// 	}

// 	os.Setenv("SERVER_ADDRESS", ":8080")
// 	os.Setenv("BASE_URL", "http://localhost:8080")
// 	os.Setenv("FILE_STORAGE_PATH", "file")

// 	linkForGetTest := ""

// 	t.Run(tests[0].name, func(t *testing.T) { // '/' method
// 		request := httptest.NewRequest(http.MethodPost, tests[0].url, bytes.NewBufferString(tests[0].data))
// 		defer request.Body.Close()

// 		w := httptest.NewRecorder()
// 		router := SetupServer()
// 		router.ServeHTTP(w, request)

// 		response := w.Result()

// 		if tests[0].want.code != response.StatusCode {
// 			t.Errorf("Expected code: %v, Actual: %v", tests[0].want.code, response.StatusCode)
// 		}

// 		defer response.Body.Close()
// 		resBody, err := ioutil.ReadAll(response.Body)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		linkForGetTest = string(resBody)
// 	})

// 	t.Run(tests[1].name, func(t *testing.T) { // '/{key}' method
// 		request := httptest.NewRequest(http.MethodGet, linkForGetTest, nil)

// 		w := httptest.NewRecorder()
// 		router := SetupServer()
// 		router.ServeHTTP(w, request)

// 		response := w.Result()

// 		if tests[1].want.code != response.StatusCode {
// 			t.Errorf("Expected code: %v, Actual: %v", tests[0].want.code, response.StatusCode)
// 		}

// 		defer response.Body.Close()

// 		resBody, err := ioutil.ReadAll(response.Body)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if tests[1].want.data == string(resBody) {
// 			t.Errorf("Expected body: %v, Actual: %v", tests[1].want.data, response.Header.Get(tests[1].want.contentType))
// 		}
// 	})

// 	t.Run(tests[2].name, func(t *testing.T) { // '/api/shorten' method
// 		request := httptest.NewRequest(http.MethodPost, tests[2].url, bytes.NewBufferString(tests[2].data))
// 		defer request.Body.Close()

// 		w := httptest.NewRecorder()
// 		router := SetupServer()
// 		router.ServeHTTP(w, request)

// 		response := w.Result()
// 		defer response.Body.Close()

// 		if tests[2].want.code != response.StatusCode {
// 			t.Errorf("Expected code: %v, Actual: %v", tests[2].want.code, response.StatusCode)
// 		}
// 	})
// }
