package internal

import (
	"github.com/lyraproj/goni/ast"
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/anchor"
	"github.com/lyraproj/goni/goni/character"
	"github.com/lyraproj/goni/goni/metachar"
	"github.com/lyraproj/goni/goni/syntax"
	"github.com/lyraproj/issue/issue"
	"strings"
)

type Lexer struct {
	scannerSupport
	regex  *Regex
	env    goni.ScanEnvironment
	syntax *goni.Syntax
	token  Token
}

func (lx *Lexer) init(regex *Regex, syntax *goni.Syntax, bytes []byte, p, end int, warnings WarnCallback) {
	lx.scannerSupport.init(regex.enc, bytes, p, end)
	lx.regex = regex
	lx.env = NewScanEnvironment(regex, syntax, warnings)
	lx.syntax = syntax
}

/**
 * @return 0: normal {n,m}, 2: fixed {n}
 * !introduce returnCode here
 */
func (lx *Lexer) fetchRangeQuantifier() int {
	lx.mark()
	snx := lx.syntax
	synAllow := snx.IsBehavior(syntax.AllowInvalidInterval)

	if !lx.left() {
		if synAllow {
			return 1 /* "....{" : OK! */
		} else {
			panic(newSyntaxException(err.EndPatternAtLeftBrace))
		}
	}

	if !synAllow {
		lx.c = lx.peek()
		if lx.c == ')' || lx.c == '(' || lx.c == '|' {
			panic(newSyntaxException(err.EndPatternAtLeftBrace))
		}
	}

	low := lx.scanUnsignedNumber()
	if low < 0 || low > config.MaxRepeatNum {
		panic(newSyntaxException(err.TooBigNumberForRepeatRange))
	}

	nonLow := false
	if lx.p == lx._p { /* can't read low */
		if snx.IsBehavior(syntax.AllowIntervalLowAbbrev) {
			low = 0
			nonLow = true
		} else {
			return lx.invalidRangeQuantifier(synAllow)
		}
	}

	if !lx.left() {
		return lx.invalidRangeQuantifier(synAllow)
	}

	lx.fetch()
	var up int
	ret := 0
	if lx.c == ',' {
		prev := lx.p // ??? last
		up := lx.scanUnsignedNumber()
		if up < 0 || up > config.MaxRepeatNum {
			panic(newSyntaxException(err.TooBigNumberForRepeatRange))
		}

		if lx.p == prev {
			if nonLow {
				return lx.invalidRangeQuantifier(synAllow)
			}
			up = ast.QuantifierRepeatInfinite /* {n,} : {n,infinite} */
		}
	} else {
		if nonLow {
			return lx.invalidRangeQuantifier(synAllow)
		}
		lx.unfetch()
		up = low /* {n} : exact n times */
		ret = 2  /* fixed */
	}

	if !lx.left() {
		return lx.invalidRangeQuantifier(synAllow)
	}
	lx.fetch()

	if snx.IsOp(syntax.OpEscBraceInterval) {
		if lx.c != snx.MetaCharTable.Esc {
			return lx.invalidRangeQuantifier(synAllow)
		}
		lx.fetch()
	}

	if lx.c != '}' {
		return lx.invalidRangeQuantifier(synAllow)
	}

	if up != ast.QuantifierRepeatInfinite && low > up {
		panic(newSyntaxException(err.UpperSmallerThanLowerInRepeatRange))
	}

	t := &lx.token
	t.typ = TkInterval
	t.setRepeatLower(low)
	t.setRepeatUpper(up)

	return ret /* 0: normal {n,m}, 2: fixed {n} */
}

func (lx *Lexer) invalidRangeQuantifier(synAllow bool) int {
	if synAllow {
		lx.restore()
		return 1
	}
	panic(newSyntaxException(err.InvalidRepeatRangePattern))
}

/* \M-, \C-, \c, or \... */
func (lx *Lexer) fetchEscapedValue() {
	if !lx.left() {
		panic(newSyntaxException(err.EndPatternAtEscape))
	}
	lx.fetch()

	snx := lx.syntax
	switch lx.c {

	case 'M':
		if snx.IsOp2(syntax.Op2EscCapitalMBarMeta) {
			if !lx.left() {
				panic(newSyntaxException(err.EndPatternAtMeta))
			}
			lx.fetch()
			if lx.c != '-' {
				panic(newSyntaxException(err.MetaCodeSyntax))
			}
			if !lx.left() {
				panic(newSyntaxException(err.EndPatternAtMeta))
			}
			lx.fetch()
			if lx.c == snx.MetaCharTable.Esc {
				lx.fetchEscapedValue()
			}
			lx.c = (lx.c & 0xff) | 0x80
		} else {
			lx.fetchEscapedValueBackSlash()
		}

	case 'C':
		if snx.IsOp2(syntax.Op2EscCapitalCBarControl) {
			if !lx.left() {
				panic(newSyntaxException(err.EndPatternAtControl))
			}
			lx.fetch()
			if lx.c != '-' {
				panic(newSyntaxException(err.ControlCodeSyntax))
			}
			lx.fetchEscapedValueControl()
		} else {
			lx.fetchEscapedValueBackSlash()
		}

	case 'c':
		if snx.IsOp(syntax.OpEscCControl) {
			lx.fetchEscapedValueControl()
		}
		fallthrough
	default:
		lx.fetchEscapedValueBackSlash()
	}
}

func (lx *Lexer) fetchEscapedValueBackSlash() {
	lx.c = lx.env.ConvertBackslashValue(lx.c)
}

