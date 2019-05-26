package internal

import (
	"github.com/lyraproj/goni/ast"
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/syntax"
	"github.com/lyraproj/issue/issue"
)

type Lexer struct {
	scannerSupport
	regex  *Regex
	env    goni.ScanEnvironment
	syntax *goni.Syntax
	token  *Token
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

	t := lx.token
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

// value implicit (rnameEnd)

type Ptr *int

func (lx *Lexer) fetchNameWithLevel(startCode int, rbackNum, rlevel Ptr) bool {
	src := lx.p
	existLevel := false
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
			*rlevel = level * flag
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
		*rbackNum = backNum * sign
	}
	lx.value = nameEnd
	return existLevel
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

func (lx *Lexer) fetchTokenInCCFor_charType(flag bool, typ int) {
	token := lx.token
	token.typ = TkCharType
	token.setPropCType(typ)
	token.setPropNot(flag)
}

var z = `

    private void fetchTokenInCCFor_p() {
        int c2 = peek(); // !!! migrate to peekIs
        if (c2 == '{' && syntax.op2EscPBraceCharProperty()) {
            inc();
            token.type = TokenType.CHAR_PROPERTY;
            token.setPropNot(c == 'P');

            if (syntax.op2EscPBraceCircumflexNot()) {
                c2 = fetchTo();
                if (c2 == '^') {
                    token.setPropNot(!token.getPropNot());
                } else {
                    unfetch();
                }
            }
        } else {
            syntaxWarn("invalid Unicode Property \\<%n>", (char)c);
        }
    }

    private void fetchTokenInCCFor_x() {
        if (!left()) return;
        int last = p;

        if (peekIs('{') && syntax.opEscXBraceHex8()) {
            inc();
            int num = scanUnsignedHexadecimalNumber(0, 8);
            if (num < 0) newValueException(ERR_TOO_BIG_WIDE_CHAR_VALUE);
            if (left()) {
                int c2 = peek();
                if (enc.isXDigit(c2)) newValueException(ERR_TOO_LONG_WIDE_CHAR_VALUE);
            }

            if (p > last + enc.length(bytes, last, stop) && left() && peekIs('}')) {
                inc();
                token.type = TokenType.CODE_POINT;
                token.base = 16;
                token.setCode(num);
            } else {
                /* can't read nothing or invalid format */
                p = last;
            }
        } else if (syntax.opEscXHex2()) {
            int num = scanUnsignedHexadecimalNumber(0, 2);
            if (num < 0) newValueException(TOO_BIG_NUMBER);
            if (p == last) { /* can't read nothing. */
                num = 0; /* but, it's not error */
            }
            token.type = TokenType.RAW_BYTE;
            token.base = 16;
            token.setC(num);
        }
    }

    private void fetchTokenInCCFor_u() {
        if (!left()) return;
        int last = p;

        if (syntax.op2EscUHex4()) {
            int num = scanUnsignedHexadecimalNumber(4, 4);
            if (num < -1) newValueException(TOO_SHORT_DIGITS);
            if (num < 0) newValueException(TOO_BIG_NUMBER);
            if (p == last) {  /* can't read nothing. */
                num = 0; /* but, it's not error */
            }
            token.type = TokenType.CODE_POINT;
            token.base = 16;
            token.setCode(num);
        }
    }

    private void fetchTokenInCCFor_digit() {
        if (syntax.opEscOctal3()) {
            unfetch();
            int last = p;
            int num = scanUnsignedOctalNumber(3);
            if (num < 0 || num > 0xff) newValueException(TOO_BIG_NUMBER);
            if (p == last) {  /* can't read nothing. */
                num = 0; /* but, it's not error */
            }
            token.type = TokenType.RAW_BYTE;
            token.base = 8;
            token.setC(num);
        }
    }

    private void fetchTokenInCCFor_posixBracket() {
        if (syntax.opPosixBracket() && peekIs(':')) {
            token.backP = p; /* point at '[' is readed */
            inc();
            if (strExistCheckWithEsc(send, send.length, ']')) {
                token.type = TokenType.POSIX_BRACKET_OPEN;
            } else {
                unfetch();
                // remove duplication, goto cc_in_cc;
                if (syntax.op2CClassSetOp()) {
                    token.type = TokenType.CC_CC_OPEN;
                } else {
                    env.ccEscWarn("[");
                }
            }
        } else { // cc_in_cc:
            if (syntax.op2CClassSetOp()) {
                token.type = TokenType.CC_CC_OPEN;
            } else {
                env.ccEscWarn("[");
            }
        }
    }

    private void fetchTokenInCCFor_and() {
        if (syntax.op2CClassSetOp() && left() && peekIs('&')) {
            inc();
            token.type = TokenType.CC_AND;
        }
    }

    protected final TokenType fetchTokenInCC() {
        if (!left()) {
            token.type = TokenType.EOT;
            return token.type;
        }

        fetch();
        token.type = TokenType.CHAR;
        token.base = 0;
        token.setC(c);
        token.escaped = false;

        if (c == ']') {
            token.type = TokenType.CC_CLOSE;
        } else if (c == '-') {
            token.type = TokenType.CC_RANGE;
        } else if (c == syntax.metaCharTable.esc) {
            if (!syntax.backSlashEscapeInCC()) return token.type;
            if (!left()) newSyntaxException(END_PATTERN_AT_ESCAPE);
            fetch();
            token.escaped = true;
            token.setC(c);

            switch (c) {
            case 'w':
                fetchTokenInCCFor_charType(false, CharacterType.WORD);
                break;
            case 'W':
                fetchTokenInCCFor_charType(true, CharacterType.WORD);
                break;
            case 'd':
                fetchTokenInCCFor_charType(false, CharacterType.DIGIT);
                break;
            case 'D':
                fetchTokenInCCFor_charType(true, CharacterType.DIGIT);
                break;
            case 's':
                fetchTokenInCCFor_charType(false, CharacterType.SPACE);
                break;
            case 'S':
                fetchTokenInCCFor_charType(true, CharacterType.SPACE);
                break;
            case 'h':
                if (syntax.op2EscHXDigit()) fetchTokenInCCFor_charType(false, CharacterType.XDIGIT);
                break;
            case 'H':
                if (syntax.op2EscHXDigit()) fetchTokenInCCFor_charType(true, CharacterType.XDIGIT);
                break;
            case 'p':
            case 'P':
                fetchTokenInCCFor_p();
                break;
            case 'x':
                fetchTokenInCCFor_x();
                break;
            case 'u':
                fetchTokenInCCFor_u();
                break;
            case '0':
            case '1':
            case '2':
            case '3':
            case '4':
            case '5':
            case '6':
            case '7':
                fetchTokenInCCFor_digit();
                break;

            default:
                unfetch();
                fetchEscapedValue();
                if (token.getC() != c) {
                    token.setCode(c);
                    token.type = TokenType.CODE_POINT;
                }
                break;
            } // switch

        } else if (c == '[') {
            fetchTokenInCCFor_posixBracket();
        } else if (c == '&') {
            fetchTokenInCCFor_and();
        }
        return token.type;
    }

    protected final int backrefRelToAbs(int relNo) {
        return env.numMem + 1 + relNo;
    }

    private void fetchTokenFor_repeat(int lower, int upper) {
        token.type = TokenType.OP_REPEAT;
        token.setRepeatLower(lower);
        token.setRepeatUpper(upper);
        greedyCheck();
    }

    private void fetchTokenFor_openBrace() {
        switch (fetchRangeQuantifier()) {
        case 0:
            greedyCheck();
            break;
        case 2:
            if (syntax.fixedIntervalIsGreedyOnly()) {
                possessiveCheck();
            } else {
                greedyCheck();
            }
            break;
        default: /* 1 : normal char */
        } // inner switch
    }

    private void fetchTokenFor_anchor(int subType) {
        token.type = TokenType.ANCHOR;
        token.setAnchorSubtype(subType);
    }

    private void fetchTokenFor_xBrace() {
        if (!left()) return;

        int last = p;
        if (peekIs('{') && syntax.opEscXBraceHex8()) {
            inc();
            int num = scanUnsignedHexadecimalNumber(0, 8);
            if (num < 0) newValueException(ERR_TOO_BIG_WIDE_CHAR_VALUE);
            if (left()) {
                if (enc.isXDigit(peek())) newValueException(ERR_TOO_LONG_WIDE_CHAR_VALUE);
            }

            if (p > last + enc.length(bytes, last, stop) && left() && peekIs('}')) {
                inc();
                token.type = TokenType.CODE_POINT;
                token.setCode(num);
            } else {
                /* can't read nothing or invalid format */
                p = last;
            }
        } else if (syntax.opEscXHex2()) {
            int num = scanUnsignedHexadecimalNumber(0, 2);
            if (num < 0) newValueException(TOO_BIG_NUMBER);
            if (p == last) { /* can't read nothing. */
                num = 0; /* but, it's not error */
            }
            token.type = TokenType.RAW_BYTE;
            token.base = 16;
            token.setC(num);
        }
    }

    private void fetchTokenFor_uHex() {
        if (!left()) return;
        int last = p;

        if (syntax.op2EscUHex4()) {
            int num = scanUnsignedHexadecimalNumber(4, 4);
            if (num < -1) newValueException(TOO_SHORT_DIGITS);
            if (num < 0) newValueException(TOO_BIG_NUMBER);
            if (p == last) { /* can't read nothing. */
                num = 0; /* but, it's not error */
            }
            token.type = TokenType.CODE_POINT;
            token.base = 16;
            token.setCode(num);
        }
    }

    private void fetchTokenFor_digit() {
        unfetch();
        int last = p;
        int num = scanUnsignedNumber();
        if (num < 0 || num > Config.MAX_BACKREF_NUM) { // goto skip_backref
        } else if (syntax.opDecimalBackref() && (num <= env.numMem || num <= 9)) { /* This spec. from GNU regex */
            if (syntax.strictCheckBackref()) {
                if (num > env.numMem || env.memNodes == null || env.memNodes[num] == null) newValueException(INVALID_BACKREF);
            }
            token.type = TokenType.BACKREF;
            token.setBackrefNum(1);
            token.setBackrefRef1(num);
            token.setBackrefByName(false);
            if (Config.USE_BACKREF_WITH_LEVEL) token.setBackrefExistLevel(false);
            return;
        }

        if (c == '8' || c == '9') { /* normal char */ // skip_backref:
            p = last;
            inc();
            return;
        }
        p = last;

        fetchTokenFor_zero(); /* fall through */
    }

    private void fetchTokenFor_zero() {
        if (syntax.opEscOctal3()) {
            int last = p;
            int num = scanUnsignedOctalNumber(c == '0' ? 2 : 3);
            if (num < 0 || num > 0xff) newValueException(TOO_BIG_NUMBER);
            if (p == last) { /* can't read nothing. */
                num = 0; /* but, it's not error */
            }
            token.type = TokenType.RAW_BYTE;
            token.base = 8;
            token.setC(num);
        } else if (c != '0') {
            inc();
        }
    }

    private void fetchTokenFor_NamedBackref() {
        if (Config.USE_NAMED_GROUP) {
            if (syntax.op2EscKNamedBackref() && left()) {
                fetch();
                if (c =='<' || c == '\'') {
                    fetchNamedBackrefToken();
                } else {
                    unfetch();
                    syntaxWarn("invalid back reference");
                }
            }
        }
    }

    private void fetchTokenFor_subexpCall() {
        if (Config.USE_NAMED_GROUP) {
            if (syntax.op2EscGBraceBackref() && left()) {
                fetch();
                if (c == '{') {
                    fetchNamedBackrefToken();
                } else {
                    unfetch();
                }
            }
        }
        if (Config.USE_SUBEXP_CALL) {
            if (syntax.op2EscGSubexpCall() && left()) {
                fetch();
                if (c == '<' || c == '\'') {
                    int gNum = -1;
                    boolean rel = false;
                    int cnext = peek();
                    int nameEnd = 0;
                    if (cnext == '0') {
                        inc();
                        if (peekIs(nameEndCodePoint(c))) { /* \g<0>, \g'0' */
                            inc();
                            nameEnd = p;
                            gNum = 0;
                        }
                    } else if (cnext == '+') {
                        inc();
                        rel = true;
                    }
                    int prev = p;
                    if (gNum < 0) {
                        gNum = fetchName(c, true);
                        nameEnd = value;
                    }
                    token.type = TokenType.CALL;
                    token.setCallNameP(prev);
                    token.setCallNameEnd(nameEnd);
                    token.setCallGNum(gNum);
                    token.setCallRel(rel);
                } else {
                    syntaxWarn("invalid subexp call");
                    unfetch();
                }
            }
        }
    }

    protected void fetchNamedBackrefToken() {
        int last = p;
        int backNum;
        if (Config.USE_BACKREF_WITH_LEVEL) {
            Ptr rbackNum = new Ptr();
            Ptr rlevel = new Ptr();
            token.setBackrefExistLevel(fetchNameWithLevel(c, rbackNum, rlevel));
            token.setBackrefLevel(rlevel.p);
            backNum = rbackNum.p;
        } else {
            backNum = fetchName(c, true);
        } // USE_BACKREF_AT_LEVEL
        int nameEnd = value; // set by fetchNameWithLevel/fetchName

        if (backNum != 0) {
            if (backNum < 0) {
                backNum = backrefRelToAbs(backNum);
                if (backNum <= 0) newValueException(INVALID_BACKREF);
            }

            if (syntax.strictCheckBackref() && (backNum > env.numMem || env.memNodes == null)) {
                newValueException(INVALID_BACKREF);
            }
            token.type = TokenType.BACKREF;
            token.setBackrefByName(false);
            token.setBackrefNum(1);
            token.setBackrefRef1(backNum);
        } else {
            NameEntry e = regex.nameToGroupNumbers(bytes, last, nameEnd);
            if (e == null) newValueException(UNDEFINED_NAME_REFERENCE, last, nameEnd);

            if (syntax.strictCheckBackref()) {
                if (e.backNum == 1) {
                    if (e.backRef1 > env.numMem ||
                        env.memNodes == null ||
                        env.memNodes[e.backRef1] == null) newValueException(INVALID_BACKREF);
                } else {
                    for (int i=0; i<e.backNum; i++) {
                        if (e.backRefs[i] > env.numMem ||
                            env.memNodes == null ||
                            env.memNodes[e.backRefs[i]] == null) newValueException(INVALID_BACKREF);
                    }
                }
            }

            token.type = TokenType.BACKREF;
            token.setBackrefByName(true);

            if (e.backNum == 1) {
                token.setBackrefNum(1);
                token.setBackrefRef1(e.backRef1);
            } else {
                token.setBackrefNum(e.backNum);
                token.setBackrefRefs(e.backRefs);
            }
        }
    }

    private void fetchTokenFor_charProperty() {
        if (peekIs('{') && syntax.op2EscPBraceCharProperty()) {
            inc();
            token.type = TokenType.CHAR_PROPERTY;
            token.setPropNot(c == 'P');

            if (syntax.op2EscPBraceCircumflexNot()) {
                fetch();
                if (c == '^') {
                    token.setPropNot(!token.getPropNot());
                } else {
                    unfetch();
                }
            }
        } else {
            syntaxWarn("invalid Unicode Property \\<%n>", (char)c);
        }
    }

    private void fetchTokenFor_metaChars() {
        if (c == syntax.metaCharTable.anyChar) {
            token.type = TokenType.ANYCHAR;
        } else if (c == syntax.metaCharTable.anyTime) {
            fetchTokenFor_repeat(0, QuantifierNode.REPEAT_INFINITE);
        }  else if (c == syntax.metaCharTable.zeroOrOneTime) {
            fetchTokenFor_repeat(0, 1);
        } else if (c == syntax.metaCharTable.oneOrMoreTime) {
            fetchTokenFor_repeat(1, QuantifierNode.REPEAT_INFINITE);
        } else if (c == syntax.metaCharTable.anyCharAnyTime) {
            token.type = TokenType.ANYCHAR_ANYTIME;
            // goto out
        }
    }

    protected final void fetchToken() {
        int src = p;
        // mark(); // out
        start:
        while(true) {
            if (!left()) {
                token.type = TokenType.EOT;
                return;
            }

            token.type = TokenType.STRING;
            token.base = 0;
            token.backP = p;

            fetch();

            if (c == syntax.metaCharTable.esc && !syntax.op2IneffectiveEscape()) { // IS_MC_ESC_CODE(code, syn)
                if (!left()) newSyntaxException(END_PATTERN_AT_ESCAPE);

                token.backP = p;
                fetch();

                token.setC(c);
                token.escaped = true;
                switch(c) {

                case '*':
                    if (syntax.opEscAsteriskZeroInf()) fetchTokenFor_repeat(0, QuantifierNode.REPEAT_INFINITE);
                    break;
                case '+':
                    if (syntax.opEscPlusOneInf()) fetchTokenFor_repeat(1, QuantifierNode.REPEAT_INFINITE);
                    break;
                case '?':
                    if (syntax.opEscQMarkZeroOne()) fetchTokenFor_repeat(0, 1);
                    break;
                case '{':
                    if (syntax.opEscBraceInterval()) fetchTokenFor_openBrace();
                    break;
                case '|':
                    if (syntax.opEscVBarAlt()) token.type = TokenType.ALT;
                    break;
                case '(':
                    if (syntax.opEscLParenSubexp()) token.type = TokenType.SUBEXP_OPEN;
                    break;
                case ')':
                    if (syntax.opEscLParenSubexp()) token.type = TokenType.SUBEXP_CLOSE;
                    break;
                case 'w':
                    if (syntax.opEscWWord()) fetchTokenInCCFor_charType(false, CharacterType.WORD);
                    break;
                case 'W':
                    if (syntax.opEscWWord()) fetchTokenInCCFor_charType(true, CharacterType.WORD);
                    break;
                case 'b':
                    if (syntax.opEscBWordBound()) {
                        fetchTokenFor_anchor(AnchorType.WORD_BOUND);
                        token.setAnchorASCIIRange(isAsciiRange(env.option) && !isWordBoundAllRange(env.option));
                    }
                    break;
                case 'B':
                    if (syntax.opEscBWordBound()) {
                        fetchTokenFor_anchor(AnchorType.NOT_WORD_BOUND);
                        token.setAnchorASCIIRange(isAsciiRange(env.option) && !isWordBoundAllRange(env.option));
                    }
                    break;
                case '<':
                    if (Config.USE_WORD_BEGIN_END && syntax.opEscLtGtWordBeginEnd()) {
                        fetchTokenFor_anchor(AnchorType.WORD_BEGIN);
                        token.setAnchorASCIIRange(isAsciiRange(env.option));
                    }
                    break;
                case '>':
                    if (Config.USE_WORD_BEGIN_END && syntax.opEscLtGtWordBeginEnd()) {
                        fetchTokenFor_anchor(AnchorType.WORD_END);
                        token.setAnchorASCIIRange(isAsciiRange(env.option));
                    }
                    break;
                case 's':
                    if (syntax.opEscSWhiteSpace()) fetchTokenInCCFor_charType(false, CharacterType.SPACE);
                    break;
                case 'S':
                    if (syntax.opEscSWhiteSpace()) fetchTokenInCCFor_charType(true, CharacterType.SPACE);
                    break;
                case 'd':
                    if (syntax.opEscDDigit()) fetchTokenInCCFor_charType(false, CharacterType.DIGIT);
                    break;
                case 'D':
                    if (syntax.opEscDDigit()) fetchTokenInCCFor_charType(true, CharacterType.DIGIT);
                    break;
                case 'h':
                    if (syntax.op2EscHXDigit()) fetchTokenInCCFor_charType(false, CharacterType.XDIGIT);
                    break;
                case 'H':
                    if (syntax.op2EscHXDigit()) fetchTokenInCCFor_charType(true, CharacterType.XDIGIT);
                    break;
                case 'A':
                    if (syntax.opEscAZBufAnchor()) fetchTokenFor_anchor(AnchorType.BEGIN_BUF);
                    break;
                case 'Z':
                    if (syntax.opEscAZBufAnchor()) fetchTokenFor_anchor(AnchorType.SEMI_END_BUF);
                    break;
                case 'z':
                    if (syntax.opEscAZBufAnchor()) fetchTokenFor_anchor(AnchorType.END_BUF);
                    break;
                case 'G':
                    if (syntax.opEscCapitalGBeginAnchor()) fetchTokenFor_anchor(AnchorType.BEGIN_POSITION);
                    break;`
var z2 = `':
                    if (syntax.op2EscGnuBufAnchor()) fetchTokenFor_anchor(AnchorType.BEGIN_BUF);
                    break;
                case '\'':
                    if (syntax.op2EscGnuBufAnchor()) fetchTokenFor_anchor(AnchorType.END_BUF);
                    break;
                case 'x':
                    fetchTokenFor_xBrace();
                    break;
                case 'u':
                    fetchTokenFor_uHex();
                    break;
                case '1':
                case '2':
                case '3':
                case '4':
                case '5':
                case '6':
                case '7':
                case '8':
                case '9':
                    fetchTokenFor_digit();
                    break;
                case '0':
                    fetchTokenFor_zero();
                    break;
                case 'k':
                    fetchTokenFor_NamedBackref();
                    break;
                case 'g':
                    fetchTokenFor_subexpCall();
                    break;
                case 'Q':
                    if (syntax.op2EscCapitalQQuote()) token.type = TokenType.QUOTE_OPEN;
                    break;
                case 'p':
                case 'P':
                    fetchTokenFor_charProperty();
                    break;
                case 'R':
                    if (syntax.op2EscCapitalRLinebreak()) token.type = TokenType.LINEBREAK;
                    break;
                case 'X':
                    if (syntax.op2EscCapitalXExtendedGraphemeCluster()) token.type = TokenType.EXTENDED_GRAPHEME_CLUSTER;
                    break;
                case 'K':
                    if (syntax.op2EscCapitalKKeep()) token.type = TokenType.KEEP;
                    break;
                default:
                    unfetch();
                    fetchEscapedValue();
                    if (token.getC() != c) { /* set_raw: */
                        token.type = TokenType.CODE_POINT;
                        token.setCode(c);
                    } else { /* string */
                        p = token.backP + enc.length(bytes, token.backP, stop);
                    }
                    break;
                } // switch (c)
            } else {
                token.setC(c);
                token.escaped = false;

                if (Config.USE_VARIABLE_META_CHARS && (c != MetaChar.INEFFECTIVE_META_CHAR && syntax.opVariableMetaCharacters())) {
                    fetchTokenFor_metaChars();
                    break;
                }

                {
                    switch(c) {
                    case '.':
                        if (syntax.opDotAnyChar()) token.type = TokenType.ANYCHAR;
                        break;
                    case '*':
                        if (syntax.opAsteriskZeroInf()) fetchTokenFor_repeat(0, QuantifierNode.REPEAT_INFINITE);
                        break;
                    case '+':
                        if (syntax.opPlusOneInf()) fetchTokenFor_repeat(1, QuantifierNode.REPEAT_INFINITE);
                        break;
                    case '?':
                        if (syntax.opQMarkZeroOne()) fetchTokenFor_repeat(0, 1);
                        break;
                    case '{':
                        if (syntax.opBraceInterval()) fetchTokenFor_openBrace();
                        break;
                    case '|':
                        if (syntax.opVBarAlt()) token.type = TokenType.ALT;
                        break;

                    case '(':
                        if (peekIs('?') && syntax.op2QMarkGroupEffect()) {
                            inc();
                            if (peekIs('#')) {
                                fetch();
                                while (true) {
                                    if (!left()) newSyntaxException(END_PATTERN_IN_GROUP);
                                    fetch();
                                    if (c == syntax.metaCharTable.esc) {
                                        if (left()) fetch();
                                    } else {
                                        if (c == ')') break;
                                    }
                                }
                                continue start; // goto start
                            }
                            unfetch();
                        }

                        if (syntax.opLParenSubexp()) token.type = TokenType.SUBEXP_OPEN;
                        break;
                    case ')':
                        if (syntax.opLParenSubexp()) token.type = TokenType.SUBEXP_CLOSE;
                        break;
                    case '^':
                        if (syntax.opLineAnchor()) fetchTokenFor_anchor(isSingleline(env.option) ? AnchorType.BEGIN_BUF : AnchorType.BEGIN_LINE);
                        break;
                    case '$':
                        if (syntax.opLineAnchor()) fetchTokenFor_anchor(isSingleline(env.option) ? AnchorType.SEMI_END_BUF : AnchorType.END_LINE);
                        break;
                    case '[':
                        if (syntax.opBracketCC()) token.type = TokenType.CC_OPEN;
                        break;
                    case ']':
                        if (src > getBegin()) { /* /].../ is allowed. */
                            env.closeBracketWithoutEscapeWarn("]");
                        }
                        break;
                    case '#':
                        if (Option.isExtend(env.option)) {
                            while (left()) {
                                fetch();
                                if (enc.isNewLine(c)) break;
                            }
                            continue start; // goto start
                        }
                        break;

                    case ' ':
                    case '\t':
                    case '\n':
                    case '\r':
                    case '\f':
                        if (Option.isExtend(env.option)) continue start; // goto start
                        break;

                    default: // string
                        break;

                    } // switch
                }
            }

            break;
        } // while
    }

    private void greedyCheck() {
        if (left() && peekIs('?') && syntax.opQMarkNonGreedy()) {

            fetch();

            token.setRepeatGreedy(false);
            token.setRepeatPossessive(false);
        } else {
            possessiveCheck();
        }
    }

    private void possessiveCheck() {
        if (left() && peekIs('+') &&
            (syntax.op2PlusPossessiveRepeat() && token.type != TokenType.INTERVAL ||
             syntax.op2PlusPossessiveInterval() && token.type == TokenType.INTERVAL)) {

            fetch();

            token.setRepeatGreedy(true);
            token.setRepeatPossessive(true);
        } else {
            token.setRepeatGreedy(true);
            token.setRepeatPossessive(false);
        }
    }

    protected final int fetchCharPropertyToCType() {
        mark();

        while (left()) {
            int last = p;
            fetch();
            if (c == '}') {
                return enc.propertyNameToCType(bytes, _p, last);
            } else if (c == '(' || c == ')' || c == '{' || c == '|') {
                throw new CharacterPropertyException(EncodingError.ERR_INVALID_CHAR_PROPERTY_NAME, bytes, _p, last);
            }
        }
        newInternalException(PARSER_BUG);
        return 0; // not reached
    }

    protected final void syntaxWarn(String message, char c) {
        syntaxWarn(message.replace("<%n>", Character.toString(c)));
    }

    protected final void syntaxWarn(String message) {
        if (env.warnings != WarnCallback.NONE) {
            env.warnings.warn(message + ": /" + new String(bytes, getBegin(), getEnd()) + "/");
        }
    }
}
`
