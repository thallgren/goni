package internal

import (
	"github.com/lyraproj/goni/ast"
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/anchor"
	"github.com/lyraproj/goni/goni/character"
	"github.com/lyraproj/goni/goni/coderange"
	"github.com/lyraproj/goni/goni/enclose"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/goni/option"
	"github.com/lyraproj/goni/goni/posix"
	"github.com/lyraproj/goni/goni/syntax"
	"github.com/lyraproj/issue/issue"
)

const (
	posixBracketNameMinLen       = 4
	posixBracketCheckLimitLength = 20
)

var bracketEnd = []byte{':', ']'}

type Parser struct {
	Lexer
	returnCode int
}

func (px *Parser) init(regex *regex, syntax *goni.Syntax, bytes []byte, p, end int, warnings goni.WarnCallback) {
	px.Lexer.init(regex, syntax, bytes, p, end, warnings)
}

func (px *Parser) parsePosixBracket(cc, ascCc *ast.CClassNode) bool {
	px.mark()

	not := false
	if px.peekIs('^') {
		px.inc()
		not = true
	}

	enc := px.enc
	if enc.StrLength(px.bytes, px.p, px.stop) >= posixBracketNameMinLen+3 {
		opt := px.env.options
		asciiRange := opt.IsAsciiRange() && !opt.IsPosixBracketAllRange()

		for i, name := range posix.PBSNamesLower {
			// hash lookup here ?
			if enc.StrNCmp(px.bytes, px.p, px.stop, name, 0, len(name)) == 0 {
				px.p = enc.Step(px.bytes, px.p, px.stop, len(name))
				if enc.StrNCmp(px.bytes, px.p, px.stop, bracketEnd, 0, len(bracketEnd)) != 0 {
					panic(newSyntaxException(err.InvalidPosixBracketType))
				}
				ctype := posix.PBSValues[i]
				cc.AddCType(ctype, not, asciiRange, px.env, &px.value)
				if ascCc != nil {
					if ctype != character.Word && ctype != character.Ascii && !asciiRange {
						ascCc.AddCType(ctype, not, asciiRange, px.env, &px.value)
					}
				}
				px.inc()
				px.inc()
				return false
			}
		}
	}

	px.c = 0
	i := 0
	for px.left() {
		px.c = px.peek()
		if px.c == ':' || px.c == ']' {
			break
		}
		px.inc()
		i++
		if i > posixBracketCheckLimitLength {
			break
		}
	}

	if px.c == ':' && px.left() {
		px.inc()
		if px.left() {
			px.fetch()
			if px.c == ']' {
				panic(newSyntaxException(err.InvalidPosixBracketType))
			}
		}
	}
	px.restore()
	return true /* 1: is not POSIX bracket, but no error. */
}

func (px *Parser) codeExistCheck(code int, ignoreEscaped bool) bool {
	px.mark()

	inEsc := false
	esc := px.syntax.MetaCharTable.Esc
	for px.left() {
		if ignoreEscaped && inEsc {
			inEsc = false
		} else {
			px.fetch()
			if px.c == code {
				px.restore()
				return true
			}
			if px.c == esc {
				inEsc = true
			}
		}
	}

	px.restore()
	return false
}

