// Package resp implements HTTP response functions.
package resp

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gesedels/stnts/stnts/tools/temp"
)

// writeHeaders writes a Content-Type and status code to a ResponseWriter.
func writeHeaders(w http.ResponseWriter, code int, ctyp string) {
	w.Header().Set("Content-Type", ctyp)
	w.WriteHeader(code)
}

// Error writes a formatted text/plain error message to a ResponseWriter.
func Error(w http.ResponseWriter, code int, text string, elems ...any) {
	writeHeaders(w, code, "text/plain; charset=utf-8")
	text = fmt.Sprintf("error %d: %s\n", code, text)
	fmt.Fprintf(w, text, elems...)
}

// Template writes a rendered HTML Template to a ResponseWriter.
func Template(w http.ResponseWriter, tobj *template.Template, code int, pipe any) {
	bytes, err := temp.Render(tobj, pipe)
	if err != nil {
		Error(w, http.StatusInternalServerError, "template error")
		return
	}

	writeHeaders(w, code, "text/html; charset=utf-8")
	w.Write(bytes)
}