func (lx *Lexer) fetchEscapedValueControl() {
	snx := lx.env.Syntax()
	if !lx.left() {
		if snx.IsOp3(syntax.Op3OptionECMAScript) {
			return
		} else {
			panic(newSyntaxException(err.EndPatternAtControl))
		}
	}
	lx.fetch()
	if lx.c == '?' {
		lx.c = 0177
	} else {
		if lx.c == snx.MetaCharTable.Esc {
			lx.fetchEscapedValue()
		}
		lx.c &= 0x9f
	}
}
func (lx *Lexer) nameEndCodePoint(start int) int {
	switch start {
	case '<':
		return '>'
	case '\'':
		return '\''
	case '(':
		return ')'
	case '{':
		return '}'
	default:
		return 0
	}
}

// USE_NAMED_GROUP && USE_BACKREF_AT_LEVEL
/*
   \k<name+n>, \k<name-n>
   \k<num+n>,  \k<num-n>
   \k<-num+n>, \k<-num-n>
*/

func (lx *Lexer) fetchNameWithLevel(startCode int) (existLevel bool, backnum, level int) {
	src := lx.p
	isNum := 0
	sign := 1

	endCode := lx.nameEndCodePoint(startCode)
	pnumHead := lx.p
	nameEnd := lx.stop

	var er issue.Code
	if !lx.left() {
		panic(newSyntaxException(err.EmptyGroupName))
	}
	lx.fetch()
	c := lx.c
	if c == endCode {
		panic(newSyntaxException(err.EmptyGroupName))
	}
	if lx.enc.IsDigit(c) {
		isNum = 1
	} else if c == '-' {
		isNum = 2
		sign = -1
		pnumHead = lx.p
	}

	for lx.left() {
		nameEnd = lx.p
		lx.fetch()
		c := lx.c
		if c == endCode || c == ')' || c == '+' || c == '-' {
			if isNum == 2 {
				er = err.InvalidGroupName
			}
			break
		}

		if isNum != 0 {
			if lx.enc.IsDigit(c) {
				isNum = 1
			} else {
				er = err.InvalidGroupName
				// isNum = 0;
			}
		}
	}

	isEndCode := false
	if er == `` && lx.c != endCode {
		if lx.c == '+' || lx.c == '-' {
			flag := 1
			if lx.c == '-' {
				flag = -1
			}

			lx.fetch()
			if !lx.enc.IsDigit(lx.c) {
				panic(lx.newValueException(err.InvalidGroupName, src, lx.stop))
			}
			lx.unfetch()
			level := lx.scanUnsignedNumber()
			if level < 0 {
				panic(newSyntaxException(err.TooBigNumber))
			}
			level = level * flag
			existLevel = true

			lx.fetch()
			isEndCode = lx.c == endCode
		}

		if !isEndCode {
			er = err.InvalidGroupName
			nameEnd = lx.stop
		}
	}

	if er != `` {
		panic(lx.newValueException(err.InvalidGroupName, src, lx.stop))
	}

	if isNum != 0 {
		lx.mark()
		lx.p = pnumHead
		backNum := lx.scanUnsignedNumber()
		lx.restore()
		if backNum < 0 {
			panic(newSyntaxException(err.TooBigNumber))
		}
		if backNum == 0 {
			panic(lx.newValueException(err.InvalidGroupName, src, lx.stop))
		}
		backNum = backNum * sign
	}
	lx.value = nameEnd
	return
}

// USE_NAMED_GROUP
// ref: 0 -> define name    (don't allow number name)
//      1 -> reference name (allow number name)
func (lx *Lexer) fetchNameForNamedGroup(startCode int, ref bool) int {
	src := lx.p
	enc := lx.enc
	lx.value = 0

	isNum := 0
	sign := 1

	endCode := lx.nameEndCodePoint(startCode)
	pnumHead := lx.p
	nameEnd := lx.stop

	var er issue.Code
	if !lx.left() {
		panic(newSyntaxException(err.EmptyGroupName))
	}
	lx.fetch()
	c := lx.c
	if c == endCode {
		panic(newSyntaxException(err.EmptyGroupName))
	}
	if enc.IsDigit(c) {
		if ref {
			isNum = 1
		} else {
			er = err.InvalidGroupName
			// isNum = 0;
		}
	} else if c == '-' {
		if ref {
			isNum = 2
			sign = -1
			pnumHead = lx.p
		} else {
			er = err.InvalidGroupName
			// isNum = 0;
		}
	}

	if er != `` {
		return lx.fetchNameTeardown(src, endCode, nameEnd, er)
	}
	for lx.left() {
		nameEnd = lx.p
		lx.fetch()
		c := lx.c
		if c == endCode || c == ')' {
			if isNum == 2 {
				er = err.InvalidGroupName
				return lx.fetchNameTeardown(src, endCode, nameEnd, er)
			}
			break
		}

		if isNum != 0 {
			if enc.IsDigit(c) {
				isNum = 1
			} else {
				if !enc.IsWord(c) {
					er = err.InvalidCharInGroupName
				} else {
					er = err.InvalidGroupName
				}
				return lx.fetchNameTeardown(src, endCode, nameEnd, er)
			}
		}
	}

	if lx.c != endCode {
		er = err.InvalidGroupName
		nameEnd = lx.stop
		return lx.fetchNameErr(src, nameEnd, er)
	}

	backNum := 0
	if isNum != 0 {
		lx.mark()
		lx.p = pnumHead
		backNum = lx.scanUnsignedNumber()
		lx.restore()
		if backNum < 0 {
			panic(newSyntaxException(err.TooBigNumber))
		}
		if backNum == 0 {
			panic(lx.newValueException(err.InvalidGroupName, src, nameEnd))
		}
		backNum *= sign
	}
	lx.value = nameEnd
	return backNum
}