func (px *Parser) parseCharClass(ascNode **ast.CClassNode) *ast.CClassNode {
	var cc, prevCc, ascCc, ascPrevCc, workCc, ascWorkCc *ast.CClassNode
	arg := &ast.CCStateArg{}
	px.fetchTokenInCC()
	neg := false
	if px.token.Type == TkChar && px.token.getC() == '^' && !px.token.escaped {
		neg = true
		px.fetchTokenInCC()
	}

	snx := px.syntax
	env := px.env
	if px.token.Type == TkCcClose && !snx.IsOp3(syntax.Op3OptionECMAScript) {
		if !px.codeExistCheck(']', true) {
			panic(newSyntaxException(err.EmptyCharClass))
		}
		env.ccEscWarn("]")
		px.token.Type = TkChar /* allow []...] */
	}

	cc = ast.NewCClassNode()
	if env.options.IsIgnoreCase() {
		ascCc = ast.NewCClassNode()
		*ascNode = ascCc
	}

	andStart := false
	arg.State = ast.CCStateStart
	enc := px.enc
	for px.token.Type != TkCcClose {
		fetched := false
		var ln int
		switch px.token.Type {
		case TkChar:
			arg.InType = ast.CCValTypeSb
			if px.token.getCode() >= goni.SingleByteSize {
				ln = enc.CodeToMbcLength(px.token.getC())
				if ln > 1 {
					arg.InType = ast.CCValTypeCodePoint
				}
			}

			arg.To = px.token.getC()
			arg.ToIsRaw = false
			px.parseCharClassValEntry2(cc, ascCc, arg) // goto val_entry2

		case TkRawByte:
			if !enc.IsSingleByte() && px.token.base != 0 { /* tok->base != 0 : octal or hexadec. */
				buf := make([]byte, config.EncMbcCaseFoldMaxlen)
				psave := px.p
				base := px.token.base
				buf[0] = byte(px.token.getC())
				var i int
				for i = 1; i < enc.MaxLength(); i++ {
					px.fetchTokenInCC()
					if px.token.Type != TkRawByte || px.token.base != base {
						fetched = true
						break
					}
					buf[i] = byte(px.token.getC())
				}
				if i < enc.MinLength() {
					panic(newSyntaxException(err.TooShortMultiByteString))
				}

				ln = enc.Length(buf, 0, i)
				if i < ln {
					panic(newSyntaxException(err.TooShortMultiByteString))
				} else if i > ln { /* fetch back */
					px.p = psave
					for i = 1; i < ln; i++ {
						px.fetchTokenInCC()
					}
					fetched = false
				}
				if i == 1 {
					arg.To = int(buf[0] & 0xff)
					arg.InType = ast.CCValTypeSb // goto raw_single
				} else {
					arg.To, _ = enc.MbcToCode(buf, 0, len(buf))
					arg.InType = ast.CCValTypeCodePoint
				}
			} else {
				arg.To = px.token.getC()
				arg.InType = ast.CCValTypeSb // raw_single:
			}
			arg.ToIsRaw = true
			px.parseCharClassValEntry2(cc, ascCc, arg) // goto val_entry2

		case TkCodePoint:
			arg.To = px.token.getCode()
			arg.ToIsRaw = true
			px.parseCharClassValEntry(cc, ascCc, arg) // val_entry:, val_entry2

		case TkPosixBracketOpen:
			if px.parsePosixBracket(cc, ascCc) { /* true: is not POSIX bracket */
				env.ccEscWarn("[")
				px.p = px.token.backP
				arg.To = px.token.getC()
				arg.ToIsRaw = false
				px.parseCharClassValEntry(cc, ascCc, arg) // goto val_entry
			} else {
				cc.NextStateClass(arg, ascCc, env) // goto next_class
			}

		case TkCharType:
			opt := env.options
			cc.AddCType(px.token.getPropCType(), px.token.getPropNot(), opt.IsAsciiRange(), env, &px.value)
			if ascCc != nil {
				if px.token.getPropCType() != character.Word {
					ascCc.AddCType(px.token.getPropCType(), px.token.getPropNot(), opt.IsAsciiRange(), env, &px.value)
				}
			}
			cc.NextStateClass(arg, ascCc, env) // next_class:

		case TkCharProperty:
			ctype := px.fetchCharPropertyToCType()
			cc.AddCType(ctype, px.token.getPropNot(), false, env, &px.value)
			if ascCc != nil {
				if ctype != character.Ascii {
					ascCc.AddCType(ctype, px.token.getPropNot(), false, env, &px.value)
				}
			}
			cc.NextStateClass(arg, ascCc, env) // goto next_class

		case TkCcRange:
			if arg.State == ast.CCStateValue {
				px.fetchTokenInCC()
				fetched = true
				if px.token.Type == TkCcClose { /* allow [x-] */
					px.parseCharClassRangeEndVal(cc, ascCc, arg) // range_end_val:, goto val_entry;
					break
				}
				if px.token.Type == TkCcAnd {
					env.ccEscWarn("-")
					px.parseCharClassRangeEndVal(cc, ascCc, arg) // goto range_end_val
					break
				}
				if arg.Type == ast.CCValTypeClass {
					panic(newSyntaxException(err.UnmatchedRangeSpecifierInCharClass))
				}
				arg.State = ast.CCStateRange
			} else if arg.State == ast.CCStateStart {
				arg.To = px.token.getC() /* [-xa] is allowed */
				arg.ToIsRaw = false
				px.fetchTokenInCC()
				fetched = true
				if px.token.Type == TkCcRange || andStart {
					env.ccEscWarn("-")
				} /* [--x] or [a&&-x] is warned. */
				px.parseCharClassValEntry(cc, ascCc, arg) // goto val_entry
			} else if arg.State == ast.CCStateRange {
				env.ccEscWarn("-")
				px.parseCharClassSbChar(cc, ascCc, arg) // goto sb_char /* [!--x] is allowed */
			} else { /* CCS_COMPLETE */
				px.fetchTokenInCC()
				fetched = true
				if px.token.Type == TkCcClose { /* allow [a-b-] */
					px.parseCharClassRangeEndVal(cc, ascCc, arg) // goto range_end_val
					break
				}
				if px.token.Type == TkCcAnd {
					env.ccEscWarn("-")
					px.parseCharClassRangeEndVal(cc, ascCc, arg) // goto range_end_val
					break
				}

				if snx.IsBehavior(syntax.AllowDoubleRangeOpInCC) {
					env.ccEscWarn("-")
					// parseCharClassSbChar(cc, ascCc, arg); // goto sb_char /* [0-9-a] is allowed as [0-9\-a] */
					px.parseCharClassRangeEndVal(cc, ascCc, arg) // goto range_end_val
					break
				}
				panic(newSyntaxException(err.UnmatchedRangeSpecifierInCharClass))
			}

		case TkCcCcOpen: /* [ */
			var ascPtr *ast.CClassNode
			acc := px.parseCharClass(&ascPtr)
			cc.Or(acc, env)
			if ascPtr != nil && ascCc != nil {
				ascCc.Or(ascPtr, env)
			}

		case TkCcAnd: /* && */
			if arg.State == ast.CCStateValue {
				arg.To = 0
				arg.ToIsRaw = false
				cc.NextStateValue(arg, ascCc, env)
			}
			/* initialize local variables */
			andStart = true
			arg.State = ast.CCStateStart
			if prevCc != nil {
				prevCc.And(cc, env)
				if ascCc != nil && ascPrevCc != nil {
					ascPrevCc.And(ascCc, env)
				}
			} else {
				prevCc = cc
				if workCc == nil {
					workCc = ast.NewCClassNode()
				}
				cc = workCc
				if ascCc != nil {
					ascPrevCc = ascCc
					if ascWorkCc == nil {
						ascWorkCc = ast.NewCClassNode()
					}
					ascCc = ascWorkCc
				}
			}
			cc.Clear()
			if ascCc != nil {
				ascCc.Clear()
			}
			break

		case TkEOT:
			panic(newSyntaxException(err.PrematureEndOfCharClass))

		default:
			panic(newSyntaxException(err.ParserBug))
		}

		if !fetched {
			px.fetchTokenInCC()
		}

	}

	if arg.State == ast.CCStateValue {
		arg.To = 0
		arg.ToIsRaw = false
		cc.NextStateValue(arg, ascCc, env)
	}

	if prevCc != nil {
		prevCc.And(cc, env)
		cc = prevCc
		if ascCc != nil && ascPrevCc != nil {
			ascPrevCc.And(ascCc, env)
			ascCc = ascPrevCc
		}
	}

	if neg {
		cc.SetNot()
		if ascCc != nil {
			ascCc.SetNot()
		}
	} else {
		cc.ClearNot()
		if ascCc != nil {
			ascCc.ClearNot()
		}
	}

	if cc.IsNot() && snx.IsBehavior(syntax.NotNewlineInNegativeCC) {
		if !cc.IsEmpty() { // ???
			newLine := 0x0a
			if enc.IsNewLine(newLine) {
				if enc.CodeToMbcLength(newLine) == 1 {
					cc.BitSet().CheckedSet(env, newLine)
				} else {
					cc.AddCodeRange(env, newLine, newLine, true)
				}
			}
		}
	}

	return cc
}

