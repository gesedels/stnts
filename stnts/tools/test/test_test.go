package test

import (
	"fmt"
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
