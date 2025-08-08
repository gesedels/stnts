package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResponse(t *testing.T) {
	// setup
	w := httptest.NewRecorder()
	fmt.Fprint(w, "test")

	// success
	code, body := GetResponse(w)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "test", body)
}

func TestNewRequest(t *testing.T) {
	// success
	r := NewRequest("GET /", "test")
	bytes, _ := io.ReadAll(r.Body)
	assert.Equal(t, "GET", r.Method)
	assert.Equal(t, "/", r.URL.Path)
	assert.Equal(t, "test", string(bytes))
}
