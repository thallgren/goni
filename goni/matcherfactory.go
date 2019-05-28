package goni

type MatcherFactory interface {
	Create(regex Regex, region Region, bytes []byte)
}
