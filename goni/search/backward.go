package search

import "github.com/lyraproj/goni/goni"

type Backward interface {
	Search(matcher goni.Matcher, text []byte, textP, adjustText, textEnd, textStart, s_, range_ int)
}
