// Package ware implements HTTP middleware functions.
package ware

import (
	"log"
	"net/http"
)

// wrapWriter is a ResponseWriter wrapper that captures outgoing request data.
type wrapWriter struct {
	http.ResponseWriter
	code int
	size int
}

// WriteHeader captures the status code before writing.
func (w *wrapWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

// Write captures the response size before writing.
func (w *wrapWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

// LogWare wraps a HandlerFunc to log each outgoing request.
func LogWare(hfun http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrap := &wrapWriter{w, 0, 0}
		hfun(wrap, r)
		log.Printf(
			"%s %s %s -> %d %d",
			r.RemoteAddr, r.Method, r.URL.Path, wrap.code, wrap.size,
		)
	}
}

// Wrap returns a HandlerFunc wrapped in middleware.
func Wrap(hfun http.HandlerFunc) http.HandlerFunc {
	return LogWare(hfun)
}
