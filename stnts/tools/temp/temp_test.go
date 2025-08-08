package temp

import (
	"testing"

	"github.com/gesedels/stnts/stnts/tools/test"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	// setup
	clear(Cache)

	// success
	temp, err := Parse(test.MockFS, "base.html", "main.html")
	assert.Equal(t, "base.html", temp.Name())
	assert.Equal(t, temp, Cache["base.html|main.html"])
	assert.NoError(t, err)
}

func TestRender(t *testing.T) {
	// setup
	temp, _ := Parse(test.MockFS, "base.html", "main.html")

	// success
	bytes, err := Render(temp, "test")
	assert.Equal(t, " pipeline=test \n", string(bytes))
	assert.NoError(t, err)
}
