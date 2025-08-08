// Package templ implements template rendering functions.
package temp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"sync"
)

// Cache is a live cache of parsed templates.
var Cache = make(map[string]*template.Template)

// Mutex is a write-locking mutex for Cache.
var Mutex = new(sync.Mutex)

// Parse returns a new or cached Template from a filesystem.
func Parse(fs embed.FS, names ...string) (*template.Template, error) {
	name := strings.Join(names, "|")

	if _, ok := Cache[name]; !ok {
		temp, err := template.ParseFS(fs, names...)
		if err != nil {
			return nil, fmt.Errorf("cannot parse template - %w", err)
		}

		Mutex.Lock()
		Cache[name] = temp
		Mutex.Unlock()
	}

	return Cache[name], nil
}

// Render returns a rendered Template as a byteslice.
func Render(temp *template.Template, pipe any) ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := temp.Execute(buff, pipe); err != nil {
		return nil, fmt.Errorf("cannot render template - %w", err)
	}

	return buff.Bytes(), nil
}