func (px *Parser) parseCharClassSbChar(cc *ast.CClassNode, ascCc *ast.CClassNode, arg *ast.CCStateArg) {
	arg.InType = ast.CCValTypeSb
	arg.To = px.token.getC()
	arg.ToIsRaw = false
	px.parseCharClassValEntry2(cc, ascCc, arg) // goto val_entry2
}

func (px *Parser) parseCharClassRangeEndVal(cc *ast.CClassNode, ascCc *ast.CClassNode, arg *ast.CCStateArg) {
	arg.To = '-'
	arg.ToIsRaw = false
	px.parseCharClassValEntry(cc, ascCc, arg) // goto val_entry
}

func (px *Parser) parseCharClassValEntry(cc *ast.CClassNode, ascCc *ast.CClassNode, arg *ast.CCStateArg) {
	ln := px.enc.CodeToMbcLength(arg.To)
	arg.InType = ast.CCValTypeCodePoint
	if ln == 1 {
		arg.InType = ast.CCValTypeSb
	}
	px.parseCharClassValEntry2(cc, ascCc, arg) // val_entry2:
}

func (px *Parser) parseCharClassValEntry2(cc *ast.CClassNode, ascCc *ast.CClassNode, arg *ast.CCStateArg) {
	cc.NextStateValue(arg, ascCc, px.env)
}