func (lx *Lexer) fetchNameErr(src, nameEnd int, er issue.Code) int {
	panic(lx.newValueException(er, src, nameEnd))
}

func (lx *Lexer) fetchNameTeardown(src, endCode, nameEnd int, er issue.Code) int {
	for lx.left() {
		nameEnd = lx.p
		lx.fetch()
		if lx.c == endCode || lx.c == ')' {
			break
		}
	}
	if !lx.left() {
		nameEnd = lx.stop
	}
	return lx.fetchNameErr(src, nameEnd, er)
}

// #else USE_NAMED_GROUP
// make it return nameEnd!
func (lx *Lexer) fetchNameForNoNamedGroup(startCode int, ref bool) int {
	src := lx.p
	enc := lx.enc
	lx.value = 0
	sign := 1

	endCode := lx.nameEndCodePoint(startCode)
	pnumHead := lx.p
	nameEnd := lx.stop

	var er issue.Code
	if !lx.left() {
		panic(newSyntaxException(err.EmptyGroupName))
	}
	lx.fetch()
	c := lx.c
	if c == endCode {
		panic(newSyntaxException(err.EmptyGroupName))
	}

	if enc.IsDigit(c) {
	} else if c == '-' {
		sign = -1
		pnumHead = lx.p
	} else {
		er = err.InvalidCharInGroupName
	}

	for lx.left() {
		nameEnd = lx.p

		lx.fetch()
		c := lx.c
		if c == endCode || c == ')' {
			break
		}
		if !enc.IsDigit(c) {
			er = err.InvalidCharInGroupName
		}
	}

	if er == `` && lx.c != endCode {
		er = err.InvalidGroupName
		nameEnd = lx.stop
	}

	if er != `` {
		panic(lx.newValueException(er, src, nameEnd))
	}
	lx.mark()
	lx.p = pnumHead
	backNum := lx.scanUnsignedNumber()
	lx.restore()
	if backNum < 0 {
		panic(newSyntaxException(err.TooBigNumber))
	}
	if backNum == 0 {
		panic(lx.newValueException(err.InvalidGroupName, src, nameEnd))
	}
	backNum *= sign

	lx.value = nameEnd
	return backNum
}

func (lx *Lexer) fetchName(startCode int, ref bool) int {
	//noinspection GoBoolExpressions
	if config.UseNamedGroup {
		return lx.fetchNameForNamedGroup(startCode, ref)
	} else {
		return lx.fetchNameForNoNamedGroup(startCode, ref)
	}
}

func (lx *Lexer) strExistCheckWithEsc(s []int, n, bad int) bool {
	p := lx.p
	to := lx.stop
	enc := lx.enc

	inEsc := false
	i := 0
	for p < to {
		if inEsc {
			inEsc = false
			p += enc.Length(lx.bytes, p, to)
		} else {
			x, cl := enc.MbcToCode(lx.bytes, p, to)
			q := p + cl
			if x == s[0] {
				for i = 1; i < n && q < to; i++ {
					x, cl = enc.MbcToCode(lx.bytes, q, to)
					if x != s[i] {
						break
					}
					q += cl
				}
				if i >= n {
					return true
				}
				p += enc.Length(lx.bytes, p, to)
			} else {
				x, _ = enc.MbcToCode(lx.bytes, p, to)
				if x == bad {
					return false
				} else if x == lx.syntax.MetaCharTable.Esc {
					inEsc = true
				}
				p = q
			}
		}
	}
	return false
}

var send = []int{':', ']'}

func (lx *Lexer) fetchTokenInCCForCharType(flag bool, typ character.Type) {
	token := &lx.token
	token.typ = TkCharType
	token.setPropCType(typ)
	token.setPropNot(flag)
}

func (lx *Lexer) fetchTokenInCCForP() {
	c2 := lx.peek() // !!! migrate to peekIs
	snx := lx.syntax
	if c2 == '{' && snx.IsOp2(syntax.Op2EscPBraceCharProperty) {
		lx.inc()
		token := &lx.token
		token.typ = TkCharProperty
		token.setPropNot(lx.c == 'P')

		if snx.IsOp2(syntax.Op2EscPBraceCircumflexNot) {
			c2 = lx.fetchTo()
			if c2 == '^' {
				token.setPropNot(!token.getPropNot())
			} else {
				lx.unfetch()
			}
		}
	} else {
		lx.syntaxCharWarn("invalid Unicode Property \\<%n>", rune(lx.c))
	}
}

func (lx *Lexer) fetchTokenInCCForX() {
	if !lx.left() {
		return
	}
	last := lx.p

	snx := lx.syntax
	token := lx.token
	if lx.peekIs('{') && snx.IsOp(syntax.OpEscXBraceHex8) {
		lx.inc()
		num := lx.scanUnsignedHexadecimalNumber(0, 8)
		if num < 0 {
			panic(newSyntaxException(err.CCTooBigWideCharValue))
		}

		enc := lx.enc
		if lx.left() {
			c2 := lx.peek()
			if enc.IsXDigit(c2) {
				panic(newSyntaxException(err.CCTooLongWideCharValue))
			}
		}

		if lx.p > last+enc.Length(lx.bytes, last, lx.stop) && lx.left() && lx.peekIs('}') {
			lx.inc()
			token.typ = TkCodePoint
			token.base = 16
			token.setCode(num)
		} else {
			/* can't read nothing or invalid format */
			lx.p = last
		}
	} else if snx.IsOp(syntax.OpEscXHex2) {
		num := lx.scanUnsignedHexadecimalNumber(0, 2)
		if num < 0 {
			panic(newSyntaxException(err.TooBigNumber))
		}
		if lx.p == last { /* can't read nothing. */
			num = 0 /* but, it's not error */
		}
		token.typ = TkRawByte
		token.base = 16
		token.setC(num)
	}
}

