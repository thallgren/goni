package node

import "bytes"

type Anchor struct {
	node
	typ        int
	asciiRange bool
	target     Node
	charLength int
}

func NewAnchor(typ int, asciiRange bool) Node {
	return &Anchor{node: node{nodeType: TypeAnchor}, typ: typ, asciiRange: asciiRange, charLength: -1}
}

func (a *Anchor) child() Node {
	return a.target
}

func (a *Anchor) setChild(child Node) {
	a.target = child
}

func (a *Anchor) String() string {
	return NodeString(a)
}

func (a *Anchor) Name() string {
	return `Anchor`
}

func (a *Anchor) levelString(level int) string {
	value := bytes.NewBufferString(``)
	value.WriteString("\ntype: ")
	a.appendType(value)
	value.WriteString(", ascii: ")
	value.WriteString(boolToString(a.asciiRange))
	value.WriteString("\ntarget: ")
	pd := Pad(value, level)
	value.Reset()
	value.WriteString(pd)
	value.WriteString(Pad(a.target, level+1))
	return value.String()
}

func (a *Anchor) appendType(s *bytes.Buffer) {
	if a.isType(BeginBuf) {
		s.WriteString(`BEGIN_BUF `)
	}
	if a.isType(BeginLine) {
		s.WriteString(`BEGIN_LINE `)
	}
	if a.isType(BeginPosition) {
		s.WriteString(`BEGIN_POSITION `)
	}
	if a.isType(EndBuf) {
		s.WriteString(`END_BUF `)
	}
	if a.isType(SemiEndBuf) {
		s.WriteString(`SEMI_END_BUF `)
	}
	if a.isType(EndLine) {
		s.WriteString(`END_LINE `)
	}
	if a.isType(WordBound) {
		s.WriteString(`WORD_BOUND `)
	}
	if a.isType(NotWordBound) {
		s.WriteString(`NOT_WORD_BOUND `)
	}
	if a.isType(WordBegin) {
		s.WriteString(`WORD_BEGIN `)
	}
	if a.isType(WordEnd) {
		s.WriteString(`WORD_END `)
	}
	if a.isType(PrecRead) {
		s.WriteString(`PREC_READ `)
	}
	if a.isType(PrecReadNot) {
		s.WriteString(`PREC_READ_NOT `)
	}
	if a.isType(LookBehind) {
		s.WriteString(`LOOK_BEHIND `)
	}
	if a.isType(LookBeindNot) {
		s.WriteString(`LOOK_BEHIND_NOT `)
	}
	if a.isType(AnycharStar) {
		s.WriteString(`ANYCHAR_STAR `)
	}
	if a.isType(AnycharStarMl) {
		s.WriteString(`ANYCHAR_STAR_ML `)
	}
}

func (a *Anchor) isType(typ int) bool {
	return (a.typ & typ) != 0
}

func boolToString(b bool) string {
	if b {
		return `true`
	}
	return `false`
}
