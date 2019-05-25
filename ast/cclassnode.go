package ast

import (
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/character"
	"github.com/lyraproj/goni/goni/coderange"
	"github.com/lyraproj/goni/goni/node"
	"github.com/lyraproj/goni/goni/syntax"
	"github.com/lyraproj/goni/util"
)

const flagNcCClassNot = 1 << 0

type CClassNode struct {
	abstractNode
	flags int
	bs    *goni.BitSet
	mbuf  *coderange.Buffer
}

func (cn *CClassNode) String() string {
	return goni.String(cn)
}

func (cn *CClassNode) AppendTo(w *util.Indenter) {
	w.NewLine()
	w.Append(`flags: `)
	cn.appendFlags(w)
	w.NewLine()
	w.Append(`bs:`)
	cn.bs.AppendTo(w.Indent())
	if cn.mbuf != nil {
		w.NewLine()
		w.Append(`mbuf: `)
		cn.mbuf.AppendTo(w.Indent())
	}
}

func (cn *CClassNode) Name() string {
	return `Character Class`
}

func NewCClassNode() *CClassNode {
	return &CClassNode{abstractNode: abstractNode{nodeType: node.CClass}, bs: goni.NewBitSet()}
}

func (cn *CClassNode) Clear() {
	cn.bs.ClearAll()
	cn.flags = 0
	cn.mbuf = nil
}

func (cn *CClassNode) isEmpty() bool {
	return cn.mbuf == nil && cn.bs.IsEmpty()
}

func (cn *CClassNode) isNot() bool {
	return (cn.flags & flagNcCClassNot) != 0
}

func (cn *CClassNode) clearNot() {
	cn.flags &= ^flagNcCClassNot
}

func (cn *CClassNode) appendFlags(w *util.Indenter) {
	if cn.isNot() {
		w.Append(`NOT`)
	}
}

func (cn *CClassNode) addCodeRangeToBuf(env goni.ScanEnvironment, from, to int, checkDup bool) {
	cn.mbuf = coderange.AddToBuffer(cn.mbuf, env, from, to, checkDup)
}

func (cn *CClassNode) addCodeRange(env goni.ScanEnvironment, from, to int, checkDup bool) {
	cn.mbuf = coderange.Add(cn.mbuf, env, from, to, checkDup)
}

func (cn *CClassNode) addAllMultiByteRange(env goni.ScanEnvironment) {
	cn.mbuf = coderange.AddAllMultiByte(env, cn.mbuf)
}

func (cn *CClassNode) clearNotFlag(env goni.ScanEnvironment) {
	if cn.isNot() {
		cn.bs.InvertAll()
		if !env.Encoding().IsSingleByte() {
			cn.mbuf = coderange.NotBuffer(env, cn.mbuf)
		}
		cn.clearNot()
	}
}

func (cn *CClassNode) isOneChar() int {
	if cn.isNot() {
		return -1
	}
	c := -1
	if cn.mbuf != nil {
		rng := cn.mbuf.Range()
		c = rng[1]
		if rng[0] == 1 && c == rng[2] {
			if c < goni.SingleByteSize && cn.bs.At(c) {
				c = -1
			}
		} else {
			return -1
		}
	}

	for i := 0; i < goni.BitSetSize; i++ {
		b1 := cn.bs.RoomAt(i)
		if b1 != 0 {
			if (b1&(b1-1)) == 0 && c == -1 {
				c = goni.BitsInRoom*i + goni.BitCount(b1-1)
			} else {
				return -1
			}
		}
	}
	return c
}

