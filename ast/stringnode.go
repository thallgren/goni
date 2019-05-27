package ast

import (
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
	"io"
)

type StringNode struct {
	abstractNode
	bytes []byte
	flag  int
}

const (
	NstrRaw            = 1 << 0
	NstrAmbig          = 1 << 1
	NstrDontGetOptInfo = 1 << 2
	NstrShared         = 1 << 3

	nodeStrMargin  = 16
	nodeStrBufSize = 24
)

var StringNodeEmpty = &StringNode{abstractNode: abstractNode{nodeType: node.Str}, bytes: []byte{}}

func NewStringNodeWithCapacity(size int) *StringNode {
	return &StringNode{abstractNode: abstractNode{nodeType: node.Str}, bytes: make([]byte, 0, size)}
}

func NewStringNode() *StringNode {
	return &StringNode{abstractNode: abstractNode{nodeType: node.Str}, bytes: make([]byte, 0, nodeStrBufSize)}
}

func (sn *StringNode) String() string {
	return goni.String(sn)
}

func (sn *StringNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	if sn.flag != 0 {
		w.Append(`flags: `)
		sn.appendFlags(w)
		w.NewLine()
	}
	w.Append(`bytes: '`)
	for _, b := range sn.bytes {
		u := uint(b)
		if u >= 0x20 && u < 0x7f {
			w.Append(string(b));
		} else {
			w.Printf("[0x%02x]", u);
		}
	}
	w.Append("'");
}

func (sn *StringNode) Name() string {
	return `String`
}

func (sn *StringNode) ClearAmbig() {
	sn.flag &= ^NstrAmbig
}

func (sn *StringNode) ClearDontGetOptInfo() {
	sn.flag &= ^NstrDontGetOptInfo
}

func (sn *StringNode) ClearShared() {
	sn.flag &= ^NstrShared
}

func (sn *StringNode) ClearRaw() {
	sn.flag &=^NstrRaw
}

func (sn *StringNode) IsAmbig() bool {
	return (sn.flag & NstrAmbig) != 0
}

func (sn *StringNode) IsDontGetOptInfo() bool {
	return (sn.flag & NstrDontGetOptInfo) != 0
}

func (sn *StringNode) IsShared() bool {
	return (sn.flag & NstrShared) != 0
}

func (sn *StringNode) IsRaw() bool {
	return (sn.flag & NstrRaw) != 0
}

func (sn *StringNode) SetAmbig() {
	sn.flag |=  NstrAmbig
}

func (sn *StringNode) SetDontGetOptInfo() {
	sn.flag |= NstrDontGetOptInfo
}

func (sn *StringNode) SetShared() {
	sn.flag |= NstrShared
}

func (sn *StringNode) SetRaw() {
	sn.flag |= NstrRaw
}

func (sn *StringNode) appendFlags(w io.Writer) {
	if sn.IsRaw() {
		util.WriteString(w, `RAW `)
	}
	if sn.IsAmbig() {
		util.WriteString(w, `AMBIG `)
	}
	if sn.IsDontGetOptInfo() {
		util.WriteString(w, `DONT_GET_OPT_INFO `)
	}
	if sn.IsShared() {
		util.WriteString(w, `SHARED `)
	}
}
