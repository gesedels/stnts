// Package test implements unit testing tools and data.
package test

import (
	"embed"
	"io"
	"net/http/httptest"
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
