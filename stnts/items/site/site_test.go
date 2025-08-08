package site

import (
	"testing"

	"github.com/gesedels/stnts/stnts/items/conf"
	"github.com/gesedels/stnts/stnts/items/link"
	"github.com/gesedels/stnts/stnts/items/list"
	"github.com/stretchr/testify/assert"
)

func mockSite() *Site {
	conf := conf.New("Name", "Desc", "Local")
	icons := []*link.Link{link.New("Name", "https://example.com", "⚪️")}
	lists := []*list.List{list.New("Name", link.New("Name", "https://example.com", "⚪️"))}
	return New(conf, icons, lists)
}

func TestNew(t *testing.T) {
	// success
	site := mockSite()
	assert.NotNil(t, site.Conf)
	assert.NotNil(t, site.Icons)
	assert.NotNil(t, site.Lists)
}

func TestParse(t *testing.T) {
	// todo
}

func TestNow(t *testing.T) {
	// setup
	site := mockSite()

	// success
	tnow := site.Now()
	assert.Equal(t, site.Conf.Location(), tnow.Location())
}
