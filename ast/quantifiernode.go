package ast

import (
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/goni/syntax"
	"github.com/lyraproj/goni/goni/targetinfo"
	"github.com/lyraproj/goni/util"
)

const QuantifierRepeatInfinite = -1

type QuantifierNode struct {
	stateNode
	target goni.Node
	lower  int
	upper  int
	greedy bool

	targetEmptyInfo int
	headExact       goni.Node
	nextHeadExact   goni.Node
	isRefered       bool /* include called node. don't eliminate even if {0} */

	// USE_COMBINATION_EXPLOSION_CHECK
	combExpCheckNum int /* 1,2,3...: check,  0: no check  */
}

func NewQuantifierNode(lower, upper int, byNumber bool) *QuantifierNode {
	qn := &QuantifierNode{stateNode: stateNode{abstractNode: abstractNode{nodeType: node.QTFR}},
		lower: lower, upper: upper, greedy: true, targetEmptyInfo: targetinfo.IsNotEmpty}
	if byNumber {
		qn.SetByNumber()
	}
	return qn
}

func (qn *QuantifierNode) Child() goni.Node {
	return qn.target
}

func (qn *QuantifierNode) String() string {
	return goni.String(qn)
}

func (qn *QuantifierNode) Name() string {
	return `Quantifier`
}

func (qn *QuantifierNode) SetChild(child goni.Node) {
	qn.target = child
}

func (qn *QuantifierNode) SetTarget(target goni.Node) {
	qn.target = target
	target.SetParent(qn)
}

func (qn *QuantifierNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	w.Append(`lower: `)
	w.AppendInt(qn.lower)
	w.Append(`, upper: `)
	w.AppendInt(qn.upper)
	w.Append(`, greedy: `)
	w.AppendBool(qn.greedy)
	w.Append(`, isRefered: `)
	w.AppendBool(qn.isRefered)
	w.Append(`, targetEmptyInfo: `)
	w.AppendInt(qn.targetEmptyInfo)
	w.Append(`, combExpCheckNum: `)
	w.AppendInt(qn.combExpCheckNum)
	if qn.headExact != nil {
		w.NewLine()
		w.Append(`headExact: `)
		qn.headExact.AppendTo(w.Indent())
	}
	if qn.nextHeadExact != nil {
		w.NewLine()
		w.Append(`nextHeadExact: `)
		qn.nextHeadExact.AppendTo(w.Indent())
	}
	if qn.target != nil {
		w.NewLine()
		w.Append(`target: `)
		qn.target.AppendTo(w.Indent())
	}
}

func (qn *QuantifierNode) isAnyCharStar() bool {
	return qn.greedy && isRepeatInfinite(qn.upper) && qn.target.Type() == node.CAny
}

func (qn *QuantifierNode) popularNum() int {
	if qn.greedy {
		if qn.lower == 0 {
			if qn.upper == 1 {
				return 0
			} else if isRepeatInfinite(qn.upper) {
				return 1
			}
		} else if qn.lower == 1 {
			if isRepeatInfinite(qn.upper) {
				return 2
			}
		}
	} else {
		if qn.lower == 0 {
			if qn.upper == 1 {
				return 3
			} else if isRepeatInfinite(qn.upper) {
				return 4
			}
		} else if qn.lower == 1 {
			if isRepeatInfinite(qn.upper) {
				return 5
			}
		}
	}
	return -1
}

type reduceType int

const (
	rtAsis = reduceType(iota) /* as is */
	rtDel                     /* delete parent */
	rtA                       /* to '*'    */
	rtAQ                      /* to '*?'   */
	rtQQ                      /* to '??'   */
	rtPQQ                     /* to '+)??' */
	rtPQPQ                    /* to '+?)?' */
)

var reduceTable = [6][6]reduceType{
	{rtDel, rtA, rtA, rtQQ, rtAQ, rtAsis},      /* '?'  */
	{rtDel, rtDel, rtDel, rtPQQ, rtPQQ, rtDel}, /* '*'  */
	{rtA, rtA, rtDel, rtAsis, rtPQQ, rtDel},    /* '+'  */
	{rtDel, rtAQ, rtAQ, rtDel, rtAQ, rtAQ},     /* '??' */
	{rtDel, rtDel, rtDel, rtDel, rtDel, rtDel}, /* '*?' */
	{rtAsis, rtPQPQ, rtDel, rtAQ, rtAQ, rtDel}} /* '+?' */

