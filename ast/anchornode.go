package ast

import (
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/anchor"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
)

type AnchorNode struct {
	abstractNode
	typ        anchor.Type
	asciiRange bool
	target     goni.Node
	charLength int
}

func NewAnchor(typ anchor.Type, asciiRange bool) goni.Node {
	return &AnchorNode{abstractNode: abstractNode{nodeType: node.Anchor}, typ: typ, asciiRange: asciiRange, charLength: -1}
}

func (a *AnchorNode) AnchorType() anchor.Type {
	return a.typ
}

func (a *AnchorNode) Name() string {
	return `AnchorNode`
}

func (a *AnchorNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	w.Append(`type: `)
	a.typ.AppendString(w)
	w.Append(`, ascii: `)
	w.AppendBool(a.asciiRange)
	if a.target != nil {
		w.NewLine()
		w.Append(`target: `)
		a.target.AppendTo(w.Indent())
	}
}

func (a *AnchorNode) String() string {
	return goni.String(a)
}

func (a *AnchorNode) Child() goni.Node {
	return a.target
}

func (a *AnchorNode) SetChild(child goni.Node) {
	a.target = child
}

func (a *AnchorNode) setTarget(target goni.Node) {
	a.target = target
	target.SetParent(a)
}
