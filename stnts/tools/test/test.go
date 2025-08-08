// Package test implements unit testing tools and data.
package test

import (
	"bytes"
	"embed"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// MockFS is an embedded filesystem containing mock files.
//
//go:embed *.*
var MockFS embed.FS

// GetResponse returns the status code and body from a ResponseRecorder.
func GetResponse(w *httptest.ResponseRecorder) (int, string) {
	rslt := w.Result()
	bytes, err := io.ReadAll(rslt.Body)
	if err != nil {
		panic(err)
	}

	return rslt.StatusCode, string(bytes)
}

// NewRequest returns a new mock Request object.
func NewRequest(path, body string) *http.Request {
	buff := bytes.NewBufferString(body)
	meth, path, _ := strings.Cut(path, " ")
	return httptest.NewRequest(meth, path, buff)
}
