package ast

import (
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
)

func NewTop(root goni.Node) goni.Node {
	t := &topNode{abstractNode{nil, node.Top}, nil}
	t.setChild(root)
	return t
}

type abstractNode struct {
	parent   goni.Node
	nodeType node.Type
}

func (n *abstractNode) SetChild(tgt goni.Node) {
	// default implementation
}

func (n *abstractNode) Child() goni.Node {
	// default implementation
	return nil
}

func (n *abstractNode) ReplaceWith(with goni.Node) {
	with.SetParent(n.parent)
	n.parent.SetChild(with)
	n.parent = nil
}

func (n *abstractNode) Parent() goni.Node {
	return n.parent
}

func (n *abstractNode) Simple() bool {
	return (n.Type2Bit() & node.Simple) != 0
}

func (n *abstractNode) Type() node.Type {
	return n.nodeType
}

func (n *abstractNode) Type2Bit() int {
	return 1 << n.nodeType
}

func (n *abstractNode) SetParent(parent goni.Node) {
	n.parent = parent
}

type topNode struct {
	abstractNode
	root goni.Node
}

func (t *topNode) String() string {
	panic("implement me")
}

func (t *topNode) Name() string {
	return `ROOT`
}

func (t *topNode) setChild(node goni.Node) {
	node.SetParent(t)
	t.root = node
}

func (t *topNode) child() goni.Node {
	return t.root
}

func (t *topNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	t.root.AppendTo(w)
}