func (cn *CClassNode) and(other *CClassNode, env goni.ScanEnvironment) {
	not1 := cn.isNot()
	bsr1 := cn.bs
	buf1 := cn.mbuf
	not2 := other.isNot()
	bsr2 := other.bs
	buf2 := other.mbuf

	if not1 {
		bs1 := goni.NewBitSet()
		bsr1.InvertTo(bs1)
		bsr1 = bs1
	}

	if not2 {
		bs2 := goni.NewBitSet()
		bsr2.InvertTo(bs2)
		bsr2 = bs2
	}

	bsr1.And(bsr2)

	if bsr1 != cn.bs {
		cn.bs.Copy(bsr1)
	}

	if not1 {
		cn.bs.InvertAll()
	}

	if !env.Encoding().IsSingleByte() {
		var pbuf *coderange.Buffer
		if not1 && not2 {
			pbuf = coderange.OrBuffer(env, buf1, false, buf2, false)
		} else {
			pbuf = coderange.AndBuffer(buf1, not1, buf2, not2, env)

			if not1 {
				pbuf = coderange.NotBuffer(env, pbuf)
			}
		}
		cn.mbuf = pbuf
	}
}

func (cn *CClassNode) or(other *CClassNode, env goni.ScanEnvironment) {
	not1 := cn.isNot()
	bsr1 := cn.bs
	buf1 := cn.mbuf
	not2 := other.isNot()
	bsr2 := other.bs
	buf2 := other.mbuf

	if not1 {
		bs1 := goni.NewBitSet()
		bsr1.InvertTo(bs1)
		bsr1 = bs1
	}

	if not2 {
		bs2 := goni.NewBitSet()
		bsr2.InvertTo(bs2)
		bsr2 = bs2
	}

	bsr1.Or(bsr2)

	if bsr1 != cn.bs {
		cn.bs.Copy(bsr1)
	}

	if not1 {
		cn.bs.InvertAll()
	}

	if !env.Encoding().IsSingleByte() {
		var pbuf *coderange.Buffer
		if not1 && not2 {
			pbuf = coderange.AndBuffer(buf1, false, buf2, false, env)
		} else {
			pbuf = coderange.OrBuffer(env, buf1, not1, buf2, not2)
			if not1 {
				pbuf = coderange.NotBuffer(env, pbuf)
			}
		}
		cn.mbuf = pbuf
	}
}

func (cn *CClassNode) addCTypeByRange(ctype character.Type, not bool, env goni.ScanEnvironment, sbOut int, mbr []int) {
	n := mbr[0]
	bs := cn.bs

	if !not {
		i := 0
		for ; i < n; i++ {
			for j := crFrom(mbr, i); j <= crTo(mbr, i); j++ {
				if j >= sbOut {
					if j > crFrom(mbr, i) {
						cn.addCodeRangeToBuf(env, j, crTo(mbr, i), true)
						i++
					}
					// !goto sb_end!, remove duplication!
					for ; i < n; i++ {
						cn.addCodeRangeToBuf(env, crFrom(mbr, i), crTo(mbr, i), true)
					}
					return
				}
				bs.CheckedSet(env, j)
			}
		}
		// !sb_end:!
		for ; i < n; i++ {
			cn.addCodeRangeToBuf(env, crFrom(mbr, i), crTo(mbr, i), true)
		}

	} else {
		prev := 0

		for i := 0; i < n; i++ {
			for j := prev; j < crFrom(mbr, i); j++ {
				if j >= sbOut {
					// !goto sb_end2!, remove duplication
					prev = sbOut
					for i = 0; i < n; i++ {
						if prev < crFrom(mbr, i) {
							cn.addCodeRangeToBuf(env, prev, crFrom(mbr, i)-1, true)
						}
						prev = crTo(mbr, i) + 1
					}
					if prev < 0x7fffffff /*!!!*/ {
						cn.addCodeRangeToBuf(env, prev, 0x7fffffff, true)
					}
					return
				}
				bs.CheckedSet(env, j)
			}
			prev = crTo(mbr, i) + 1
		}

		for j := prev; j < sbOut; j++ {
			bs.CheckedSet(env, j)
		}

		// !sb_end2:!
		prev = sbOut
		for i := 0; i < n; i++ {
			if prev < crFrom(mbr, i) {
				cn.addCodeRangeToBuf(env, prev, crFrom(mbr, i)-1, true)
			}
			prev = crTo(mbr, i) + 1
		}
		if prev < 0x7fffffff /*!!!*/ {
			cn.addCodeRangeToBuf(env, prev, 0x7fffffff, true)
		}
	}
}

