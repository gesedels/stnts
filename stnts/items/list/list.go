// Package list implements the List type and methods.
package list

import "github.com/gesedels/stnts/stnts/items/link"

// List is a single ordered list of Links.
type List struct {
	Name  string       `json:"name"`
	Links []*link.Link `json:"links"`
}

// New returns a new List.
func New(name string, links ...*link.Link) *List {
	return &List{name, links}
}

// String returns the List's name as a string.
func (l *List) String() string {
	return l.Name
}
