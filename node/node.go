package node

import (
	"bytes"
	"fmt"
	"strings"
)

type Node interface {
	fmt.Stringer

	Simple() bool

	Name() string

	Parent() Node

	ReplaceWith(with Node)

	Type() Type

	Type2Bit() int

	setChild(tgt Node)

	child() Node

	levelString(level int) string

	setParent(parent Node)
}

func AddressName(n Node) string {
	s := bytes.NewBufferString(``);
	AppendAddressName(n, s)
	return s.String()
}

func AppendAddressName(n Node, b *bytes.Buffer) {
	fmt.Fprintf(b, `%s:%p`, n.Name(), n);
}

func NodeString(n Node) string {
	s := bytes.NewBufferString(``);
	AppendNodeString(n, s)
	return s.String()
}

func AppendNodeString(n Node, s *bytes.Buffer) {
	s.WriteString(`<`)
	AppendAddressName(n, s)
	s.WriteString(` (`)
	if n.Parent() == nil {
		s.WriteString(`NULL`)
	} else {
		AppendAddressName(n.Parent(), s)
	}
	s.WriteString(`)>`)
}

func NewTop(root Node) Node {
	t := &topNode{node{nil, TypeTop}, nil}
	t.setChild(root)
	return t
}

func Pad(value fmt.Stringer, level int) string {
	if value == nil {
		return `NULL`
	}
	vs := value.String()
	if level == 0 {
		return vs
	}

	s := bytes.NewBufferString("\n")
	for i := 0; i < level; i++ {
		s.WriteString(`  `);
	}
	return strings.ReplaceAll(vs, "\n", s.String());
}

type node struct {
	parent   Node
	nodeType Type
}

func (n *node) setChild(tgt Node) {
	// default implementation
}

func (n *node) child() Node {
	// default implementation
	return nil
}

func (n *node) ReplaceWith(with Node) {
	with.setParent(n.parent)
	n.parent.setChild(with)
	n.parent = nil
}

func (n *node) Parent() Node {
	return n.parent
}

func (n *node) Simple() bool {
	return (n.Type2Bit() & BitSimple) != 0
}

func (n *node) Type() Type {
	return n.nodeType
}

func (n *node) Type2Bit() int {
	return 1 << n.nodeType
}

func (n *node) setParent(parent Node) {
	n.parent = parent
}

type topNode struct {
	node
	root Node
}

func (t *topNode) String() string {
	panic("implement me")
}

func (t * topNode) Name() string {
	return `ROOT`
}

func (t *topNode) setChild(node Node) {
	node.setParent(t)
	t.root = node
}

func (t *topNode) child() Node {
	return t.root
}

func (t * topNode) levelString(level int) string {
	return "\n" + Pad(t.root, level)
}
