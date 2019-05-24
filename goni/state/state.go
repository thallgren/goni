package state

import (
	"github.com/lyraproj/goni/util"
	"io"
)

type Type int
const (
	MinFixed           = Type(1 << 0)
	MaxFixed           = Type(1 << 1)
	CLenFixed          = Type(1 << 2)
	Mark1              = Type(1 << 3)
	Mark2              = Type(1 << 4)
	MemBackRefed       = Type(1 << 5)
	StopBtSimpleRepeat = Type(1 << 6)
	Recursion          = Type(1 << 7)
	Called             = Type(1 << 8)
	AddrFixed          = Type(1 << 9)
	NamedGroup         = Type(1 << 10)
	NameRef            = Type(1 << 11)
	InRepeat           = Type(1 << 12) /* STK_REPEAT is nested in stack. */
	NestLevel          = Type(1 << 13)
	ByNumber           = Type(1 << 14) /* {n,m} */
)

func (s Type) AppendString(w io.Writer) {
	if s.IsMinFixed() {
		util.WriteString(w, `MIN_FIXED `)
	}
	if s.IsMaxFixed() {
		util.WriteString(w, `MAX_FIXED `)
	}
	if s.IsMark1() {
		util.WriteString(w, `MARK1 `)
	}
	if s.IsMark2() {
		util.WriteString(w, `MARK2 `)
	}
	if s.IsMemBackRefed() {
		util.WriteString(w, `MEM_BACKREFED `)
	}
	if s.IsStopBtSimpleRepeat() {
		util.WriteString(w, `STOP_BT_SIMPLE_REPEAT `)
	}
	if s.IsRecursion() {
		util.WriteString(w, `RECURSION `)
	}
	if s.IsCalled() {
		util.WriteString(w, `CALLED `)
	}
	if s.IsAddrFixed() {
		util.WriteString(w, `ADDR_FIXED `)
	}
	if s.IsNamedGroup() {
		util.WriteString(w, `NAMED_GROUP `)
	}
	if s.IsNameRef() {
		util.WriteString(w, `NAME_REF `)
	}
	if s.IsInRepeat() {
		util.WriteString(w, `IN_REPEAT `)
	}
	if s.IsNestLevel() {
		util.WriteString(w, `NEST_LEVEL `)
	}
	if s.IsByNumber() {
		util.WriteString(w, `BY_NUMBER `)
	}
}

func (t Type) IsType(ot Type) bool {
	return (t & ot) != 0
}

func (s Type) IsMinFixed() bool {
	return (s & MinFixed) != 0
}

func (s Type) IsMaxFixed() bool {
	return (s & MaxFixed) != 0
}

func (s Type) IsCLenFixed() bool {
	return (s & CLenFixed) != 0
}

func (s Type) IsMark1() bool {
	return (s & Mark1) != 0
}

func (s Type) IsMark2() bool {
	return (s & Mark2) != 0
}

func (s Type) IsMemBackRefed() bool {
	return (s & MemBackRefed) != 0
}

func (s Type) IsStopBtSimpleRepeat() bool {
	return (s & StopBtSimpleRepeat) != 0
}

func (s Type) IsRecursion() bool {
	return (s & Recursion) != 0
}

func (s Type) IsCalled() bool {
	return (s & Called) != 0
}

func (s Type) IsAddrFixed() bool {
	return (s & AddrFixed) != 0
}

func (s Type) IsNamedGroup() bool {
	return (s & NamedGroup) != 0
}

func (s Type) IsNameRef() bool {
	return (s & NameRef) != 0
}

func (s Type) IsInRepeat() bool {
	return (s & InRepeat) != 0
}

func (s Type) IsNestLevel() bool {
	return (s & NestLevel) != 0
}

func (s Type) IsByNumber() bool {
	return (s & ByNumber) != 0
}
