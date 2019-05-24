package ast

import (
	"github.com/lyraproj/goni/goni/state"
	"github.com/lyraproj/goni/util"
)

type stateNode struct {
	abstractNode
	state state.Type
}

func (s *stateNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	w.Append(`state: `)
	s.state.AppendString(w)
}

func (s *stateNode) StateType() state.Type {
	return s.state
}

func (s *stateNode) SetMinFixed() {
	s.state |= state.MinFixed
}

func (s *stateNode) ClearMinFixed() {
	s.state &= ^state.MinFixed
}

func (s *stateNode) SetMaxFixed() {
	s.state |= state.MaxFixed
}

func (s *stateNode) ClearMaxFixed() {
	s.state &= ^state.MaxFixed
}

func (s *stateNode) SetCLenFixed() {
	s.state |= state.CLenFixed
}

func (s *stateNode) ClearCLenFixed() {
	s.state &= ^state.CLenFixed
}

func (s *stateNode) SetMark1() {
	s.state |= state.Mark1
}

func (s *stateNode) ClearMark1() {
	s.state &= ^state.Mark1
}

func (s *stateNode) SetMark2() {
	s.state |= state.Mark2
}

func (s *stateNode) ClearMark2() {
	s.state &= ^state.Mark2
}

func (s *stateNode) SetMemBackRefed() {
	s.state |= state.MemBackRefed
}

func (s *stateNode) ClearMemBackRefed() {
	s.state &= ^state.MemBackRefed
}

func (s *stateNode) SetStopBtSimpleRepeat() {
	s.state |= state.StopBtSimpleRepeat
}

func (s *stateNode) ClearStopBtSimpleRepeat() {
	s.state &= ^state.StopBtSimpleRepeat
}

func (s *stateNode) SetRecursion() {
	s.state |= state.Recursion
}

func (s *stateNode) ClearRecursion() {
	s.state &= ^state.Recursion
}

func (s *stateNode) SetCalled() {
	s.state |= state.Called
}

func (s *stateNode) ClearCalled() {
	s.state &= ^state.Called
}

func (s *stateNode) SetAddrFixed() {
	s.state |= state.AddrFixed
}

func (s *stateNode) ClearAddrFixed() {
	s.state &= ^state.AddrFixed
}

func (s *stateNode) SetNamedGroup() {
	s.state |= state.NamedGroup
}

func (s *stateNode) ClearNamedGroup() {
	s.state &= ^state.NamedGroup
}

func (s *stateNode) SetNameRef() {
	s.state |= state.NameRef
}

func (s *stateNode) ClearNameRef() {
	s.state &= ^state.NameRef
}

func (s *stateNode) SetInRepeat() {
	s.state |= state.InRepeat
}

func (s *stateNode) ClearInRepeat() {
	s.state &= ^state.InRepeat
}

func (s *stateNode) SetNestLevel() {
	s.state |= state.NestLevel
}

func (s *stateNode) ClearNestLevel() {
	s.state &= ^state.NestLevel
}

func (s *stateNode) SetByNumber() {
	s.state |= state.ByNumber
}

func (s *stateNode) ClearByNumber() {
	s.state &= ^state.ByNumber
}