func crFrom(rng []int, i int) int {
	return rng[(i*2)+1]
}

func crTo(rng []int, i int) int {
	return rng[(i*2)+2]
}

func (cn *CClassNode) addCType(ctype character.Type, not, asciiRange bool, env goni.ScanEnvironment, sbOut *int) {
	enc := env.Encoding()
	ranges := enc.CTypeCodeRange(ctype, sbOut)
	if ranges != nil {
		if asciiRange {
			ccWork := NewCClassNode()
			ccWork.addCTypeByRange(ctype, not, env, *sbOut, ranges)
			if not {
				ccWork.addCodeRangeToBuf(env, 0x80, coderange.LastCodePoint, false)
			} else {
				ccAscii := NewCClassNode()
				if enc.MinLength() > 1 {
					ccAscii.addCodeRangeToBuf(env, 0x00, 0x7F, true)
				} else {
					ccAscii.bs.CheckedSetRange(env, 0x00, 0x7F)
				}
				ccWork.and(ccAscii, env)
			}
			cn.or(ccWork, env)
		} else {
			cn.addCTypeByRange(ctype, not, env, *sbOut, ranges)
		}
		return
	}

	maxCode := goni.SingleByteSize
	if asciiRange {
		maxCode = 0x80
	}
	switch ctype {
	case character.Alpha,
		character.Blank,
		character.Cntrl,
		character.Digit,
		character.Lower,
		character.Punct,
		character.Space,
		character.Upper,
		character.Xdigit,
		character.Ascii,
		character.Alnum:
		if not {
			for c := 0; c < goni.SingleByteSize; c++ {
				if !enc.IsCodeCType(c, ctype) {
					cn.bs.CheckedSet(env, c)
				}
			}
			cn.addAllMultiByteRange(env)
		} else {
			for c := 0; c < goni.SingleByteSize; c++ {
				if enc.IsCodeCType(c, ctype) {
					cn.bs.CheckedSet(env, c)
				}
			}
		}

	case character.Graph, character.Print:
		if not {
			for c := 0; c < goni.SingleByteSize; c++ {
				if !enc.IsCodeCType(c, ctype) || c >= maxCode {
					cn.bs.CheckedSet(env, c)
				}
			}
			if asciiRange {
				cn.addAllMultiByteRange(env)
			}
		} else {
			for c := 0; c < maxCode; c++ {
				if enc.IsCodeCType(c, ctype) {
					cn.bs.CheckedSet(env, c)
				}
			}
			if !asciiRange {
				cn.addAllMultiByteRange(env)
			}
		}

	case character.Word:
		if !not {
			for c := 0; c < maxCode; c++ {
				if enc.IsSbWord(c) {
					cn.bs.CheckedSet(env, c)
				}
			}
			if !asciiRange {
				cn.addAllMultiByteRange(env)
			}
		} else {
			for c := 0; c < goni.SingleByteSize; c++ {
				if enc.CodeToMbcLength(c) > 0 && /* check invalid code point */
					!(enc.IsWord(c) || c >= maxCode) {
					cn.bs.CheckedSet(env, c)
				}
			}
			if asciiRange {
				cn.addAllMultiByteRange(env)
			}
		}

	default:
		panic(err.NoArgs(err.ParserBug))
	}
}

type CCValType int
type CCState int

const (
	CCValTypeSb        = CCValType(0)
	CCValTypeCodePoint = CCValType(1)
	CCValTypeClass     = CCValType(2)

	CCStateValue    = CCState(0)
	CCStateRange    = CCState(1)
	CCStateComplete = CCState(2)
	CCStateStart    = CCState(3)
)

type CCStateArg struct {
	From      int
	To        int
	FromIsRaw bool
	ToIsRaw   bool
	InType    CCValType
	Type      CCValType
	State     CCState
}

