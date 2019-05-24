package enclose

import (
	"bytes"
	"github.com/lyraproj/goni/util"
	"io"
)

type Type int

const (
	Memory        = Type(1 << 0)
	Option        = Type(1 << 1)
	StopBacktrack = Type(1 << 2)
	Condition     = Type(1 << 3)
	Absent        = Type(1 << 4)

	AllowedInLb    = Memory | Option
	AllowedInLbNot = Option
)

func (t Type) String() string {
	w := &bytes.Buffer{}
	t.AppendString(w)
	return w.String()
}

func (t Type) AppendString(w io.Writer) {
	if t.IsStopBacktrack() {
		util.WriteString(w, `STOP_BACKTRACK `)
	}
	if t.IsMemory() {
		util.WriteString(w, `MEMORY `)
	}
	if t.IsOption() {
		util.WriteString(w, `OPTION `)
	}
	if t.IsCondition() {
		util.WriteString(w, `CONDITION `)
	}
	if t.IsAbsent() {
		util.WriteString(w, `ABSENT `)
	}
}

func (t Type) IsType(ot Type) bool {
	return (t & ot) != 0
}

func (t Type) IsMemory() bool {
	return (t & Memory) != 0
}

func (t Type) IsOption() bool {
	return (t & Option) != 0
}

func (t Type) IsCondition() bool {
	return (t & Condition) != 0
}

func (t Type) IsStopBacktrack() bool {
	return (t & StopBacktrack) != 0
}

func (t Type) IsAbsent() bool {
	return (t & Absent) != 0
}
