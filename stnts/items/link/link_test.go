package link

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockLink() *Link {
	return New("https://example.com", "⚪️", "Example")
}

func TestNew(t *testing.T) {
	// success
	link := mockLink()
	assert.Equal(t, "https://example.com", link.Addr)
	assert.Equal(t, "⚪️", link.Icon)
	assert.Equal(t, "Example", link.Name)
}

func TestURL(t *testing.T) {
	// setup
	link := mockLink()

	// success
	url := link.URL()
	assert.Equal(t, "https://example.com", url.String())
}
