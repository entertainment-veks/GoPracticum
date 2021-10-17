package shortener

import (
	"bytes"
	"encoding/json"
	"go_practicum/internal/app/store/teststore"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testBaseURL = "http://localhost:8080"

func TestServer_HandleLinkCreateJSON(t *testing.T) {
	s := newServer(teststore.New(), testBaseURL)
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "Valid",
			payload: map[string]string{
				"url":    "www.google.com/",
				"userid": "userid",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid params",
			payload: map[string]string{
				"url": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/shorten", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
