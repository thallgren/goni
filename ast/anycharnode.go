package ast

import (
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
)

type AnyCharNode struct {
	abstractNode
}

func (a *AnyCharNode) String() string {
	return goni.String(a)
}

func (a *AnyCharNode) Name() string {
	return `Any Char`
}

func (a *AnyCharNode) AppendTo(_ *util.Indenter) {
}

func NewAnyCharNode() goni.Node {
	return &AnyCharNode{abstractNode: abstractNode{nodeType: node.CAny}}
}
