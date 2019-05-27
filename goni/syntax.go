package goni

import (
	"github.com/lyraproj/goni/goni/metachar"
	"github.com/lyraproj/goni/goni/option"
	"github.com/lyraproj/goni/goni/syntax"
)

type Syntax struct {
	name string
	op syntax.Op
	op2 syntax.Op2
	op3 syntax.Op3
	behavior syntax.Behavior
	options option.Type
	MetaCharTable *metachar.Table
}

func (s *Syntax) IsOp(op syntax.Op) bool {
	return s.op.IsSet(op)
}

func (s *Syntax) IsOp2(op syntax.Op2) bool {
	return s.op2.IsSet(op)
}

func (s *Syntax) IsOp3(op syntax.Op3) bool {
	return s.op3.IsSet(op)
}

func (s *Syntax) IsBehavior(op syntax.Behavior) bool {
	return s.behavior.IsSet(op)
}
