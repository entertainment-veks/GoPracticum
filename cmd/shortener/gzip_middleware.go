package main

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// type gzipWriter struct {
// 	http.ResponseWriter
// 	w io.Writer
// }

type gzipReader struct {
	*http.Request
	rcer io.ReadCloser
}

func (gz *gzipReader) Body() io.ReadCloser {
	return gz.rcer
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mr *http.Request

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			result, _ := gzip.NewReader(r.Body)

			tempo, _ := ioutil.ReadAll(result)
			tempo = tempo

			mr = gzipReader{
				r,
				result,
			}
		} else {
			mr = r
		}

		next.ServeHTTP(w, mr)
	})
}

// if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
// 	var b bytes.Buffer
// 	tempoWriter := gzip.NewWriter(&b)

// 	data, _ := ioutil.ReadAll(r.Body)
// 	tempoWriter.Write(data)

// 	tempo := b.String()
// 	tempo = tempo

// 	rw = gzipWriter{
// 		w,
// 		tempoWriter,
// 	}
// }
