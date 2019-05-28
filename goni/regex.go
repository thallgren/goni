package goni

type Regex interface {
	Matcher(bytes []byte) Matcher
}
