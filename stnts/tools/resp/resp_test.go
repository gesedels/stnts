package resp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gesedels/stnts/stnts/tools/test"
	"github.com/gesedels/stnts/stnts/tools/tpls"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	// setup
	w := httptest.NewRecorder()

	// success
	Error(w, http.StatusNotFound, "%s", "test")
	code, body := test.Response(w)
	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, "error 404: test", body)
}

func TestHTML(t *testing.T) {
	// setup
	clear(tpls.Cache)
	w := httptest.NewRecorder()
	tobj, _ := tpls.Parse(test.MockFS, "base.html", "main.html")

	// success
	HTML(w, tobj, http.StatusOK, "test")
	code, body := test.Response(w)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, " pipeline=test \n", body)
}