func (px *Parser) parseEnclose(term TokenType) goni.Node {
	var nd goni.Node

	if !px.left() {
		panic(newSyntaxException(err.EndPatternWithUnmatchedParenthesis))
	}

	env := px.env
	enc := px.enc
	opt := env.options
	snx := px.syntax

	if px.peekIs('?') && snx.IsOp2(syntax.Op2QMarkGroupEffect) {
		px.inc()
		if !px.left() {
			panic(newSyntaxException(err.EndPatternInGroup))
		}

		listCapture := false

		px.fetch()
		switch px.c {
		case ':': /* (?:...) grouping only */
			px.fetchToken() // group:
			nd = px.parseSubExp(term)
			px.returnCode = 1 /* group */
			return nd
		case '=':
			nd = ast.NewAnchorNode(anchor.PrecRead, false)
		case '!': /*         preceding read */
			nd = ast.NewAnchorNode(anchor.PrecReadNot, false)
			if snx.IsOp3(syntax.Op3OptionECMAScript) {
				env.pushPrecReadNotNode(nd)
			}
		case '>': /* (?>...) stop backtrack */
			nd = ast.NewEncloseNode(enclose.StopBacktrack) // node_new_enclose
		case '~': /* (?~...) absent operator */
			if snx.IsOp2(syntax.Op2QMarkTildeAbsent) {
				nd = ast.NewEncloseNode(enclose.Absent)
			} else {
				panic(newSyntaxException(err.UndefinedGroupOption))
			}
		case '\'':
			//noinspection GoBoolExpressions
			if config.UseNamedGroup && snx.IsOp2(syntax.Op2QMarkLtNamedGroup) {
				listCapture = false // goto named_group1
				nd = px.parseEncloseNamedGroup2(listCapture)
			} else {
				panic(newSyntaxException(err.UndefinedGroupOption))
			}
		case '<': /* look behind (?<=...), (?<!...) */
			px.fetch()
			if px.c == '=' {
				nd = ast.NewAnchorNode(anchor.LookBehind, false)
			} else if px.c == '!' {
				nd = ast.NewAnchorNode(anchor.LookBehindNot, false)
			} else {
				//noinspection GoBoolExpressions
				if config.UseNamedGroup {
					if snx.IsOp2(syntax.Op2QMarkLtNamedGroup) {
						px.unfetch()
						px.c = '<'

						listCapture = false                          // named_group1:
						nd = px.parseEncloseNamedGroup2(listCapture) // named_group2:
					} else {
						panic(newSyntaxException(err.UndefinedGroupOption))
					}
				} else { // USE_NAMED_GROUP
					panic(newSyntaxException(err.UndefinedGroupOption))
				} // USE_NAMED_GROUP
			}
		case '@':
			if snx.IsOp2(syntax.Op2AtMarkCaptureHistory) {
				//noinspection GoBoolExpressions
				if config.UseNamedGroup {
					if snx.IsOp2(syntax.Op2QMarkLtNamedGroup) {
						px.fetch()
						if px.c == '<' || px.c == '\'' {
							listCapture = true
							nd = px.parseEncloseNamedGroup2(listCapture) // goto named_group2 /* (?@<name>...) */
						}
						px.unfetch()
					}
				} // USE_NAMED_GROUP
				en := ast.NewMemory(opt, false)
				num := env.addMemEntry()
				if num >= option.BitsNum {
					panic(newSyntaxException(err.GroupNumberOverForCaptureHistory))
				}
				en.RegNum = num
				nd = en
			} else {
				panic(newSyntaxException(err.UndefinedGroupOption))
			}
		case '(': /* conditional expression: (?(cond)yes), (?(cond)yes|no) */
			if snx.IsOp2(syntax.Op2QMarkLParenCondition) {
				num := -1
				name := -1
				px.fetch()
				if enc.IsDigit(px.c) { /* (n) */
					px.unfetch()
					num = px.fetchName('(', true)
					if snx.IsBehavior(syntax.StrictCheckBackref) {
						if num > env.NumMem() || env.MemNodes() == nil || env.MemNodes()[num] == nil {
							panic(newSyntaxException(err.InvalidBackref))
						}
					}
				} else {
					//noinspection GoBoolExpressions
					if config.UseNamedGroup {
						if px.c == '<' || px.c == '\'' { /* (<name>), ('name') */
							name = px.p
							px.fetchNamedBackrefToken()
							px.inc()
							if px.token.getBackrefNum() > 1 {
								num = px.token.getBackrefRefs()[0]
							} else {
								num = px.token.getBackrefRef1()
							}
						}
					} else { // USE_NAMED_GROUP
						panic(newSyntaxException(err.InvalidConditionPattern))
					}
				}
				en := ast.NewEncloseNode(enclose.Condition)
				en.RegNum = num
				if name != -1 {
					en.SetNameRef()
				}
				nd = en
			} else {
				panic(newSyntaxException(err.UndefinedGroupOption))
			}

		case '^': /* loads default options */
			if px.left() && snx.IsOp2(syntax.Op2OptionPerl) {
				/* d-imsx */
				opt = option.OnOff(opt, option.AsciiRange, true)
				opt = option.OnOff(opt, option.IgnoreCase, true)
				opt = option.OnOff(opt, option.SingleLine, false)
				opt = option.OnOff(opt, option.MultiLine, true)
				opt = option.OnOff(opt, option.Extend, true)
				px.fetch()
			} else {
				panic(newSyntaxException(err.UndefinedGroupOption))
			}
			fallthrough

		// case 'p': #ifdef USE_POSIXLINE_OPTION
		case '-', 'i', 'm', 's', 'x', 'a', 'd', 'l', 'u':
			neg := false
			for {
				switch px.c {
				case ':', ')':

				case '-':
					neg = true
				case 'x':
					opt = option.OnOff(opt, option.Extend, neg)
				case 'i':
					opt = option.OnOff(opt, option.IgnoreCase, neg)
				case 's':
					if snx.IsOp2(syntax.Op2OptionPerl) {
						opt = option.OnOff(opt, option.MultiLine, neg)
					} else {
						panic(newSyntaxException(err.UndefinedGroupOption))
					}
				case 'm':
					if snx.IsOp2(syntax.Op2OptionPerl) {
						opt = option.OnOff(opt, option.SingleLine, !neg)
					} else if snx.IsOp2(syntax.Op2OptionRuby) {
						opt = option.OnOff(opt, option.MultiLine, neg)
					} else {
						panic(newSyntaxException(err.UndefinedGroupOption))
					}
				case 'a': /* limits \d, \s, \w and POSIX brackets to ASCII range */
					if (snx.IsOp2(syntax.Op2OptionPerl) || snx.IsOp2(syntax.Op2OptionRuby)) && !neg {
						opt = option.OnOff(opt, option.AsciiRange, false)
						opt = option.OnOff(opt, option.PosixBracketAllRange, true)
						opt = option.OnOff(opt, option.WordBoundAllRange, true)
					} else {
						panic(newSyntaxException(err.UndefinedGroupOption))
					}
				case 'u':
					if (snx.IsOp2(syntax.Op2OptionPerl) || snx.IsOp2(syntax.Op2OptionRuby)) && !neg {
						opt = option.OnOff(opt, option.AsciiRange, true)
						opt = option.OnOff(opt, option.PosixBracketAllRange, true)
						opt = option.OnOff(opt, option.WordBoundAllRange, true)
					} else {
						panic(newSyntaxException(err.UndefinedGroupOption))
					}

				case 'd':
					if snx.IsOp2(syntax.Op2OptionPerl) && !neg {
						opt = option.OnOff(opt, option.AsciiRange, true)
					} else if snx.IsOp2(syntax.Op2OptionRuby) && !neg {
						opt = option.OnOff(opt, option.AsciiRange, false)
						opt = option.OnOff(opt, option.PosixBracketAllRange, false)
						opt = option.OnOff(opt, option.WordBoundAllRange, false)
					} else {
						panic(newSyntaxException(err.UndefinedGroupOption))
					}
				case 'l':
					if snx.IsOp2(syntax.Op2OptionPerl) && !neg {
						opt = option.OnOff(opt, option.AsciiRange, true)
					} else {
						panic(newSyntaxException(err.UndefinedGroupOption))
					}
				default:
					panic(newSyntaxException(err.UndefinedGroupOption))
				}

				if px.c == ')' {
					en := ast.NewOption(opt)
					nd = en
					px.returnCode = 2 /* option only */
					return nd
				}
				if px.c == ':' {
					prev := env.options
					env.options = opt
					px.fetchToken()
					target := px.parseSubExp(term)
					env.options = prev
					en := ast.NewOption(opt)
					en.SetTarget(target)
					nd = en
					px.returnCode = 0
					return nd
				}
				if !px.left() {
					panic(newSyntaxException(err.EndPatternInGroup))
				}
				px.fetch()
			}

		default:
			panic(newSyntaxException(err.UndefinedGroupOption))
		}
	} else {
		if opt.IsDontCaptureGroup() {
			px.fetchToken() // goto group
			nd = px.parseSubExp(term)
			px.returnCode = 1 /* group */
			return nd
		}
		en := ast.NewMemory(env.options, false)
		num := env.addMemEntry()
		en.RegNum = num
		nd = en
	}

	px.fetchToken()
	target := px.parseSubExp(term)

	if nd.Type() == node.Anchor {
		an := nd.(*ast.AnchorNode)
		an.SetTarget(target)
		if snx.IsOp3(syntax.Op3OptionECMAScript) && an.AnchorType() == anchor.PrecReadNot {
			env.popPrecReadNotNode(an)
		}
	} else {
		en := nd.(*ast.EncloseNode)
		en.SetTarget(target)
		if en.EncloseType() == enclose.Memory {
			if snx.IsOp3(syntax.Op3OptionECMAScript) {
				en.ContainingAnchor = env.currentPrecReadNotNode()
			}
			/* Don't move this to previous of parse_subexp() */
			env.setMemNode(en.RegNum, en)
		} else if en.EncloseType() == enclose.Condition {
			if target.Type() != node.Alt { /* convert (?(cond)yes) to (?(cond)yes|empty) */
				en.SetTarget(ast.NewAlt(target, ast.NewAlt(ast.StringNodeEmpty, nil)))
			}
		}
	}
	px.returnCode = 0
	return nd // ??
}

