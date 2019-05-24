package node

import "bytes"

type State struct {
	node
	state int
}

func (s *State) levelString(level int) string {
	b := bytes.NewBufferString("\nstate: ")
	s.appendStateString(b)
	return Pad(b, level)
}

func (s *State) appendStateString(states *bytes.Buffer) {
	if s.IsMinFixed() {
		states.WriteString(`MIN_FIXED `)
	}
	if s.IsMaxFixed() {
		states.WriteString(`MAX_FIXED `)
	}
	if s.IsMark1() {
		states.WriteString(`MARK1 `)
	}
	if s.IsMark2() {
		states.WriteString(`MARK2 `)
	}
	if s.IsMemBackrefed() {
		states.WriteString(`MEM_BACKREFED `)
	}
	if s.IsStopBtSimpleRepeat() {
		states.WriteString(`STOP_BT_SIMPLE_REPEAT `)
	}
	if s.IsRecursion() {
		states.WriteString(`RECURSION `)
	}
	if s.IsCalled() {
		states.WriteString(`CALLED `)
	}
	if s.IsAddrFixed() {
		states.WriteString(`ADDR_FIXED `)
	}
	if s.IsNamedGroup() {
		states.WriteString(`NAMED_GROUP `)
	}
	if s.IsNameRef() {
		states.WriteString(`NAME_REF `)
	}
	if s.IsInRepeat() {
		states.WriteString(`IN_REPEAT `)
	}
	if s.IsNestLevel() {
		states.WriteString(`NEST_LEVEL `)
	}
	if s.IsByNumber() {
		states.WriteString(`BY_NUMBER `)
	}
}

func (s *State) IsMinFixed() bool {
	return (s.state & StateMinFixed) != 0
}

func (s *State) SetMinFixed() {
	s.state |= StateMinFixed
}

func (s *State) ClearMinFixed() {
	s.state &= ^StateMinFixed
}

func (s *State) IsMaxFixed() bool {
	return (s.state & StateMaxFixed) != 0
}

func (s *State) SetMaxFixed() {
	s.state |= StateMaxFixed
}

func (s *State) ClearMaxFixed() {
	s.state &= ^StateMaxFixed
}

func (s *State) IsCLenFixed() bool {
	return (s.state & StateCLenFixed) != 0
}

func (s *State) SetCLenFixed() {
	s.state |= StateCLenFixed
}

func (s *State) ClearCLenFixed() {
	s.state &= ^StateCLenFixed
}

func (s *State) IsMark1() bool {
	return (s.state & StateMark1) != 0
}

func (s *State) SetMark1() {
	s.state |= StateMark1
}

func (s *State) ClearMark1() {
	s.state &= ^StateMark1
}

func (s *State) IsMark2() bool {
	return (s.state & StateMark2) != 0
}

func (s *State) SetMark2() {
	s.state |= StateMark2
}

func (s *State) ClearMark2() {
	s.state &= ^StateMark2
}

func (s *State) IsMemBackrefed() bool {
	return (s.state & StateMemBackrefed) != 0
}

func (s *State) SetMemBackrefed() {
	s.state |= StateMemBackrefed
}

func (s *State) ClearMemBackrefed() {
	s.state &= ^StateMemBackrefed
}

func (s *State) IsStopBtSimpleRepeat() bool {
	return (s.state & StateStopBtSimpleRepeat) != 0
}

func (s *State) SetStopBtSimpleRepeat() {
	s.state |= StateStopBtSimpleRepeat
}

func (s *State) ClearStopBtSimpleRepeat() {
	s.state &= ^StateStopBtSimpleRepeat
}

func (s *State) IsRecursion() bool {
	return (s.state & StateRecursion) != 0
}

func (s *State) SetRecursion() {
	s.state |= StateRecursion
}

func (s *State) ClearRecursion() {
	s.state &= ^StateRecursion
}

func (s *State) IsCalled() bool {
	return (s.state & StateCalled) != 0
}

func (s *State) SetCalled() {
	s.state |= StateCalled
}

func (s *State) ClearCAlled() {
	s.state &= ^StateCalled
}

func (s *State) IsAddrFixed() bool {
	return (s.state & StateAddrFixed) != 0
}

func (s *State) SetAddrFixed() {
	s.state |= StateAddrFixed
}

func (s *State) ClearAddrFixed() {
	s.state &= ^StateAddrFixed
}

func (s *State) IsNamedGroup() bool {
	return (s.state & StateNamedGroup) != 0
}

func (s *State) SetNamedGroup() {
	s.state |= StateNamedGroup
}

func (s *State) ClearNamedGroup() {
	s.state &= ^StateNamedGroup
}

func (s *State) IsNameRef() bool {
	return (s.state & StateNameRef) != 0
}

func (s *State) SetNameRef() {
	s.state |= StateNameRef
}

func (s *State) ClearNameRef() {
	s.state &= ^StateNameRef
}

func (s *State) IsInRepeat() bool {
	return (s.state & StateInRepeat) != 0
}

func (s *State) SetInRepeat() {
	s.state |= StateInRepeat
}

func (s *State) ClearInRepeat() {
	s.state &= ^StateInRepeat
}

func (s *State) IsNestLevel() bool {
	return (s.state & StateNestLevel) != 0
}

func (s *State) SetNestLevel() {
	s.state |= StateNestLevel
}

func (s *State) ClearNestLevel() {
	s.state &= ^StateNestLevel
}

func (s *State) IsByNumber() bool {
	return (s.state & StateByNumber) != 0
}

func (s *State) SetByNumber() {
	s.state |= StateByNumber
}

func (s *State) ClearByNumber() {
	s.state &= ^StateByNumber
}
