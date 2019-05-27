package internal

import (
	"github.com/lyraproj/goni/goni/anchor"
	"github.com/lyraproj/goni/goni/character"
)

type TokenType int

const (
	TkEOT = TokenType(iota) /* end of token */
	TkRawByte
	TkChar
	TkString
	TkCodePoint
	TkAnyChar
	TkCharType
	TkBackRef
	TkCall
	TkAnchor
	TkOpRepeat
	TkInterval
	TkAnycharAnytime /* SQL '%' == .* */
	TkAlt
	TkSubexpOpen
	TkSubexpClose
	TkCcOpen
	TkQuoteOpen
	TkCharProperty /* \p{...}, \P{...} */
	TkLineBreak
	TkExtendedGraphemeCluster
	TkKeep
	/* in cc */
	TkCcClose
	TkCcRange
	TkPosixBracketOpen
	TkCcAnd    /* && */
	TkCcCcOpen /* [ */
)

type Token struct {
	Type    TokenType
	escaped bool
	base    int
	backP   int

	// union fields
	int1 int
	int2 int
	int3 int
	int4 int
	int5 int
	inta1 []int
}

// union accessors
func (t *Token) getC() int {
	return t.int1
}
func (t *Token) setC(c int) {
	t.int1 = c
}
func (t *Token) getCode() int {
	return t.int1
}
func (t *Token) setCode(code int) {
	t.int1 = code
}
func (t *Token) getAnchorSubtype() anchor.Type {
	return anchor.Type(t.int1)
}
func (t *Token) setAnchorSubtype(anchor anchor.Type) {
	t.int1 = int(anchor)
}
func (t *Token) getAnchorASCIIRange() bool {
	return t.int2 == 1
}
func (t *Token) setAnchorASCIIRange(ascii bool) {
	if ascii {
		t.int2 = 1
	} else {
		t.int2 = 0
	}
}

// repeat union member
func (t *Token) getRepeatLower() int {
	return t.int1
}
func (t *Token) setRepeatLower(lower int) {
	t.int1 = lower
}
func (t *Token) getRepeatUpper() int {
	return t.int2
}
func (t *Token) setRepeatUpper(upper int) {
	t.int2 = upper
}
func (t *Token) getRepeatGreedy() bool {
	return t.int3 != 0
}
func (t *Token) setRepeatGreedy(greedy bool) {
	if greedy {
		t.int3 = 1
	} else {
		t.int3 = 0
	}
}
func (t *Token) getRepeatPossessive() bool {
	return t.int4 != 0
}
func (t *Token) setRepeatPossessive(possessive bool) {
	if possessive {
		t.int4 = 1
	} else {
		t.int4 = 0
	}
}

// backref union member
func (t *Token) getBackrefNum() int {
	return t.int1
}
func (t *Token) setBackrefNum(num int) {
	t.int1 = num
}
func (t *Token) getBackrefRef1() int {
	return t.int2
}
func (t *Token) setBackrefRef1(ref1 int) {
	t.int2 = ref1
}

func (t *Token) getBackrefRefs() []int {
	return t.inta1
}
func (t *Token) setBackrefRefs(refs []int) {
	t.inta1 = refs
}
func (t *Token) getBackrefByName() bool {
	return t.int3 != 0
}
func (t *Token) setBackrefByName(byName bool) {
	if byName {
		t.int3 = 1
	} else {
		t.int3 = 0
	}
}

// USE_BACKREF_AT_LEVEL
func (t *Token) getBackrefExistLevel() bool {
	return t.int4 != 0
}
func (t *Token) setBackrefExistLevel(existLevel bool) {
	if existLevel {
		t.int4 = 1
	} else {
		t.int4 = 0
	}
}
func (t *Token) getBackrefLevel() int {
	return t.int5
}
func (t *Token) setBackrefLevel(level int) {
	t.int5 = level
}

// call union member
func (t *Token) getCallNameP() int {
	return t.int1
}
func (t *Token) setCallNameP(nameP int) {
	t.int1 = nameP
}
func (t *Token) getCallNameEnd() int {
	return t.int2
}
func (t *Token) setCallNameEnd(nameEnd int) {
	t.int2 = nameEnd
}
func (t *Token) getCallGNum() int {
	return t.int3
}
func (t *Token) setCallGNum(gnum int) {
	t.int3 = gnum
}
func (t *Token) getCallRel() bool {
	return t.int4 != 0
}
func (t *Token) setCallRel(rel bool) {
	if rel {
		t.int4 = 1
	} else {
		t.int4 = 0
	}
}

// prop union member
func (t *Token) getPropCType() character.Type {
	return character.Type(t.int1)
}
func (t *Token) setPropCType(ctype character.Type) {
	t.int1 = int(ctype)
}
func (t *Token) getPropNot() bool {
	return t.int2 != 0
}
func (t *Token) setPropNot(not bool) {
	if not {
		t.int2 = 1
	} else {
		t.int2 = 0
	}
}
