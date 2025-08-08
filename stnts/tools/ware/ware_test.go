package ware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gesedels/stnts/stnts/tools/test"
	"github.com/stretchr/testify/assert"
)

func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "test")
}

func TestWriteHeader(t *testing.T) {
	// setup
	w := httptest.NewRecorder()
	wrap := &wrapWriter{w, 0, 0}

	// success
	wrap.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, wrap.code)
}

func TestWrite(t *testing.T) {
	// setup
	w := httptest.NewRecorder()
	wrap := &wrapWriter{w, 0, 0}

	// success
	wrap.Write([]byte("test"))
	assert.Equal(t, 4, wrap.size)
}

func TestLogWare(t *testing.T) {
	// setup
	r := test.NewRequest("GET /", "")
	w := httptest.NewRecorder()
	buff := new(bytes.Buffer)
	log.SetOutput(buff)
	log.SetFlags(0)

	// success
	LogWare(mockHandler)(w, r)
	code, body := test.GetResponse(w)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "test", body)
	assert.Equal(t, "192.0.2.1:1234 GET / -> 200 4\n", buff.String())
}

func TestWrap(t *testing.T) {
	// success
	hfun := Wrap(mockHandler)
	assert.NotNil(t, hfun)
}
