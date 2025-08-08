// Package resp implements HTTP response writing functions.
package resp

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gesedels/stnts/stnts/tools/tpls"
)

// Error writes a formatted text/plain error message to a ResponseWriter.
func Error(w http.ResponseWriter, code int, text string, elems ...any) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	text = fmt.Sprintf("error %d: %s", code, text)
	fmt.Fprintf(w, text, elems...)
}

// HTML writes a rendered HTML template to a ResponseWriter.
func HTML(w http.ResponseWriter, tobj *template.Template, code int, pipe any) {
	bytes, err := tpls.Render(tobj, pipe)
	if err != nil {
		Error(w, http.StatusInternalServerError, "template error")
		log.Printf("error: %s", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	w.Write(bytes)
}