func (px *Parser) parseEncloseNamedGroup2(listCapture bool) goni.Node {
	nm := px.p
	num := px.fetchName(px.c, false)
	nameEnd := px.value
	num = px.env.addMemEntry()
	if listCapture && num >= option.BitsNum {
		panic(newSyntaxException(err.GroupNumberOverForCaptureHistory))
	}

	px.regex.nameAdd(px.bytes, nm, nameEnd, num, px.syntax)
	en := ast.NewMemory(px.env.options, true)
	en.RegNum = num
	if listCapture {
		px.env.captureHistory = option.OnAtSimple(px.env.captureHistory, num)
	}
	px.env.numNamed++
	return en
}

func (px *Parser) findStrPosition(s []int, n, from, to int, nextChar *int) int {
	p := from
	enc := px.enc
	for p < to {
		x, l := enc.MbcToCode(px.bytes, p, to)
		q := p + l
		if x == s[0] {
			i := 1
			for ; i < n && q < to; i++ {
				x, l = enc.MbcToCode(px.bytes, q, to)
				if x != s[i] {
					break
				}
				q += l
			}
			if i >= n {
				if px.bytes[*nextChar] != 0 {
					*nextChar = q // we may need zero term semantics...
				}
				return p
			}
		}
		p = q
	}
	return -1
}

