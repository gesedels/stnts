// Package test implements unit testing data and functions.
package test

import (
	"io"
	"net/http/httptest"
)

// Response returns the status code and body string from a ResponseRecorder.
func Response(w *httptest.ResponseRecorder) (int, string) {
	rslt := w.Result()
	bytes, err := io.ReadAll(rslt.Body)
	if err != nil {
		panic(err)
	}

	return rslt.StatusCode, string(bytes)
}
