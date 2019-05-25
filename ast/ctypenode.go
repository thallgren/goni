package ast

import (
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/character"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
)

type CTypeNode struct {
	abstractNode

	ctype character.Type
	not bool
	asciiRange bool
}

func NewCTypeNode(ctype character.Type, not, asciiRange bool) *CTypeNode {
	return &CTypeNode{abstractNode: abstractNode{nodeType: node.CType}, ctype: ctype, not: not, asciiRange: asciiRange}
}

func (cn *CTypeNode) String() string {
	return goni.String(cn)
}

func (cn *CTypeNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	w.Append(`ctype: `)
	w.AppendInt(int(cn.ctype))
	w.Append(`, not: `)
	w.AppendBool(cn.not)
	w.Append(`, ascii: `)
	w.AppendBool(cn.asciiRange)
}

func (cn *CTypeNode) Name() string {
	return `Character Type`
}