func (px *Parser) parseExp(term TokenType) goni.Node {
	if px.token.Type == term {
		return ast.StringNodeEmpty
	}
	var nd goni.Node
	group := false

	switch px.token.Type {
	case TkAlt, TkEOT:
		return ast.StringNodeEmpty // end_of_token:, node_new_empty

	case TkSubexpOpen:
		nd = px.parseEnclose(TkSubexpClose)
		if px.returnCode == 1 {
			group = true
		} else if px.returnCode == 2 { /* option only */
			env := px.env
			prev := env.options
			en := nd.(*ast.EncloseNode)
			env.options = en.Option
			px.fetchToken()
			target := px.parseSubExp(term)
			env.options = prev
			en.SetTarget(target)
			return nd
		}
	case TkSubexpClose:
		if !px.syntax.IsBehavior(syntax.AllowUnmatchedCloseSubexp) {
			panic(newSyntaxException(err.UnmatchedCloseParenthesis))
		}
		if px.token.escaped {
			return px.parseExpTkRawByte(group) // goto tk_raw_byte
		}
		return px.parseExpTkByte(group) // goto tk_byte
	case TkLineBreak:
		nd = px.parseLineBreak()
	case TkExtendedGraphemeCluster:
		nd = px.parseExtendedGraphemeCluster()
	case TkKeep:
		nd = ast.NewAnchorNode(anchor.Keep, false)
	case TkString:
		return px.parseExpTkByte(group) // tk_byte:
	case TkRawByte:
		return px.parseExpTkRawByte(group) // tk_raw_byte:
	case TkCodePoint:
		return px.parseStringLoop(StringNode.fromCodePoint(token.getCode(), enc), group)
	case TkQuoteOpen:
		nd = px.parseQuoteOpen()
	case TkCharType:
		nd = px.parseCharType(nd)
	case TkCharProperty:
		nd = px.parseCharProperty()
	case TkCcOpen:
		var ascPtr *ast.CClassNode
		cc := px.parseCharClass(&ascPtr)
		code := cc.IsOneChar()
		if code != -1 {
			return px.parseStringLoop(StringNode.fromCodePoint(code, enc), group)
		}

		nd = cc
		if px.env.options.IsIgnoreCase() {
			nd = cClassCaseFold(nd, cc, ascPtr)
		}

	case TkAnyChar:
		nd = ast.NewAnyCharNode()
	case TkAnycharAnytime:
		nd = px.parseAnycharAnytime()
	case TkBackRef:
		nd = px.parseBackref()
	case TkCall:
		//noinspection GoBoolExpressions
		if config.UseSubExpCall {
			nd = px.parseCall()
		}
	case TkAnchor:
		nd = ast.NewAnchorNode(px.token.getAnchorSubtype(), px.token.getAnchorASCIIRange())
	case TkOpRepeat, TkInterval:
		snx := px.syntax
		if snx.IsBehavior(syntax.ContextIndepRepeatOps) {
			if snx.IsBehavior(syntax.ContextInvalidRepeatOps) {
				panic(newSyntaxException(err.TargetOfRepeatOperatorNotSpecified))
			} else {
				nd = ast.StringNodeEmpty // node_new_empty
			}
		} else {
			return px.parseExpTkByte(group) // goto tk_byte
		}
		break

	default:
		panic(newSyntaxException(err.ParserBug))
	}

	//targetp = node;

	px.fetchToken() // re_entry:

	return px.parseExpRepeat(nd, group) // repeat:
}

func (px *Parser) parseLineBreak() goni.Node {
	enc := px.enc
	buflb := make([]byte, 0, config.EncCodeToMbcMaxlen*2)
	buflb = enc.CodeToMbc(0x0D, buflb)
	buflb = enc.CodeToMbc(0x0A, buflb)
	left := ast.NewStringNodeShared(buflb, 0)
	left.SetRaw()
	/* [\x0A-\x0D] or [\x0A-\x0D\x{85}\x{2028}\x{2029}] */
	right := ast.NewCClassNode()
	env := px.env
	if enc.MinLength() > 1 {
		right.AddCodeRange(env, 0x0A, 0x0D, true)
	} else {
		right.Bs.CheckedSetRange(env, 0x0A, 0x0D)
	}

	if enc.IsUnicode() {
		/* UTF-8, UTF-16BE/LE, UTF-32BE/LE */
		right.AddCodeRange(env, 0x85, 0x85, true)
		right.AddCodeRange(env, 0x2028, 0x2029, true)
	}
	/* (?>...) */
	en := ast.NewEncloseNode(enclose.StopBacktrack)
	en.SetTarget(ast.NewAlt(left, ast.NewAlt(right, nil)))
	return en
}

var graphemeClusterBreakExtend = []byte(`Grapheme_Cluster_Break=Extend`)
var graphemeClusterBreakControl = []byte(`Grapheme_Cluster_Break=Control`)
var graphemeClusterBreakPrepend = []byte(`Grapheme_Cluster_Break=Prepend`)
var graphemeClusterBreakL = []byte(`Grapheme_Cluster_Break=L`)
var graphemeClusterBreakV = []byte(`Grapheme_Cluster_Break=V`)
var graphemeClusterBreakLV = []byte(`Grapheme_Cluster_Break=LV`)
var graphemeClusterBreakLVT = []byte(`Grapheme_Cluster_Break=LVT`)
var graphemeClusterBreakT = []byte(`Grapheme_Cluster_Break=T`)
var regionalIndicator = []byte(`Regional_Indicator`)
var extendedPictographic = []byte(`Extended_Pictographic`)
var graphemeClusterBreakSpacingMark = []byte(`Grapheme_Cluster_Break=SpacingMark`)

