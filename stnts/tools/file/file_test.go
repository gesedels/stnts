package file

import (
	"testing"

	"github.com/gesedels/stnts/stnts/tools/test"
	"github.com/stretchr/testify/assert"
)

func TestReadJSON(t *testing.T) {
	// setup
	var data int
	orig := test.TempFile(t, "test.json", "1234")

	// success
	err := ReadJSON(orig, &data)
	assert.Equal(t, 1234, data)
	assert.NoError(t, err)
}
