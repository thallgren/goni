package internal

import (
	"github.com/lyraproj/goni/ast"
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/anchor"
	"github.com/lyraproj/goni/goni/character"
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

func (px *Parser) init(regex *Regex, syntax *goni.Syntax, bytes []byte, p, end int, warnings WarnCallback) {
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
		opt := px.env.Option()
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
		env.CCEscWarn("]")
		px.token.Type = TkChar /* allow []...] */
	}

	cc = ast.NewCClassNode()
	if env.Option().IsIgnoreCase() {
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
				env.CCEscWarn("[")
				px.p = px.token.backP
				arg.To = px.token.getC()
				arg.ToIsRaw = false
				px.parseCharClassValEntry(cc, ascCc, arg) // goto val_entry
			} else {
				cc.NextStateClass(arg, ascCc, env) // goto next_class
			}

		case TkCharType:
			opt := env.Option()
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
					env.CCEscWarn("-")
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
					env.CCEscWarn("-")
				} /* [--x] or [a&&-x] is warned. */
				px.parseCharClassValEntry(cc, ascCc, arg) // goto val_entry
			} else if arg.State == ast.CCStateRange {
				env.CCEscWarn("-")
				px.parseCharClassSbChar(cc, ascCc, arg) // goto sb_char /* [!--x] is allowed */
			} else { /* CCS_COMPLETE */
				px.fetchTokenInCC()
				fetched = true
				if px.token.Type == TkCcClose { /* allow [a-b-] */
					px.parseCharClassRangeEndVal(cc, ascCc, arg) // goto range_end_val
					break
				}
				if px.token.Type == TkCcAnd {
					env.CCEscWarn("-")
					px.parseCharClassRangeEndVal(cc, ascCc, arg) // goto range_end_val
					break
				}

				if snx.IsBehavior(syntax.AllowDoubleRangeOpInCC) {
					env.CCEscWarn("-")
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
	opt := env.Option()
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
				env.PushPrecReadNotNode(nd)
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
				num := env.AddMemEntry()
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
					prev := env.Option()
					env.SetOption(opt)
					px.fetchToken()
					target := px.parseSubExp(term)
					env.SetOption(prev)
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
		en := ast.NewMemory(env.Option(), false)
		num := env.AddMemEntry()
		en.RegNum = num
		nd = en
	}

	px.fetchToken()
	target := px.parseSubExp(term)

	if nd.Type() == node.Anchor {
		an := nd.(*ast.AnchorNode)
		an.SetTarget(target)
		if snx.IsOp3(syntax.Op3OptionECMAScript) && an.AnchorType() == anchor.PrecReadNot {
			env.PopPrecReadNotNode(an)
		}
	} else {
		en := nd.(*ast.EncloseNode)
		en.SetTarget(target)
		if en.EncloseType() == enclose.Memory {
			if snx.IsOp3(syntax.Op3OptionECMAScript) {
				en.ContainingAnchor = env.CurrentPrecReadNotNode()
			}
			/* Don't move this to previous of parse_subexp() */
			env.SetMemNode(en.RegNum, en)
		} else if en.EncloseType() == enclose.Condition {
			if target.Type() != node.Alt { /* convert (?(cond)yes) to (?(cond)yes|empty) */
				en.SetTarget(ast.NewAlt(target, ast.NewAlt(ast.StringNodeEmpty, nil)))
			}
		}
	}
	px.returnCode = 0
	return nd // ??
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

func (px *Parser) parseEncloseNamedGroup2(b bool) goni.Node {
	return nil // TODO
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