func (lx *Lexer) fetchTokenInCCForU() {
	if !lx.left() {
		return
	}

	last := lx.p
	if lx.syntax.IsOp2(syntax.Op2EscUHex4) {
		num := lx.scanUnsignedHexadecimalNumber(4, 4)
		if num < -1 {
			panic(newSyntaxException(err.TooShortDigits))
		}
		if num < 0 {
			panic(newSyntaxException(err.TooBigNumber))
		}
		if lx.p == last { /* can't read nothing. */
			num = 0 /* but, it's not error */
		}
		token := lx.token
		token.typ = TkCodePoint
		token.base = 16
		token.setCode(num)
	}
}

func (lx *Lexer) fetchTokenInCCForDigit() {
	if lx.syntax.IsOp(syntax.OpEscOctal3) {
		lx.unfetch()
		last := lx.p
		num := lx.scanUnsignedOctalNumber(3)
		if num < 0 || num > 0xff {
			panic(newSyntaxException(err.TooBigNumber))
		}
		if lx.p == last { /* can't read nothing. */
			num = 0 /* but, it's not error */
		}
		token := lx.token
		token.typ = TkRawByte
		token.base = 8
		token.setC(num)
	}
}

func (lx *Lexer) fetchTokenInCCForPosixBracket() {
	snx := lx.syntax
	if snx.IsOp(syntax.OpPosixBracket) && lx.peekIs(':') {
		lx.token.backP = lx.p /* point at '[' is readed */
		lx.inc()
		if lx.strExistCheckWithEsc(send, len(send), ']') {
			lx.token.typ = TkPosixBracketOpen
		} else {
			lx.unfetch()
			// remove duplication, goto cc_in_cc;
			if snx.IsOp2(syntax.Op2CClassSetOp) {
				lx.token.typ = TkCcCcOpen
			} else {
				lx.env.CCEscWarn("[")
			}
		}
	} else { // cc_in_cc:
		if snx.IsOp2(syntax.Op2CClassSetOp) {
			lx.token.typ = TkCcCcOpen
		} else {
			lx.env.CCEscWarn("[")
		}
	}
}

func (lx *Lexer) fetchTokenInCCForAnd() {
	if lx.syntax.IsOp2(syntax.Op2CClassSetOp) && lx.left() && lx.peekIs('&') {
		lx.inc()
		lx.token.typ = TkCcAnd
	}
}

func (lx *Lexer) fetchTokenInCC() TokenType {
	if !lx.left() {
		lx.token.typ = TkEOT
		return TkEOT
	}

	lx.fetch()
	c := lx.c
	token := lx.token
	token.typ = TkChar
	token.base = 0
	token.setC(c)
	token.escaped = false

	if c == ']' {
		token.typ = TkCcClose
	} else if c == '-' {
		token.typ = TkCcRange
	} else if c == lx.syntax.MetaCharTable.Esc {
		snx := lx.syntax
		if !snx.IsBehavior(syntax.BackslashEscapeInCC) {
			return token.typ
		}
		if !lx.left() {
			panic(newSyntaxException(err.EndPatternAtEscape))
		}
		lx.fetch()
		c = lx.c
		token.escaped = true
		token.setC(c)

		switch c {
		case 'w':
			lx.fetchTokenInCCForCharType(false, character.Word)
		case 'W':
			lx.fetchTokenInCCForCharType(true, character.Word)
		case 'd':
			lx.fetchTokenInCCForCharType(false, character.Digit)
		case 'D':
			lx.fetchTokenInCCForCharType(true, character.Digit)
		case 's':
			lx.fetchTokenInCCForCharType(false, character.Space)
		case 'S':
			lx.fetchTokenInCCForCharType(true, character.Space)
		case 'h':
			if snx.IsOp2(syntax.Op2EscHXDigit) {
				lx.fetchTokenInCCForCharType(false, character.XDigit)
			}
		case 'H':
			if snx.IsOp2(syntax.Op2EscHXDigit) {
				lx.fetchTokenInCCForCharType(true, character.XDigit)
			}
		case 'p', 'P':
			lx.fetchTokenInCCForP()
		case 'x':
			lx.fetchTokenInCCForX()
		case 'u':
			lx.fetchTokenInCCForU()
		case '0', '1', '2', '3', '4', '5', '6', '7':
			lx.fetchTokenInCCForDigit()
		default:
			lx.unfetch()
			lx.fetchEscapedValue()
			if token.getC() != c {
				token.setCode(c)
				token.typ = TkCodePoint
			}
		}
	} else if c == '[' {
		lx.fetchTokenInCCForPosixBracket()
	} else if c == '&' {
		lx.fetchTokenInCCForAnd()
	}
	return token.typ
}

func (lx *Lexer) backrefRelToAbs(relNo int) int {
	return lx.env.NumMem() + 1 + relNo
}

func (lx *Lexer) fetchTokenForRepeat(lower, upper int) {
	token := &lx.token
	token.typ = TkOpRepeat
	token.setRepeatLower(lower)
	token.setRepeatUpper(upper)
	lx.greedyCheck()
}

