package main

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGunzipMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		serverAddr string
		baseAddr   string
		fileName   string
		input      string
	}{
		{
			"Simple test #1",
			":8080",
			"http://localhost:8080",
			"file",
			"https://www.google.com/",
		},
	}

	for _, tt := range tests {
		os.Setenv("BASE_URL", tt.baseAddr)
		os.Setenv("FILE_STORAGE_PATH", tt.fileName)
		os.Setenv("SERVER_ADDRESS", tt.serverAddr)

		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			zw := gzip.NewWriter(&buf)

			_, err := zw.Write([]byte(tt.input))
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if err := zw.Close(); err != nil {
				t.Errorf("Error: %v", err)
			}

			req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/", &buf)
			req.Header.Add("Content-Encoding", "gzip")
			writer := httptest.NewRecorder()
			router := SetupServer()
			router.ServeHTTP(writer, req)

			resp := writer.Result()

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("got = %v, want = %v", resp.StatusCode, http.StatusOK)
			}
		})
	}
}
