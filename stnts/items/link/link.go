// Package link implements the Link type and methods.
package link

import "net/url"

// Link is a single web address with a name and icon.
type Link struct {
	Addr string `json:"addr"`
	Icon string `json:"icon"`
	Name string `json:"name"`
}

// New returns a new Link.
func New(addr, icon, name string) *Link {
	return &Link{addr, icon, name}
}

// URL returns the Link's address as a parsed URL.
func (l *Link) URL() *url.URL {
	url, _ := url.Parse(l.Addr)
	return url
}