func (lx *Lexer) fetchTokenForOpenBrace() {
	switch lx.fetchRangeQuantifier() {
	case 0:
		lx.greedyCheck()
	case 2:
		if lx.syntax.IsBehavior(syntax.FixedIntervalIsGreedyOnly) {
			lx.possessiveCheck()
		} else {
			lx.greedyCheck()
		}
	default: /* 1 : normal char */
	}
}

func (lx *Lexer) fetchTokenForAnchor(subType anchor.Type) {
	lx.token.typ = TkAnchor
	lx.token.setAnchorSubtype(subType)
}

func (lx *Lexer) fetchTokenForXBrace() {
	if !lx.left() {
		return
	}

	last := lx.p
	snx := lx.syntax
	if lx.peekIs('{') && snx.IsOp(syntax.OpEscXBraceHex8) {
		lx.inc()
		num := lx.scanUnsignedHexadecimalNumber(0, 8)
		if num < 0 {
			panic(newSyntaxException(err.CCTooBigWideCharValue))
		}
		enc := lx.enc
		if lx.left() {
			if enc.IsXDigit(lx.peek()) {
				panic(newSyntaxException(err.CCTooLongWideCharValue))
			}
		}

		if lx.p > last+enc.Length(lx.bytes, last, lx.stop) && lx.left() && lx.peekIs('}') {
			lx.inc()
			lx.token.typ = TkCodePoint
			lx.token.setCode(num)
		} else {
			/* can't read nothing or invalid format */
			lx.p = last
		}
	} else if snx.IsOp(syntax.OpEscXHex2) {
		num := lx.scanUnsignedHexadecimalNumber(0, 2)
		if num < 0 {
			panic(newSyntaxException(err.TooBigNumber))
		}
		if lx.p == last { /* can't read nothing. */
			num = 0 /* but, it's not error */
		}
		lx.token.typ = TkRawByte
		lx.token.base = 16
		lx.token.setC(num)
	}
}

func (lx *Lexer) fetchTokenForUHex() {
	if !lx.left() {
		return
	}

	last := lx.p
	if lx.syntax.IsOp2(syntax.Op2EscUHex4) {
		num := lx.scanUnsignedHexadecimalNumber(4, 4)
		if num < -1 {
			panic(newSyntaxException(err.TooShortDigits))
		}
		if num < 0 {
			panic(newSyntaxException(err.TooBigNumber))
		}
		if lx.p == last { /* can't read nothing. */
			num = 0 /* but, it's not error */
		}
		lx.token.typ = TkCodePoint
		lx.token.base = 16
		lx.token.setCode(num)
	}
}

func (lx *Lexer) fetchTokenForDigit() {
	lx.unfetch()
	last := lx.p
	num := lx.scanUnsignedNumber()
	if num < 0 || num > config.MaxBackrefNum { // goto skip_backref
	} else {
		snx := lx.syntax
		env := lx.env
		if snx.IsOp(syntax.OpDecimalBackref) && (num <= env.NumMem() || num <= 9) { /* This spec. from GNU regex */
			if snx.IsBehavior(syntax.StrictCheckBackref) {
				if num > env.NumMem() || env.MemNodes() == nil || env.MemNodes()[num] == nil {
					panic(newSyntaxException(err.InvalidBackref))
				}
			}
			lx.token.typ = TkBackRef
			lx.token.setBackrefNum(1)
			lx.token.setBackrefRef1(num)
			lx.token.setBackrefByName(false)
			//noinspection GoBoolExpressions
			if config.UseBackrefWithLevel {
				lx.token.setBackrefExistLevel(false)
			}
			return
		}
	}

	if lx.c == '8' || lx.c == '9' { /* normal char */ // skip_backref:
		lx.p = last
		lx.inc()
		return
	}
	lx.p = last

	lx.fetchTokenForZero() /* fall through */
}

func (lx *Lexer) fetchTokenForZero() {
	if lx.syntax.IsOp(syntax.OpEscOctal3) {
		last := lx.p
		n := 3
		if lx.c == '0' {
			n = 2
		}
		num := lx.scanUnsignedOctalNumber(n)
		if num < 0 || num > 0xff {
			panic(newSyntaxException(err.TooBigNumber))
		}
		if lx.p == last { /* can't read nothing. */
			num = 0 /* but, it's not error */
		}
		lx.token.typ = TkRawByte
		lx.token.base = 8
		lx.token.setC(num)
	} else if lx.c != '0' {
		lx.inc()
	}
}

func (lx *Lexer) fetchTokenForNamedBackref() {
	//noinspection GoBoolExpressions
	if config.UseNamedGroup {
		if lx.syntax.IsOp2(syntax.Op2EscKNamedBackref) && lx.left() {
			lx.fetch()
			if lx.c == '<' || lx.c == '\'' {
				lx.fetchNamedBackrefToken()
			} else {
				lx.unfetch()
				lx.syntaxWarn("invalid back reference")
			}
		}
	}
}

