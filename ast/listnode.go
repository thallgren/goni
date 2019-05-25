package ast

import (
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
)

type ListNode struct {
	abstractNode
	value goni.Node
	tail *ListNode
}

func newListNode(value goni.Node, tail *ListNode, typ node.Type) *ListNode {
	ln := &ListNode{abstractNode: abstractNode{nodeType: typ}, value: value, tail: tail}
	if value != nil {
		value.SetParent(ln)
	}
	if tail != nil {
		tail.SetParent(ln)
	}
	return ln
}

func NewList(value goni.Node, tail *ListNode) *ListNode {
	return newListNode(value, tail, node.List)
}

func NewAlt(value goni.Node, tail *ListNode) *ListNode {
	return newListNode(value, tail, node.Alt)
}

func ListAdd(list *ListNode, value goni.Node) *ListNode {
	n := newListNode(value, nil, node.List)
	if list != nil {
		for list.tail != nil {
			list = list.tail
		}
		list.tail = n
	}
	return n
}

func (ln *ListNode) ToListNode() {
	ln.nodeType = node.List
}

func (ln *ListNode) Child() goni.Node {
	return ln.value
}

func (ln *ListNode) SetChild(child goni.Node) {
	ln.value = child
}

func (ln *ListNode) setValue(value goni.Node) {
	ln.value = value
	value.SetParent(ln)
}

func (ln *ListNode) setTail(tail *ListNode) {
	ln.tail = tail
}

func (ln *ListNode) String() string {
	return goni.String(ln)
}

func (ln *ListNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	w.Append(`value: `)
	if ln.value == nil {
		w.Append(`NULL`)
	} else {
		ln.value.AppendTo(w.Indent())
	}
	w.Append(`tail: `)
	if ln.tail == nil {
		w.Append(`NULL`)
	} else {
		ln.tail.AppendTo(w.Indent())
	}
}

func (ln *ListNode) Name() string {
	if ln.nodeType == node.List {
		return `List`
	}
	return `Alt`
}
