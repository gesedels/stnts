package list

import (
	"testing"

	"github.com/gesedels/stnts/stnts/items/link"
	"github.com/stretchr/testify/assert"
)

func mockList() *List {
	return New("Name", link.New("Name", "https://example.com", "⚪️"))
}

func TestNew(t *testing.T) {
	// success
	list := mockList()
	assert.Equal(t, "Name", list.Name)
	assert.Len(t, list.Links, 1)
	assert.Equal(t, "https://example.com", list.Links[0].Addr)
}

func TestString(t *testing.T) {
	// setup
	list := mockList()

	// success
	text := list.String()
	assert.Equal(t, "Name", text)
}
