// Package conf implements the Conf type and methods.
package conf

import (
	"time"
)

// Conf is a single static configuration map.
type Conf struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	TimeZone string `json:"timezone"`
}

// New returns a new Conf.
func New(name, desc, zone string) *Conf {
	return &Conf{name, desc, zone}
}

// Location returns the Conf's parsed timezone.
func (c *Conf) Location() *time.Location {
	loca, _ := time.LoadLocation(c.TimeZone)
	return loca
}
