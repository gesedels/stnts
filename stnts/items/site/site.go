// Package site implements the Site type and methods.
package site

import (
	"encoding/json"
	"fmt"
	"os"
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

// Parse returns a new Site from a parsed JSON file.
func Parse(orig string) (*Site, error) {
	bytes, err := os.ReadFile(orig)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %q - %w", orig, err)
	}

	site := new(Site)
	if err := json.Unmarshal(bytes, site); err != nil {
		return nil, fmt.Errorf("cannot parse file %q - %w", orig, err)
	}

	return site, nil
}

// Now returns the current time in the Site's configured timezone.
func (s *Site) Now() time.Time {
	return time.Now().In(s.Conf.Location())
}
