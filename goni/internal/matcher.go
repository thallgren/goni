package internal

import "github.com/lyraproj/goni/goni"

type matcher struct {
	value int
	regex Regex
	enc goni.Encoding

	bytes []byte
	str int
	end int

	msaStart int
	msaOptions int
	msaRegion *Region
	msaBestLen int
	msaBestS int

	msaBegin int
	msaEnd int
}