func (qn *QuantifierNode) copy(other *QuantifierNode) {
	qn.state = other.state
	qn.SetTarget(other.target)
	other.target = nil
	qn.lower = other.lower
	qn.upper = other.upper
	qn.greedy = other.greedy
	qn.targetEmptyInfo = other.targetEmptyInfo
	qn.headExact = other.headExact
	qn.nextHeadExact = other.nextHeadExact
	qn.isRefered = other.isRefered
	qn.combExpCheckNum = other.combExpCheckNum
}

func (qn *QuantifierNode) reduceNestedQuantifier(other *QuantifierNode) {
	pnum := qn.popularNum()
	cnum := other.popularNum()

	if pnum < 0 || cnum < 0 {
		return
	}

	switch reduceTable[cnum][pnum] {
	case rtDel:
		// no need to set the parent here...
		qn.copy(other)

	case rtA:
		qn.SetTarget(other.target)
		qn.lower = 0
		qn.upper = QuantifierRepeatInfinite
		qn.greedy = true

	case rtAQ:
		qn.SetTarget(other.target)
		qn.lower = 0
		qn.upper = QuantifierRepeatInfinite
		qn.greedy = false

	case rtQQ:
		qn.SetTarget(other.target)
		qn.lower = 0
		qn.upper = 1
		qn.greedy = false

	case rtPQQ:
		qn.SetTarget(other)
		qn.lower = 0
		qn.upper = 1
		qn.greedy = false
		other.lower = 1
		other.upper = QuantifierRepeatInfinite
		other.greedy = true
		return

	case rtPQPQ:
		qn.SetTarget(other)
		qn.lower = 0
		qn.upper = 1
		qn.greedy = true
		other.lower = 1
		other.upper = QuantifierRepeatInfinite
		other.greedy = false
		return

	case rtAsis:
		qn.SetTarget(other)
		return
	}
	other.target = nil // remove target from reduced quantifier
}

var popularQStr = []string{"?", "*", "+", "??", "*?", "+?"}
var reduceQStr = []string{"", "", "*", "*?", "??", "+ and ??", "+? and ?"}

func (qn *QuantifierNode) setQuantifier(tgt goni.Node, group bool, env goni.ScanEnvironment, bytes []byte, p, end int) int {
	if qn.lower == 1 && qn.upper == 1 {
		if env.Syntax().IsOp3(syntax.Op3OptionECMAScript) {
			qn.SetTarget(tgt)
		}
		return 1
	}

	switch tgt.Type() {

	case node.Str:
		if !group {
			sn := tgt.(*StringNode)
			enc := env.Encoding()
			if sn.canBeSplit(enc) {
				n := sn.splitLastChar(enc)
				if n != nil {
					qn.SetTarget(n)
					return 2
				}
			}
		}
		break

	case node.QTFR:
		/* check redundant double repeat. */
		/* verbose warn (?:.?)? etc... but not warn (.?)? etc... */
		qnt := tgt.(*QuantifierNode)
		nestQNum := qn.popularNum()
		targetQNum := qnt.popularNum()

		//noinspection GoBoolExpressions
		if config.UseWarningRedundantNestedRepeatOperator {
			if nestQNum >= 0 && targetQNum >= 0 && env.Syntax().IsBehavior(syntax.WarnRedundantNestedRepeat) {
				switch reduceTable[targetQNum][nestQNum] {
				case rtAsis:
				case rtDel:
					env.Warnings().Warn("regular expression has redundant nested repeat operator " + popularQStr[targetQNum] + " /" + string(bytes[p:end]) + "/")
				default:
					env.Warnings().Warn("nested repeat operator '" + popularQStr[targetQNum] + "' and '" + popularQStr[nestQNum] +
						"' was replaced with '" + reduceQStr[reduceTable[targetQNum][nestQNum]] + "' in regular expression " + "/" + string(bytes[p:end]) + "/")
				}
			}
		}

		if targetQNum >= 0 {
			if nestQNum >= 0 {
				qn.reduceNestedQuantifier(qnt)
				return 0
			} else if targetQNum == 1 || targetQNum == 2 { /* * or + */
				/* (?:a*){n,m}, (?:a+){n,m} => (?:a*){n,n}, (?:a+){n,n} */
				if !isRepeatInfinite(qn.upper) && qn.upper > 1 && qn.greedy {
					if qn.lower == 0 {
						qn.upper = 1
					} else {
						qn.upper = qn.lower
					}
				}
			}
		}
	}

	qn.SetTarget(tgt)
	return 0
}

func isRepeatInfinite(n int) bool {
	return n == QuantifierRepeatInfinite
}
