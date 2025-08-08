// Package temp implements template rendering functions.
package temp

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"sync"
)

// Cache is a live cache of parsed templates.
var Cache = make(map[string]*template.Template)

// Mutex is a locking mutex for Cache.
var Mutex = new(sync.Mutex)

// Parse returns a new or cached Template from a filesystem.
func Parse(fs fs.FS, base, name string) (*template.Template, error) {
	if _, ok := Cache[name]; !ok {
		temp, err := template.ParseFS(fs, base, name)
		if err != nil {
			return nil, fmt.Errorf("cannot parse template %q - %w", name, err)
		}

		Mutex.Lock()
		Cache[name] = temp
		Mutex.Unlock()
	}

	return Cache[name], nil
}

// Render returns a rendered Template as a byteslice.
func Render(temp *template.Template, pipe any) ([]byte, error) {
	var buff = new(bytes.Buffer)
	if err := temp.Execute(buff, pipe); err != nil {
		return nil, fmt.Errorf("cannot render template %q - %w", temp.Name(), err)
	}

	return buff.Bytes(), nil
}

// RenderFrom returns a rendered Template from a filesystem as a byteslice.
func RenderFrom(fs fs.FS, base, name string, pipe any) ([]byte, error) {
	temp, err := Parse(fs, base, name)
	if err != nil {
		return nil, err
	}

	return Render(temp, pipe)
}