func (lx *Lexer) fetchTokenForSubexpCall() {
	snx := lx.syntax
	//noinspection GoBoolExpressions
	if config.UseNamedGroup {
		if snx.IsOp2(syntax.Op2EscGBraceBackref) && lx.left() {
			lx.fetch()
			if lx.c == '{' {
				lx.fetchNamedBackrefToken()
			} else {
				lx.unfetch()
			}
		}
	}
	//noinspection GoBoolExpressions
	if config.UseSubExpCall {
		if snx.IsOp2(syntax.Op2EscGSubexpCall) && lx.left() {
			lx.fetch()
			if lx.c == '<' || lx.c == '\'' {
				gNum := -1
				rel := false
				cnext := lx.peek()
				nameEnd := 0
				if cnext == '0' {
					lx.inc()
					if lx.peekIs(lx.nameEndCodePoint(lx.c)) { /* \g<0>, \g'0' */
						lx.inc()
						nameEnd = lx.p
						gNum = 0
					}
				} else if cnext == '+' {
					lx.inc()
					rel = true
				}
				prev := lx.p
				if gNum < 0 {
					gNum = lx.fetchName(lx.c, true)
					nameEnd = lx.value
				}
				lx.token.typ = TkCall
				lx.token.setCallNameP(prev)
				lx.token.setCallNameEnd(nameEnd)
				lx.token.setCallGNum(gNum)
				lx.token.setCallRel(rel)
			} else {
				lx.syntaxWarn("invalid subexp call")
				lx.unfetch()
			}
		}
	}
}

func (lx *Lexer) fetchNamedBackrefToken() {
	last := lx.p
	snx := lx.syntax
	env := lx.env

	var backNum int
	//noinspection GoBoolExpressions
	if config.UseBackrefWithLevel {
		var level int
		var existLevel bool
		existLevel, backNum, level = lx.fetchNameWithLevel(lx.c)
		lx.token.setBackrefExistLevel(existLevel)
		lx.token.setBackrefLevel(level)
	} else {
		backNum = lx.fetchName(lx.c, true)
	} // USE_BACKREF_AT_LEVEL
	nameEnd := lx.value // set by fetchNameWithLevel/fetchName

	if backNum != 0 {
		if backNum < 0 {
			backNum = lx.backrefRelToAbs(backNum)
			if backNum <= 0 {
				panic(newSyntaxException(err.InvalidBackref))
			}
		}

		if snx.IsBehavior(syntax.StrictCheckBackref) && (backNum > env.NumMem() || env.MemNodes() == nil) {
			panic(newSyntaxException(err.InvalidBackref))
		}
		lx.token.typ = TkBackRef
		lx.token.setBackrefByName(false)
		lx.token.setBackrefNum(1)
		lx.token.setBackrefRef1(backNum)
	} else {
		e := lx.regex.nameToGroupNumbers(lx.bytes, last, nameEnd)
		if e == nil {
			panic(lx.newValueException(err.UndefinedNameReference, last, nameEnd))
		}
		backRefs := e.BackRefs()
		backNum = len(backRefs)

		env := lx.env
		if snx.IsBehavior(syntax.StrictCheckBackref) {
			memNodes := env.MemNodes()
			numMem := env.NumMem()
			if backNum == 1 {
				backRef1 := backRefs[0]
				if backRef1 > numMem ||
					memNodes == nil ||
					memNodes[backRef1] == nil {
					panic(newSyntaxException(err.InvalidBackref))
				}
			} else {
				for _, br := range backRefs {
					if br > numMem ||
						memNodes == nil ||
						memNodes[br] == nil {
						panic(newSyntaxException(err.InvalidBackref))
					}
				}
			}
		}

		lx.token.typ = TkBackRef
		lx.token.setBackrefByName(true)

		if backNum == 1 {
			lx.token.setBackrefNum(1)
			lx.token.setBackrefRef1(backRefs[0])
		} else {
			lx.token.setBackrefNum(backNum)
			lx.token.setBackrefRefs(backRefs)
		}
	}
}

func (lx *Lexer) fetchTokenForCharProperty() {
	snx := lx.syntax
	if lx.peekIs('{') && snx.IsOp2(syntax.Op2EscPBraceCharProperty) {
		lx.inc()
		lx.token.typ = TkCharProperty
		lx.token.setPropNot(lx.c == 'P')

		if snx.IsOp2(syntax.Op2EscPBraceCircumflexNot) {
			lx.fetch()
			if lx.c == '^' {
				lx.token.setPropNot(!lx.token.getPropNot())
			} else {
				lx.unfetch()
			}
		}
	} else {
		lx.syntaxCharWarn("invalid Unicode Property \\<%n>", rune(lx.c))
	}
}

func (lx *Lexer) fetchTokenForMetaChars() {
	metaCharTable := lx.syntax.MetaCharTable
	c := lx.c
	if c == metaCharTable.AnyChar {
		lx.token.typ = TkAnyChar
	} else if c == metaCharTable.AnyTime {
		lx.fetchTokenForRepeat(0, ast.QuantifierRepeatInfinite)
	} else if c == metaCharTable.ZeroOrOneTime {
		lx.fetchTokenForRepeat(0, 1)
	} else if c == metaCharTable.OneOrMoreTime {
		lx.fetchTokenForRepeat(1, ast.QuantifierRepeatInfinite)
	} else if c == metaCharTable.AnyCharAnyTime {
		lx.token.typ = TkAnycharAnytime
	}
}

