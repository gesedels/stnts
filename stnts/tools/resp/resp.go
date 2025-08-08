// Package resp implements HTTP response writing functions.
package resp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Cache is a live cache of parsed templates.
var Cache = make(map[string]*template.Template)

// Mutex is a locking mutex for Cache.
var Mutex = new(sync.Mutex)

// Error writes a formatted text/plain error message to a ResponseWriter.
func Error(w http.ResponseWriter, code int, text string, elems ...any) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	text = fmt.Sprintf("error %d: %s", code, text)
	fmt.Fprintf(w, text, elems...)
}

// Render caches and writes a text/html Template from a filesystem to a ResponseWriter.
func Render(w http.ResponseWriter, fs embed.FS, code int, base, name string, pipe any) {
	if _, ok := Cache[name]; !ok {
		temp, err := template.ParseFS(fs, base, name)
		if err != nil {
			Error(w, http.StatusInternalServerError, "template error")
			log.Printf("error: cannot parse template %q - %s", name, err)
			return
		}

		Mutex.Lock()
		Cache[name] = temp
		Mutex.Unlock()
	}

	buff := bytes.NewBuffer(nil)
	if err := Cache[name].Execute(buff, pipe); err != nil {
		Error(w, http.StatusInternalServerError, "template error")
		log.Printf("error: cannot render template %q - %s", name, err)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buff.Bytes())
}
