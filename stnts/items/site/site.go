// Package site implements the Site type and methods.
package site

import (
	"time"

	"github.com/gesedels/stnts/stnts/items/conf"
	"github.com/gesedels/stnts/stnts/items/link"
	"github.com/gesedels/stnts/stnts/items/list"
)

// Site is a global site data container.
type Site struct {
	Conf  *conf.Conf   `json:"conf"`
	Icons []*link.Link `json:"icons"`
	Lists []*list.List `json:"lists"`
}

// New returns a new Site.
func New(conf *conf.Conf, icons []*link.Link, lists []*list.List) *Site {
	return &Site{conf, icons, lists}
}

// Now returns the current time in the Site's configured timezone.
func (s *Site) Now() time.Time {
	return time.Now().In(s.Conf.Location())
}
