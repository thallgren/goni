package node

import "github.com/lyraproj/goni/goni"

type BackRef struct {
	State
	back []int
	backNum int
	nestLevel int
}

func NewBackRef(backNum int, backRefs []int, byName bool, env *goni.ScanEnvironment) Node {
	b := &BackRef{State: State{node: node{nodeType: TypeBref}},
		back: backRefs, backNum: backNum}
	if byName {
		b.SetNameRef()
	}

	for i := 0; i < backNum; i++ {
		if backRefs[i] <= env.NumMem() && env.MemNodes()[backRefs[i]] == nil {
			b.SetRecursion(); /* /...(\1).../ */
			break;
		}
	}
	return b
}

func NewBackRef2(backNum int, backRefs []int, byName, existLevel bool, nestLevel int, env *goni.ScanEnvironment) Node {
	b := NewBackRef(backNum, backRefs, byName, env)
	if goni.UseBackrefWithLevel && existLevel {
		bi := b.(*BackRef)
		bi.SetNestLevel()
		bi.nestLevel = nestLevel
	}
	return b
}
func (b *BackRef) Back() []int {
	return b.back
}

func (b *BackRef) BackNum() int {
	return b.backNum
}

func (b *BackRef) NestLevel() int {
	return b.nestLevel
}

func (b *BackRef) Renumber(m []int) {
	// TODO: Implement
}

func (b *BackRef) String() string {
	return NodeString(b)
}

func (b *BackRef) Name() string {
	return `Back Ref`
}

