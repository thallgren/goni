package goni

import (
	"bytes"
	"fmt"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
	"io"
)

type Node interface {
	fmt.Stringer
	util.Indentable

	Simple() bool

	Name() string

	Parent() Node

	ReplaceWith(with Node)

	Type() node.Type

	Type2Bit() int

	SetChild(tgt Node)

	Child() Node

	SetParent(parent Node)
}

func AddressName(n Node) string {
	s := bytes.NewBufferString(``)
	AppendAddressName(n, s)
	return s.String()
}

func AppendAddressName(n Node, w io.Writer) {
	util.Fprintf(w, `%s:%p`, n.Name(), n)
}

func String(n Node) string {
	s := util.NewIndenter()
	AppendNodeString(n, s)
	return s.String()
}

func AppendNodeString(n Node, w io.Writer) {
	util.WriteString(w, `<`)
	AppendAddressName(n, w)
	util.WriteString(w, ` (`)
	if n.Parent() == nil {
		util.WriteString(w, `NULL`)
	} else {
		AppendAddressName(n.Parent(), w)
	}
	util.WriteString(w, `)>`)
}
