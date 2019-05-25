package internal

import "github.com/lyraproj/goni/goni"

type Analyzer struct {
	regex *Regex
	enc goni.Encoding
}