// Package link implements the Link type and methods.
package link

import "net/url"

// Link is a single web address with a name and icon.
type Link struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	Icon string `json:"icon"`
}

// New returns a new Link.
func New(name, addr, icon string) *Link {
	return &Link{name, addr, icon}
}

// String returns the Link's address as a string.
func (l *Link) String() string {
	return l.Addr
}

// URL returns the Link's address as a parsed URL.
func (l *Link) URL() *url.URL {
	url, _ := url.Parse(l.Addr)
	return url
}