func (lx *Lexer) fetchToken() {
	snx := lx.syntax
	env := lx.env
	option := env.Option()
	metaCharTable := snx.MetaCharTable

	src := lx.p
	// mark(); // out
start:
	for {
		if !lx.left() {
			lx.token.typ = TkEOT
			return
		}

		lx.token.typ = TkString
		lx.token.base = 0
		lx.token.backP = lx.p

		lx.fetch()

		if lx.c == metaCharTable.Esc && !snx.IsOp2(syntax.Op2IneffectiveEscape) { // IS_MC_ESC_CODE(code, syn)
			if !lx.left() {
				panic(newSyntaxException(err.EndPatternAtEscape))
			}

			lx.token.backP = lx.p
			lx.fetch()

			lx.token.setC(lx.c)
			lx.token.escaped = true
			switch lx.c {

			case '*':
				if snx.IsOp(syntax.OpEscAsteriskZeroInf) {
					lx.fetchTokenForRepeat(0, ast.QuantifierRepeatInfinite)
				}
			case '+':
				if snx.IsOp(syntax.OpEscPlusOneInf) {
					lx.fetchTokenForRepeat(1, ast.QuantifierRepeatInfinite)
				}
			case '?':
				if snx.IsOp(syntax.OpEscQMarkZeroOne) {
					lx.fetchTokenForRepeat(0, 1)
				}
			case '{':
				if snx.IsOp(syntax.OpEscBraceInterval) {
					lx.fetchTokenForOpenBrace()
				}
			case '|':
				if snx.IsOp(syntax.OpEscVBarAlt) {
					lx.token.typ = TkAlt
				}
			case '(':
				if snx.IsOp(syntax.OpEscLParenSubexp) {
					lx.token.typ = TkSubexpOpen
				}
			case ')':
				if snx.IsOp(syntax.OpEscLParenSubexp) {
					lx.token.typ = TkSubexpClose
				}
			case 'w':
				if snx.IsOp(syntax.OpEscWWord) {
					lx.fetchTokenInCCForCharType(false, character.Word)
				}
			case 'W':
				if snx.IsOp(syntax.OpEscWWord) {
					lx.fetchTokenInCCForCharType(true, character.Word)
				}
			case 'b':
				if snx.IsOp(syntax.OpEscBWordBound) {
					lx.fetchTokenForAnchor(anchor.WordBound)
					lx.token.setAnchorASCIIRange(option.IsAsciiRange() && !option.IsWordBoundAllRange())
				}
			case 'B':
				if snx.IsOp(syntax.OpEscBWordBound) {
					lx.fetchTokenForAnchor(anchor.NotWordBound)
					lx.token.setAnchorASCIIRange(option.IsAsciiRange() && !option.IsWordBoundAllRange())
				}
			case '<':
				//noinspection GoBoolExpressions
				if config.UseWordBeginEnd && snx.IsOp(syntax.OpEscLtGtWordBeginEnd) {
					lx.fetchTokenForAnchor(anchor.WordBegin)
					lx.token.setAnchorASCIIRange(option.IsAsciiRange())
				}
			case '>':
				//noinspection GoBoolExpressions
				if config.UseWordBeginEnd && snx.IsOp(syntax.OpEscLtGtWordBeginEnd) {
					lx.fetchTokenForAnchor(anchor.WordEnd)
					lx.token.setAnchorASCIIRange(option.IsAsciiRange())
				}
			case 's':
				if snx.IsOp(syntax.OpEscSWhiteSpace) {
					lx.fetchTokenInCCForCharType(false, character.Space)
				}
			case 'S':
				if snx.IsOp(syntax.OpEscSWhiteSpace) {
					lx.fetchTokenInCCForCharType(true, character.Space)
				}
			case 'd':
				if snx.IsOp(syntax.OpEscDDigit) {
					lx.fetchTokenInCCForCharType(false, character.Digit)
				}
			case 'D':
				if snx.IsOp(syntax.OpEscDDigit) {
					lx.fetchTokenInCCForCharType(true, character.Digit)
				}
			case 'h':
				if snx.IsOp2(syntax.Op2EscHXDigit) {
					lx.fetchTokenInCCForCharType(false, character.XDigit)
				}
			case 'H':
				if snx.IsOp2(syntax.Op2EscHXDigit) {
					lx.fetchTokenInCCForCharType(true, character.XDigit)
				}
			case 'A':
				if snx.IsOp(syntax.OpEscAZBufAnchor) {
					lx.fetchTokenForAnchor(anchor.BeginBuf)
				}
			case 'Z':
				if snx.IsOp(syntax.OpEscAZBufAnchor) {
					lx.fetchTokenForAnchor(anchor.SemiEndBuf)
				}
			case 'z':
				if snx.IsOp(syntax.OpEscAZBufAnchor) {
					lx.fetchTokenForAnchor(anchor.EndBuf)
				}
			case 'G':
				if snx.IsOp(syntax.OpEscCapitalGBeginAnchor) {
					lx.fetchTokenForAnchor(anchor.BeginPosition)
				}
			case '`':
				if snx.IsOp2(syntax.Op2EscGnuBufAnchor) {
					lx.fetchTokenForAnchor(anchor.BeginBuf)
				}
			case '\'':
				if snx.IsOp2(syntax.Op2EscGnuBufAnchor) {
					lx.fetchTokenForAnchor(anchor.EndBuf)
				}
			case 'x':
				lx.fetchTokenForXBrace()
			case 'u':
				lx.fetchTokenForUHex()
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				lx.fetchTokenForDigit()
			case '0':
				lx.fetchTokenForZero()
			case 'k':
				lx.fetchTokenForNamedBackref()
			case 'g':
				lx.fetchTokenForSubexpCall()
			case 'Q':
				if snx.IsOp2(syntax.Op2EscCapitalQQuote) {
					lx.token.typ = TkQuoteOpen
				}
			case 'p', 'P':
				lx.fetchTokenForCharProperty()
			case 'R':
				if snx.IsOp2(syntax.Op2EscCapitalRLinebreak) {
					lx.token.typ = TkLineBreak
				}
			case 'X':
				if snx.IsOp2(syntax.Op2EscCapitalXExtendedGraphemeCluster) {
					lx.token.typ = TkExtendedGraphemeCluster
				}
			case 'K':
				if snx.IsOp2(syntax.Op2EscCapitalKKeep) {
					lx.token.typ = TkKeep
				}
			default:
				lx.unfetch()
				lx.fetchEscapedValue()
				if lx.token.getC() != lx.c { /* set_raw: */
					lx.token.typ = TkCodePoint
					lx.token.setCode(lx.c)
				} else { /* string */
					lx.p = lx.token.backP + lx.enc.Length(lx.bytes, lx.token.backP, lx.stop)
				}
				break
			}
		} else {
			lx.token.setC(lx.c)
			lx.token.escaped = false

			//noinspection GoBoolExpressions
			if config.UseVariableMetaChars && (lx.c != metachar.InnefectiveMetaChar && snx.IsOp(syntax.OpVariableMetaCharacters)) {
				lx.fetchTokenForMetaChars()
				break
			}

			switch lx.c {
			case '.':
				if snx.IsOp(syntax.OpDotAnyChar) {
					lx.token.typ = TkAnyChar
				}
			case '*':
				if snx.IsOp(syntax.OpAsteriskZeroInf) {
					lx.fetchTokenForRepeat(0, ast.QuantifierRepeatInfinite)
				}
			case '+':
				if snx.IsOp(syntax.OpPlusOneInf) {
					lx.fetchTokenForRepeat(1, ast.QuantifierRepeatInfinite)
				}
			case '?':
				if snx.IsOp(syntax.OpQMarkZeroOne) {
					lx.fetchTokenForRepeat(0, 1)
				}
			case '{':
				if snx.IsOp(syntax.OpBraceInterval) {
					lx.fetchTokenForOpenBrace()
				}
			case '|':
				if snx.IsOp(syntax.OpVBarAlt) {
					lx.token.typ = TkAlt
				}
			case '(':
				if lx.peekIs('?') && snx.IsOp2(syntax.Op2QMarkGroupEffect) {
					lx.inc()
					if lx.peekIs('#') {
						lx.fetch()
						for {
							if !lx.left() {
								panic(newSyntaxException(err.EndPatternInGroup))
							}
							lx.fetch()
							if lx.c == metaCharTable.Esc {
								if lx.left() {
									lx.fetch()
								}
							} else {
								if lx.c == ')' {
									break
								}
							}
						}
						continue start // goto start
					}
					lx.unfetch()
				}

				if snx.IsOp(syntax.OpLParenSubexp) {
					lx.token.typ = TkSubexpOpen
				}
			case ')':
				if snx.IsOp(syntax.OpLParenSubexp) {
					lx.token.typ = TkSubexpClose
				}
			case '^':
				if snx.IsOp(syntax.OpLineAnchor) {
					var at anchor.Type
					if option.IsSingleLine() {
						at = anchor.BeginBuf
					} else {
						at = anchor.BeginLine
					}
					lx.fetchTokenForAnchor(at)
				}
			case '$':
				if snx.IsOp(syntax.OpLineAnchor) {
					var at anchor.Type
					if option.IsSingleLine() {
						at = anchor.SemiEndBuf
					} else {
						at = anchor.EndLine
					}
					lx.fetchTokenForAnchor(at)
				}
			case '[':
				if snx.IsOp(syntax.OpBracketCC) {
					lx.token.typ = TkCcOpen
				}
			case ']':
				if src > lx.begin { /* /].../ is allowed. */
					lx.env.CloseBracketWithoutEscapeWarning("]")
				}
			case '#':
				if option.IsExtend() {
					for lx.left() {
						lx.fetch()
						if lx.enc.IsNewLine(lx.c) {
							break
						}
					}
					continue start // goto start
				}
			case ' ', '\t', '\n', '\r', '\f':
				if option.IsExtend() {
					continue start
				}
			}
		}
		break
	}
}

