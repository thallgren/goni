package goni

import "github.com/lyraproj/goni/util"

const initNameBackrefsAllocNum = 8

type NameEntry struct {
	name string
	backRefs []int
}

func NewNameEntry(bytes []byte, p, end int) *NameEntry {
	return &NameEntry{string(bytes[p:end]), nil}
}

func (ne *NameEntry) AddBackRef(backRef int) {
	ne.backRefs = append(ne.backRefs, backRef)
}

func (ne *NameEntry) BackRefs() []int {
	return ne.backRefs
}

func (ne *NameEntry) Len() int {
	return len(ne.backRefs)
}

func (ne *NameEntry) String() string {
	return util.String(ne)
}

func (ne *NameEntry) AppendTo(w *util.Indenter) {
	w.Append(ne.name)
	if len(ne.backRefs) == 0 {
		w.Append(` -`)
	} else {
		w.Append(` `)
		for i, br := range ne.backRefs {
			if i > 0 {
				w.Append(`, `)
			}
			w.AppendInt(br)
		}
	}
}