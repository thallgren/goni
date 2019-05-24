package ast

import (
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
)

type CallNode struct {
	stateNode
	name []byte
	nameP int
	nameEnd int

	groupNum int
	target *EncloseNode
	unsetAddrList *UnsetAddrList
}

func (c *CallNode) String() string {
	return goni.String(c)
}

func (c *CallNode) Name() string {
	return `CallNode`
}

func NewCall(name []byte, nameP, nameEnd, gnum int) *CallNode {
	return &CallNode{
		stateNode: stateNode{abstractNode: abstractNode{nodeType: node.Call}},
		name:      name, nameP: nameP, nameEnd: nameEnd, groupNum: gnum}
}

func (c *CallNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	w.Append(`name: `)
	w.Append(string(c.name[c.nameP:c.nameEnd]))
	w.Append(`, groupNum: `)
	w.AppendInt(c.groupNum)
	if c.unsetAddrList != nil {
		w.NewLine()
		w.Append(`unsetAddrList: `)
		c.unsetAddrList.AppendTo(w.Indent())
	}
	if c.target != nil {
		w.NewLine()
		w.Append(`target: `)
		c.target.AppendTo(w.Indent())
	}
}

func (c *CallNode) Child() goni.Node {
	return c.target
}

func (c *CallNode) SetChild(child goni.Node) {
	c.target = child.(*EncloseNode)
}

func (c *CallNode) setTarget(target *EncloseNode) {
	c.target = target
	target.SetParent(c)
}
