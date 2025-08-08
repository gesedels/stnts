package temp

import (
	"embed"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed *.html
var mockFS embed.FS

func TestParse(t *testing.T) {
	// success - no cache
	temp, err := Parse(mockFS, "base.html", "name.html")
	assert.Equal(t, "base.html", temp.Name())
	assert.Equal(t, temp, Cache["name.html"])
	assert.NoError(t, err)

	// success - with cache
	temp, err = Parse(mockFS, "base.html", "name.html")
	assert.Equal(t, temp, Cache["name.html"])
	assert.NoError(t, err)
}

func TestRender(t *testing.T) {
	// setup
	temp, _ := template.ParseFS(mockFS, "base.html", "name.html")

	// success
	bytes, err := Render(temp, "test")
	assert.Equal(t, " pipeline=test \n", string(bytes))
	assert.NoError(t, err)
}

func TestRenderFrom(t *testing.T) {
	// success
	bytes, err := RenderFrom(mockFS, "base.html", "name.html", "test")
	assert.Equal(t, " pipeline=test \n", string(bytes))
	assert.NoError(t, err)
}
