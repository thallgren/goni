package ast

import (
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/util"
)

type BackRefNode struct {
	stateNode
	back      []int
	backNum   int
	nestLevel int
}

func NewBackRef(backNum int, backRefs []int, byName bool, env goni.ScanEnvironment) goni.Node {
	b := &BackRefNode{stateNode: stateNode{abstractNode: abstractNode{nodeType: node.Bref}},
		back: backRefs, backNum: backNum}
	if byName {
		b.SetNameRef()
	}

	for i := 0; i < backNum; i++ {
		if backRefs[i] <= env.NumMem() && env.MemNodes()[backRefs[i]] == nil {
			b.SetRecursion() /* /...(\1).../ */
			break
		}
	}
	return b
}

//noinspection GoBoolExpressions
func NewBackRef2(backNum int, backRefs []int, byName, existLevel bool, nestLevel int, env goni.ScanEnvironment) goni.Node {
	b := NewBackRef(backNum, backRefs, byName, env)
	if config.UseBackrefWithLevel && existLevel {
		bi := b.(*BackRefNode)
		bi.SetNestLevel()
		bi.nestLevel = nestLevel
	}
	return b
}

func (b *BackRefNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	w.Append(`backNum: `)
	w.AppendInt(b.backNum)
	w.NewLine()
	w.Append(`back: `)
	for i, bk := range b.back {
		if i > 0 {
			w.Append(`, `)
		}
		w.AppendInt(bk)
	}
	w.NewLine()
	w.Append(`nestLevel: `)
	w.AppendInt(b.nestLevel)
}

func (b *BackRefNode) Back() []int {
	return b.back
}

func (b *BackRefNode) BackNum() int {
	return b.backNum
}

func (b *BackRefNode) NestLevel() int {
	return b.nestLevel
}

func (b *BackRefNode) Renumber(m []int) {
	if !b.state.IsNameRef() {
		panic(err.NoArgs(err.NumberedBackrefOrCallNotAllowed))
	}

	pos := 0
	for _, bk := range b.back {
		n := m[bk]
		if n > 0 {
			b.back[pos] = n
			pos++
		}
	}
	b.backNum = pos
}

func (b *BackRefNode) String() string {
	return goni.String(b)
}

func (b *BackRefNode) Name() string {
	return `Back Ref`
}