func (cn *CClassNode) nextStateClass(arg *CCStateArg, ascCC *CClassNode, env goni.ScanEnvironment) {
	if arg.State == CCStateRange {
		panic(err.NoArgs(err.CharClassValueAtEndOfRange))
	}

	if arg.State == CCStateValue && arg.Type != CCValTypeClass {
		if arg.Type == CCValTypeSb {
			cn.bs.CheckedSet(env, arg.From)
			if ascCC != nil {
				ascCC.bs.Set(arg.From)
			}
		} else if arg.Type == CCValTypeCodePoint {
			cn.addCodeRange(env, arg.From, arg.From, true)
			if ascCC != nil {
				ascCC.addCodeRange(env, arg.From, arg.From, false)
			}
		}
	}
	arg.State = CCStateValue
	arg.Type = CCValTypeClass
}

func (cn *CClassNode) nextStateValue(arg *CCStateArg, ascCC *CClassNode, env goni.ScanEnvironment) {
	bs := cn.bs
	switch arg.State {
	case CCStateValue:
		if arg.Type == CCValTypeSb {
			bs.CheckedSet(env, arg.From)
			if ascCC != nil {
				ascCC.bs.Set(arg.From)
			}
		} else if arg.Type == CCValTypeCodePoint {
			cn.addCodeRange(env, arg.From, arg.From, true)
			if ascCC != nil {
				ascCC.addCodeRange(env, arg.From, arg.From, false)
			}
		}

	case CCStateRange:
		if arg.InType == arg.Type {
			if arg.InType == CCValTypeSb {
				if arg.From > 0xff || arg.To > 0xff {
					panic(err.NoArgs(err.CCInvalidCodePointValue))
				}

				if arg.From > arg.To {
					if env.Syntax().IsBehavior(syntax.AllowEmptyRangeInCC) {
						// goto ccs_range_end
						arg.State = CCStateComplete
						break
					} else {
						panic(err.NoArgs(err.EmptyRangeInCharClass))
					}
				}
				bs.CheckedSetRange(env, arg.From, arg.To)
				if ascCC != nil {
					ascCC.bs.SetRange(arg.From, arg.To)
				}
			} else {
				cn.addCodeRange(env, arg.From, arg.To, true)
				if ascCC != nil {
					ascCC.addCodeRange(env, arg.From, arg.To, false)
				}
			}
		} else {
			if arg.From > arg.To {
				if env.Syntax().IsBehavior(syntax.AllowEmptyRangeInCC) {
					// goto ccs_range_end
					arg.State = CCStateComplete
					break
				} else {
					panic(err.NoArgs(err.EmptyRangeInCharClass))
				}
			}
			to := 0xff
			if arg.To < 0xff {
				to = arg.To
			}
			bs.CheckedSetRange(env, arg.From, to)
			cn.addCodeRange(env, arg.From, arg.To, true)
			if ascCC != nil {
				ascCC.bs.SetRange(arg.From, to)
				ascCC.addCodeRange(env, arg.From, arg.To, false)
			}
		}
		// ccs_range_end:
		arg.State = CCStateComplete

	case CCStateComplete, CCStateStart:
		arg.State = CCStateValue
	} // switch

	arg.FromIsRaw = arg.ToIsRaw
	arg.From = arg.To
	arg.Type = arg.InType
}

func (cn *CClassNode) isCodeInCCLength(encLength, code int) bool {
	var found bool
	if encLength > 1 || code >= goni.SingleByteSize {
		if cn.mbuf == nil {
			found = false
		} else {
			found = coderange.IsInCodeRange(cn.mbuf.Range(), code)
		}
	} else {
		found = cn.bs.At(code)
	}

	if cn.isNot() {
		return !found
	}
	return found
}

func (cn *CClassNode) isCodeInCC(enc goni.Encoding, code int) bool {
	var ln int
	if enc.MinLength() > 1 {
		ln = 2
	} else {
		ln = enc.CodeToMbcLength(code)
	}
	return cn.isCodeInCCLength(ln, code)
}