func (px *Parser) addPropertyToCC(cc *ast.CClassNode, propName []byte, not bool) {
	ctype := px.enc.PropertyNameToCType(propName, 0, len(propName))
	cc.AddCType(ctype, not, false, px.env, &px.value)
}

func (px *Parser) createPropertyNode(nodes []goni.Node, np int, propName []byte) {
	cc := ast.NewCClassNode()
	px.addPropertyToCC(cc, propName, false)
	nodes[np] = cc
}

func (px *Parser) quantifierNode(nodes []goni.Node, np, lower, upper int) {
	qnf := ast.NewQuantifierNode(lower, upper, false)
	qnf.SetTarget(nodes[np])
	nodes[np] = qnf
}

func (px *Parser) quantifierPropertyNode(nodes []goni.Node, np int, propName []byte, repetitions rune) {
	lower := 0
	upper := ast.QuantifierRepeatInfinite

	px.createPropertyNode(nodes, np, propName)
	switch repetitions {
	case '?':
		upper = 1
	case '+':
		lower = 1
	case '*':
		// No op
	case '2':
		lower = 2
		upper = lower
	default:
		panic(newSyntaxException(err.ParserBug))
	}

	px.quantifierNode(nodes, np, lower, upper)
}

func (px *Parser) createNodeFromArray(list bool, nodes []goni.Node, np, nodeArray int) {
	i := 0
	for nodes[nodeArray+i] != nil {
		i++
	}

	var tmp *ast.ListNode
	for i--; i >= 0; i-- {
		n := nodes[nodeArray+i]
		if list {
			tmp = ast.NewList(n, tmp)
		} else {
			tmp = ast.NewAlt(n, tmp)
		}
		nodes[np] = tmp
		nodes[nodeArray+i] = nil
	}
}

func (px *Parser) createNodeFromArray2(nodes []goni.Node, nodeArray int) *ast.ListNode {
	i := 0
	for nodes[nodeArray+i] != nil {
		i++
	}

	var np *ast.ListNode
	for i--; i >= 0; i-- {
		np = ast.NewAlt(nodes[nodeArray+i], np)
		nodes[nodeArray+i] = nil
	}
	return np
}

const nodeCommonSize = 16

func (px *Parser) parseExtendedGraphemeCluster() goni.Node {
	nodes := make([]goni.Node, nodeCommonSize)
	var anyTargetPosition int
	alts := 0

	enc := px.enc
	env := px.env

	strNode := ast.NewStringNodeWithCapacity(config.EncCodeToMbcMaxlen * 2)
	strNode.SetRaw()
	strNode.CatCode(0x0D, enc)
	strNode.CatCode(0x0A, enc)
	nodes[alts] = strNode

	//noinspection GoBoolExpressions
	if config.UseUnicodeProperties && enc.IsUnicode() {
		cc := ast.NewCClassNode()
		nodes[alts+1] = cc
		px.addPropertyToCC(cc, graphemeClusterBreakControl, false)
		if enc.MinLength() > 1 {
			cc.AddCodeRange(env, 0x000A, 0x000A, true)
			cc.AddCodeRange(env, 0x000D, 0x000D, true)
		} else {
			cc.Bs.Set(0x0A)
			cc.Bs.Set(0x0D)
		}

		list := alts + 3
		px.quantifierPropertyNode(nodes, list+0, graphemeClusterBreakPrepend, '*')
		coreAlts := list + 2

		HList := coreAlts + 1
		px.quantifierPropertyNode(nodes, HList+0, graphemeClusterBreakL, '*')

		HAlt2 := HList + 2
		px.quantifierPropertyNode(nodes, HAlt2+0, graphemeClusterBreakV, '+')

		HList2 := HAlt2 + 2
		px.createPropertyNode(nodes, HList2+0, graphemeClusterBreakLV)
		px.quantifierPropertyNode(nodes, HList2+1, graphemeClusterBreakV, '*')
		px.createNodeFromArray(true, nodes, HAlt2+1, HList2)
		px.createPropertyNode(nodes, HAlt2+2, graphemeClusterBreakLVT)
		px.createNodeFromArray(false, nodes, HList+1, HAlt2)
		px.quantifierPropertyNode(nodes, HList+2, graphemeClusterBreakT, '*')
		px.createNodeFromArray(true, nodes, coreAlts+0, HList)

		px.quantifierPropertyNode(nodes, coreAlts+1, graphemeClusterBreakL, '+')
		px.quantifierPropertyNode(nodes, coreAlts+2, graphemeClusterBreakT, '+')
		px.quantifierPropertyNode(nodes, coreAlts+3, regionalIndicator, '2')

		XPList := coreAlts + 5
		px.createPropertyNode(nodes, XPList+0, extendedPictographic)

		ExList := XPList + 2
		px.quantifierPropertyNode(nodes, ExList+0, graphemeClusterBreakExtend, '*')
		strNode = ast.NewStringNodeWithCapacity(config.EncCodeToMbcMaxlen)
		strNode.SetRaw()
		strNode.CatCode(0x200D, enc)
		nodes[ExList+1] = strNode
		px.createPropertyNode(nodes, ExList+2, extendedPictographic)
		px.createNodeFromArray(true, nodes, XPList+1, ExList)

		px.quantifierNode(nodes, XPList+1, 0, ast.QuantifierRepeatInfinite)
		px.createNodeFromArray(true, nodes, coreAlts+4, XPList)

		cc = ast.NewCClassNode()
		nodes[coreAlts+5] = cc
		if enc.MinLength() > 1 {
			px.addPropertyToCC(cc, graphemeClusterBreakControl, false)
			cc.AddCodeRange(env, 0x000A, 0x000A, true)
			cc.AddCodeRange(env, 0x000D, 0x000D, true)
			cc.Mbuf = coderange.NotBuffer(env, cc.Mbuf)
		} else {
			px.addPropertyToCC(cc, graphemeClusterBreakControl, true)
			cc.Bs.Clear(0x0A)
			cc.Bs.Clear(0x0D)
		}
		px.createNodeFromArray(false, nodes, list+1, coreAlts)

		px.createPropertyNode(nodes, list+2, graphemeClusterBreakExtend)
		cc = nodes[list+2].(*ast.CClassNode)
		px.addPropertyToCC(cc, graphemeClusterBreakSpacingMark, false)
		cc.AddCodeRange(px.env, 0x200D, 0x200D, true)
		px.quantifierNode(nodes, list+2, 0, ast.QuantifierRepeatInfinite)
		px.createNodeFromArray(true, nodes, alts+2, list)

		anyTargetPosition = 3
	} else { // enc.isUnicode()
		anyTargetPosition = 1
	}

	any := ast.NewAnyCharNode()
	opt := ast.NewOption(option.OnOff(env.options, option.MultiLine, false))
	opt.SetTarget(any)
	nodes[anyTargetPosition] = opt

	topAlt := px.createNodeFromArray2(nodes, alts)
	encl := ast.NewEncloseNode(enclose.StopBacktrack)
	encl.SetTarget(topAlt)

	//noinspection GoBoolExpressions
	if config.UseUnicodeProperties && enc.IsUnicode() {
		opt = ast.NewOption(option.OnOff(env.options, option.IgnoreCase, true))
		opt.SetTarget(encl)
		return opt
	}
	return encl
}

