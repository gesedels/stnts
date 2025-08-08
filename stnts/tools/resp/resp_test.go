package resp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gesedels/stnts/stnts/tools/temp"
	"github.com/gesedels/stnts/stnts/tools/test"
	"github.com/stretchr/testify/assert"
)

func TestWriteHeaders(t *testing.T) {
	// setup
	w := httptest.NewRecorder()

	// success
	writeHeaders(w, http.StatusOK, "test")
	code, _ := test.GetResponse(w)
	ctyp := w.Header().Get("Content-Type")
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "test", ctyp)
}

func TestError(t *testing.T) {
	// setup
	w := httptest.NewRecorder()

	// success
	Error(w, http.StatusInternalServerError, "%s", "test")
	code, body := test.GetResponse(w)
	ctyp := w.Header().Get("Content-Type")
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, "error 500: test\n", body)
	assert.Equal(t, "text/plain; charset=utf-8", ctyp)
}

func TestTemplate(t *testing.T) {
	// setup
	w := httptest.NewRecorder()
	tobj, _ := temp.Parse(test.MockFS, "base.html", "main.html")

	// success
	Template(w, tobj, http.StatusOK, "test")
	code, body := test.GetResponse(w)
	ctyp := w.Header().Get("Content-Type")
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, " pipeline=test \n", body)
	assert.Equal(t, "text/html; charset=utf-8", ctyp)
}