func (lx *Lexer) greedyCheck() {
	if lx.left() && lx.peekIs('?') && lx.syntax.IsOp(syntax.OpQMarkNonGreedy) {
		lx.fetch()

		lx.token.setRepeatGreedy(false)
		lx.token.setRepeatPossessive(false)
	} else {
		lx.possessiveCheck()
	}
}

func (lx *Lexer) possessiveCheck() {
	token := &lx.token
	if lx.left() && lx.peekIs('+') &&
		(lx.syntax.IsOp2(syntax.Op2PlusPossessiveRepeat) && token.typ != TkInterval ||
			lx.syntax.IsOp2(syntax.Op2PlusPossessiveInterval) && token.typ == TkInterval) {

		lx.fetch()

		token.setRepeatGreedy(true)
		lx.token.setRepeatPossessive(true)
	} else {
		token.setRepeatGreedy(true)
		token.setRepeatPossessive(false)
	}
}

func (lx *Lexer) fetchCharPropertyToCType() character.Type {
	lx.mark()

	for lx.left() {
		last := lx.p
		lx.fetch()
		c := lx.c
		if c == '}' {
			return lx.enc.PropertyNameToCType(lx.bytes, lx._p, last)
		}
		if c == '(' || c == ')' || c == '{' || c == '|' {
			panic(err.WithArgs(err.CCInvalidCharPropertyName, issue.H{`n`: string(lx.bytes[lx._p:last])}))
		}
	}
	panic(err.NoArgs(err.ParserBug))
}

func (lx *Lexer) syntaxCharWarn(message string, c rune) {
	lx.syntaxWarn(strings.ReplaceAll(message, "<%n>", string(c)))
}

func (lx *Lexer) syntaxWarn(message string) {
	if (lx.env.Warnings() != &goni.WarnNone{}) {
		lx.env.Warnings().Warn(message + ": /" + string(lx.bytes[lx.begin:lx.end]) + "/")
	}
}
