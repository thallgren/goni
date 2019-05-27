package ast

import (
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/enclose"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/goni/option"
	"github.com/lyraproj/goni/goni/state"
	"github.com/lyraproj/goni/util"
)

type EncloseNode struct {
	stateNode
	typ              enclose.Type
	RegNum           int
	option           option.Type
	target           goni.Node     /* EncloseNode : ENCLOSE_MEMORY */
	callAddr         int    // AbsAddrType
	minLength        int
	maxLength        int
	charLength       int  // referenced count in optimize_node_left()
	optCount         int
	ContainingAnchor goni.Node
}

func NewEncloseNode(typ enclose.Type) *EncloseNode {
	return &EncloseNode{stateNode: stateNode{abstractNode: abstractNode{nodeType: node.Enclose}}, typ: typ, callAddr: -1}
}

//noinspection GoBoolExpressions
func NewMemory(option option.Type, isNamed bool) *EncloseNode {
	en := NewEncloseNode(enclose.Memory)
	if config.UseSubExpCall {
		en.option = option
	}
	if isNamed {
		en.SetNamedGroup()
	}
	return en
}

func NewOption(option option.Type) *EncloseNode {
	en := NewEncloseNode(enclose.Option)
	en.option = option
	return en
}

func (en *EncloseNode) AppendTo(w *util.Indenter) {
	en.stateNode.AppendTo(w)
	w.NewLine()
	w.Append(`type: `)
	en.typ.AppendString(w)
	w.NewLine()
	w.Append(`regNum: `)
	w.AppendInt(en.RegNum)
	w.Append(`, option: `)
	en.option.AppendString(w)
	w.Append(`, callAddr: `)
	w.AppendInt(en.callAddr)
	w.Append(`, minLength: `)
	w.AppendInt(en.minLength)
	w.Append(`, maxLength: `)
	w.AppendInt(en.maxLength)
	w.Append(`, charLength: `)
	w.AppendInt(en.charLength)
	w.Append(`, optCount: `)
	w.AppendInt(en.optCount)
	if en.target != nil {
		w.NewLine()
		w.Append(`target: `)
		en.target.AppendTo(w.Indent())
	}
}

func (en *EncloseNode) Child() goni.Node {
	return en.target
}

func (en *EncloseNode) SetChild(child goni.Node) {
	en.target = child
}

func (en *EncloseNode) SetTarget(target goni.Node) {
	en.target = target
	target.SetParent(en)
}

func (en *EncloseNode) String() string {
	return goni.String(en)
}

func (en *EncloseNode) Name() string {
	return `EncloseNode`
}

func (en *EncloseNode) SetEncloseStatus(flag state.Type) {
	en.state |= flag
}

func (en *EncloseNode) ClearEncloseStatus(flag state.Type) {
	en.state &= ^flag
}

func (en *EncloseNode) IsType(et enclose.Type) bool {
	return en.typ.IsType(et)
}

func (en *EncloseNode) EncloseType() enclose.Type {
	return en.typ
}

type UnsetAddrList struct {
	targets []*EncloseNode
	offsets []int
}

func NewUnsetAddrList(size int) *UnsetAddrList {
	return &UnsetAddrList{targets: make([]*EncloseNode, 0, size), offsets: make([]int, 0, size)}
}

func (u* UnsetAddrList) Add(offset int, node *EncloseNode) {
	u.targets = append(u.targets, node)
	u.offsets = append(u.offsets, offset)
}

func (u* UnsetAddrList) Fix(regexCode []int) {
	for i, o := range u.offsets {
		en := u.targets[i]
		if !en.state.IsAddrFixed() {
			panic(err.NoArgs(err.ParserBug))
		}
		regexCode[o] = en.callAddr
	}
}

func (u *UnsetAddrList) String() string {
	w := util.NewIndenter()
	u.AppendTo(w)
	return w.String()
}

func (u *UnsetAddrList) AppendTo(w *util.Indenter) {
	for i, o := range u.offsets {
		w.NewLine()
		w.Append(`offset + `)
		w.AppendInt(o)
		w.Append(`target: `)
		goni.AppendAddressName(u.targets[i], w)
	}
}