func (px *Parser) parseExpTkByte(group bool) goni.Node {
	backP := px.token.backP
	nd := ast.NewStringNodeShared(px.bytes[backP:px.p], backP); // tk_byte:
	return px.parseStringLoop(nd, group);
}

func (px *Parser) parseStringLoop(node *ast.StringNode, group bool) goni.Node {
	enc := px.enc
	for {
		px.fetchToken();
		if (px.token.Type == TkString) {
			if (px.token.backP == node.End()) {
				node.SetEnd(px.p) // non escaped character, remain shared, just increase shared range
			} else {
				node.CatBytes(px.bytes[px.token.backP:px.p]) // non continuous string stream, need to COW
			}
		} else if (px.token.Type == TkCodePoint) {
			node.CatCode(px.token.getCode(), enc);
		} else {
			break;
		}
	}
	// targetp = node;
	return px.parseExpRepeat(node, group); // string_end:, goto repeat
}

func (px *Parser) parseExpTkRawByte(group bool) goni.Node {
// tk_raw_byte:
node := ast.NewStringNode();
node.SetRaw();
node.CatByte(byte(px.token.getC()))

len := 1;
enc := px.enc
for {
if (len >= enc.MinLength()) {
if (len == enc.Length(node.bytes, 0, len(node.bytes)) {
fetchToken();
node.clearRaw();
// !goto string_end;!
return parseExpRepeat(node, group);
}
}

fetchToken();
if (token.type != TokenType.RAW_BYTE) {
/* Don't use this, it is wrong for little endian encodings. */
// USE_PAD_TO_SHORT_BYTE_CHAR ...
newValueException(TOO_SHORT_MULTI_BYTE_STRING);
}
node.catByte((byte)token.getC());
len++;
} // while
}


func (px *Parser) parseSubExp(term TokenType) goni.Node {
	nd := px.parseBranch(term)

	if px.token.Type == term {
		return nd
	}
	if px.token.Type == TkAlt {
		top := ast.NewAlt(nd, nil)
		t := top
		for px.token.Type == TkAlt {
			px.fetchToken()
			nd = px.parseBranch(term)

			t.SetTail(ast.NewAlt(nd, nil))
			t = t.Tail
		}

		if px.token.Type != term {
			panic(parseSubExpError(term))
		}
		return top
	}
	panic(parseSubExpError(term))
}

func (px *Parser) parseBranch(tokenType TokenType) goni.Node {
	return nil // TODO
}

func parseSubExpError(term TokenType) issue.Reported {
	if term == TkSubexpClose {
		return newSyntaxException(err.EndPatternWithUnmatchedParenthesis)
	}
	return newSyntaxException(err.ParserBug)
}
