package internal

import (
	"github.com/lyraproj/goni/ast"
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/option"
	"github.com/lyraproj/goni/goni/syntax"
)

type scanEnvironment struct {
	options        option.Type
	caseFoldFlag   int
	enc            goni.Encoding
	syntax         *goni.Syntax
	captureHistory option.Type
	btMemStart     int
	btMemEnd       int
	backrefedMem   int

	warnings goni.WarnCallback

	numCall       int
	unsetAddrList *ast.UnsetAddrList // USE_SUBEXP_CALL

	numNamed int // USE_NAMED_GROUP

	memNodes []goni.Node

	// USE_COMBINATION_EXPLOSION_CHECK
	numCombExpCheck  int
	combExpMaxRegNum int
	currMaxRegNum    int
	hasRecursion     bool
	warningsFlag     syntax.Behavior

	numPrecReadNotNodes int
	precReadNotNodes    []goni.Node
}

func newScanEnvironment(regex *regex, syntax *goni.Syntax, warnings goni.WarnCallback) *scanEnvironment {
	return &scanEnvironment{
		syntax:       syntax,
		warnings:     warnings,
		options:      regex.options,
		caseFoldFlag: regex.caseFoldFlag,
		enc:          regex.enc}
}

func (se *scanEnvironment) addMemEntry() int {
	numMem := len(se.memNodes) + 1
	if numMem >= config.MaxCaptureGroupNum {
		panic(err.NoArgs(err.TooManyCaptureGroups))
	}
	se.memNodes = append(se.memNodes, nil)
	return numMem
}

func (se *scanEnvironment) setMemNode(i int, node *ast.EncloseNode) {
	se.memNodes[i] = node
}

func (se *scanEnvironment) pushPrecReadNotNode(node goni.Node) {
	se.precReadNotNodes = append(se.precReadNotNodes, node)
}

func (se *scanEnvironment) popPrecReadNotNode(node goni.Node) {
	se.precReadNotNodes = se.precReadNotNodes[:len(se.precReadNotNodes)-1]
}

func (se *scanEnvironment) currentPrecReadNotNode() goni.Node {
	n := len(se.precReadNotNodes) - 1
	if n >= 0 {
		return se.precReadNotNodes[n]
	}
	return nil
}

func (se *scanEnvironment) convertBackslashValue(c int) int {
	snx := se.syntax
	if snx.IsOp(syntax.OpEscControlChars) {
		switch c {
		case 'n':
			return '\n'
		case 't':
			return '\t'
		case 'r':
			return '\r'
		case 'f':
			return '\f'
		case 'a':
			return '\007'
		case 'b':
			return '\010'
		case 'e':
			return '\033'
		case 'v':
			if snx.IsOp2(syntax.Op2EscVVtab) {
				return 11 // '\v'
			}
		default:
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
				se.unknownEscWarn(string(rune(c)))
			}
		}
	}
	return c
}

func (se *scanEnvironment) ccEscWarn(s string) {
	if se.warnings != nil {
		snx := se.syntax
		if snx.IsBehavior(syntax.WarnCCOpNotEscaped) && snx.IsBehavior(syntax.BackslashEscapeInCC) {
			se.warnings.Warn(`character class has '` + s +  `' without escape`);
		}
	}
}

func (se *scanEnvironment) unknownEscWarn(s string) {
	if se.warnings != nil {
		se.warnings.Warn(`Unknown escape \` + s + ` is ignored`)
	}
}

func (se *scanEnvironment) closeBracketWithoutEscapeWarning(s string) {
	if se.warnings != nil && se.syntax.IsBehavior(syntax.WarnCCOpNotEscaped) {
		se.warnings.Warn(`regular expression has '` + s +  `' without escape`);
	}
}

func (se *scanEnvironment) ccDuplicateWarning() {
	if (se.syntax.IsBehavior(syntax.WarnCCDup) && (se.warningsFlag & syntax.WarnCCDup) == 0) {
		se.warningsFlag |= syntax.WarnCCDup;
		// FIXME: jruby/joni#34 points out problem and what it will take to uncomment this (we were getting erroneous versions of this)
		// se.warnings.Warn("character class has duplicated range");
	}
}

func (se *scanEnvironment) Encoding() goni.Encoding {
	return se.enc
}

func (se *scanEnvironment) Option() option.Type {
	return se.options
}

func (se *scanEnvironment) Syntax() *goni.Syntax {
	return se.syntax
}

func (se *scanEnvironment) Warnings() goni.WarnCallback {
	return se.warnings
}

func (se *scanEnvironment) NumMem() int {
	return len(se.memNodes)
}

func (se *scanEnvironment) MemNodes() []goni.Node {
	return se.memNodes
}
