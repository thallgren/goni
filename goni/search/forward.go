package search

import "github.com/lyraproj/goni/goni"

type Forward interface {
	Name() string
	Search(matcher goni.Matcher, text []byte, textP, textEnd, textRange int)
}
