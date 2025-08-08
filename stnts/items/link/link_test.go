package link

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockLink() *Link {
	return New("Name", "https://example.com", "⚪️")
}

func TestNew(t *testing.T) {
	// success
	link := mockLink()
	assert.Equal(t, "Name", link.Name)
	assert.Equal(t, "https://example.com", link.Addr)
	assert.Equal(t, "⚪️", link.Icon)
}

func TestString(t *testing.T) {
	// setup
	link := mockLink()

	// success
	text := link.String()
	assert.Equal(t, "https://example.com", text)
}

func TestURL(t *testing.T) {
	// setup
	link := mockLink()

	// success
	url := link.URL()
	assert.Equal(t, "https://example.com", url.String())
}
