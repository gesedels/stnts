// Package temp implements template rendering functions.
package temp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"sync"
)

// Cache is a map of cached parsed Templates.
var Cache = make(map[string]*template.Template)

// CacheLock is a write mutex for Cache.
var CacheLock = new(sync.Mutex)

// Parse returns a new or cached Template from an embedded filesystem.
func Parse(fs embed.FS, names ...string) (*template.Template, error) {
	name := strings.Join(names, "|")
	if _, ok := Cache[name]; !ok {
		tobj, err := template.ParseFS(fs, names...)
		if err != nil {
			return nil, fmt.Errorf("cannot parse template - %w", err)
		}

		CacheLock.Lock()
		Cache[name] = tobj
		CacheLock.Unlock()
	}

	return Cache[name], nil
}

// Render returns a rendered Template as a byteslice.
func Render(tobj *template.Template, pipe any) ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := tobj.Execute(buff, pipe); err != nil {
		return nil, fmt.Errorf("cannot render template - %w", err)
	}

	return buff.Bytes(), nil
}
