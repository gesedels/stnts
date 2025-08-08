package resp

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gesedels/stnts/stnts/tools/test"
	"github.com/stretchr/testify/assert"
)

//go:embed *.html
var mockFS embed.FS

func TestError(t *testing.T) {
	// setup
	w := httptest.NewRecorder()

	// success
	Error(w, http.StatusNotFound, "%s", "test")
	code, body := test.Response(w)
	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, "error 404: test", body)
}

func TestRender(t *testing.T) {
	// setup
	clear(Cache)
	w := httptest.NewRecorder()

	// success
	Render(w, mockFS, http.StatusOK, "base.html", "main.html", "test")
	code, body := test.Response(w)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, " pipeline=test \n", body)
	assert.NotNil(t, Cache["main.html"])
}
