package temp

import (
	"testing"

	"github.com/gesedels/stnts/stnts/tools/test"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	// success
	tobj, err := Parse(test.MockFS, "base.html", "main.html")
	assert.Equal(t, "base.html", tobj.Name())
	assert.Equal(t, tobj, Cache["base.html|main.html"])
	assert.NoError(t, err)
}

func TestRender(t *testing.T) {
	// setup
	tobj, _ := Parse(test.MockFS, "base.html", "main.html")

	// success
	bytes, err := Render(tobj, "test")
	assert.Equal(t, " pipeline=test \n", string(bytes))
	assert.NoError(t, err)
}
