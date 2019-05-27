package internal

import "github.com/lyraproj/goni/goni"

type Regex struct {
	enc         goni.Encoding // fast access to encoding
}

func (rx *Regex) nameToGroupNumbers([]byte, int, int) *goni.NameEntry {
	// TODO
	return nil
}


