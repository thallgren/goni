package coderange

import (
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/goni/err"
	"github.com/lyraproj/goni/goni"
	"github.com/lyraproj/goni/goni/syntax"
	"github.com/lyraproj/goni/util"
)

const (
	initMultiByteRangeSize = 5
	LastCodePoint          = 0x7fffffff
)

type Buffer struct {
	p []int
}

func NewBuffer() *Buffer {
	c := &Buffer{make([]int, 0, initMultiByteRangeSize)}
	c.writeCodePoint(0, 0)
	return c
}

func (c *Buffer) AppendTo(w *util.Indenter) {
	w.Append(`CodeRange`)
	w = w.Indent()
	w.NewLine()
	w.Append(`used: `)
	w.AppendInt(len(c.p))
	w.Append(`, size: `)
	sz := c.p[0]
	w.AppendInt(sz)
	w.NewLine()
	w.Append(`ranges:`)
	w = w.Indent()
	for i := 0; i < sz; i++ {
		if i%6 == 0 {
			w.NewLine()
		}
		w.Printf(`[0x%x..0x%x]`, c.p[i*2+1], c.p[i*2+2])
	}
}

func (b *Buffer) Copy() *Buffer {
	l := len(b.p)
	np := make([]int, l)
	copy(np, b.p)
	return &Buffer{np}
}

func (b *Buffer) Range() []int {
	return b.p
}

func (b *Buffer) move(from, to, n int) {
	if to+n > cap(b.p) {
		b.expand(to + n)
	}
	if to+n > len(b.p) {
		b.p = b.p[:to+n]
	}
	copy(b.p[to:], b.p[from:from+n])
}

func (b *Buffer) moveLeftAndReduce(from, to int) {
	copy(b.p[to:], b.p[from:])
	b.p = b.p[:len(b.p)-(from-to)]
}

func AddToBuffer(pbuf *Buffer, env goni.ScanEnvironment, from, to int, checkDup bool) *Buffer {
	if from > to {
		n := from
		from = to
		to = n
	}

	if pbuf == nil {
		pbuf = NewBuffer()
	}

	p := pbuf.p
	n := p[0]

	bound := n
	if from == 0 {
		bound = 0
	}
	low := 0
	for low < bound {
		x := (low + bound) >> 1
		if from-1 > p[x*2+2] {
			low = x + 1
		} else {
			bound = x
		}
	}

	high := low
	if to == LastCodePoint {
		high = n
	}
	bound = n
	for high < bound {
		x := (high + bound) >> 1
		if to+1 >= p[x*2+1] {
			high = x + 1
		} else {
			bound = x
		}
	}

	incN := low + 1 - high

	if n+incN > config.MaxMultiByteRangesNum {
		panic(err.NoArgs(err.TooManyMultiByteRanges))
	}

	if incN != 1 {
		if checkDup {
			if from <= p[low*2+2] && (p[low*2+1] <= from || p[low*2+2] <= to) {
				env.CCDuplicateWarning()
			}
		}

		if from > p[low*2+1] {
			from = p[low*2+1]
		}
		if to < p[(high-1)*2+2] {
			to = p[(high-1)*2+2]
		}
	}

	if incN != 0 {
		fromPos := 1 + high*2
		toPos := 1 + (low+1)*2

		if incN > 0 {
			if high < n {
				size := (n - high) * 2
				pbuf.move(fromPos, toPos, size)
			}
		} else {
			pbuf.moveLeftAndReduce(fromPos, toPos)
		}
	}

	pos := 1 + low*2
	// pbuf.ensureSize(pos + 2);
	pbuf.writeCodePoint(pos, from)
	pbuf.writeCodePoint(pos+1, to)
	n += incN
	pbuf.writeCodePoint(0, n)

	return pbuf

}

func Add(pbuf *Buffer, env goni.ScanEnvironment, from, to int, checkDup bool) *Buffer {
	if from > to {
		if env.Syntax().IsBehavior(syntax.AllowEmptyRangeInCC) {
			return pbuf
		} else {
			panic(err.NoArgs(err.EmptyRangeInCharClass))
		}
	}
	return AddToBuffer(pbuf, env, from, to, checkDup)
}

func mbcodeStartPosition(enc goni.Encoding) int {
	if enc.MinLength() > 1 {
		return 0
	}
	return 0x80
}

func (b *Buffer) writeCodePoint(pos, p int) {
	u := pos + 1
	if cap(b.p) < u {
		b.expand(u)
	}
	if len(b.p) < u {
		b.p = b.p[:u]
	}
	b.p[pos] = p
}

func (b *Buffer) expand(low int) {
	var c int
	for c = cap(b.p); c < low; c <<= 1 {
	}
	l := len(b.p)
	np := make([]int, l, c)
	copy(np, b.p)
	b.p = np
}

func setAllMultiByteRange(env goni.ScanEnvironment, pbuf *Buffer) *Buffer {
	return AddToBuffer(pbuf, env, mbcodeStartPosition(env.Encoding()), LastCodePoint, true)
}

func AddAllMultiByte(env goni.ScanEnvironment, pbuf *Buffer) *Buffer {
	if !env.Encoding().IsSingleByte() {
		return setAllMultiByteRange(env, pbuf)
	}
	return pbuf
}

func NotBuffer(env goni.ScanEnvironment, bbuf *Buffer) *Buffer {
	if bbuf == nil {
		return setAllMultiByteRange(env, nil)
	}

	p := bbuf.p
	n := p[0]

	if n <= 0 {
		return setAllMultiByteRange(env, nil)
	}

	pre := mbcodeStartPosition(env.Encoding())

	var pbuf *Buffer
	to := 0
	for i := 0; i < n; i++ {
		from := p[i*2+1]
		to = p[i*2+2]
		if pre <= from-1 {
			pbuf = AddToBuffer(pbuf, env, pre, from-1, true)
		}
		if to == LastCodePoint {
			break
		}
		pre = to + 1
	}

	if to < LastCodePoint {
		pbuf = AddToBuffer(pbuf, env, to+1, LastCodePoint, true)
	}
	return pbuf
}

func OrBuffer(env goni.ScanEnvironment, bbuf1 *Buffer, not1 bool, bbuf2 *Buffer, not2 bool) (pbuf *Buffer) {
	if bbuf1 == nil && bbuf2 == nil {
		if not1 || not2 {
			pbuf = setAllMultiByteRange(env, pbuf)
		}
		return
	}

	if bbuf2 == nil {
		// swap
		tnot := not1
		not1 = not2
		not2 = tnot
		tbuf := bbuf1
		bbuf1 = bbuf2
		bbuf2 = tbuf
	}

	if bbuf1 == nil {
		if not1 {
			pbuf = setAllMultiByteRange(env, pbuf)
		} else {
			if !not2 {
				pbuf = bbuf2.Copy()
			} else {
				pbuf = NotBuffer(env, bbuf2)
			}
		}
		return
	}

	if not1 {
		// swap
		tnot := not1
		not1 = not2
		not2 = tnot
		tbuf := bbuf1
		bbuf1 = bbuf2
		bbuf2 = tbuf
	}

	if !not2 && !not1 { /* 1 OR 2 */
		pbuf = bbuf2.Copy()
	} else if !not1 { /* 1 OR (not 2) */
		pbuf = NotBuffer(env, bbuf2)
	}

	p1 := bbuf1.p
	n1 := p1[0]

	for i := 0; i < n1; i++ {
		from := p1[i*2+1]
		to := p1[i*2+2]
		pbuf = AddToBuffer(pbuf, env, from, to, true)
	}
	return
}

func andCodeRange1(pbuf *Buffer, env goni.ScanEnvironment, from1, to1 int, data []int, n int) *Buffer {
	for i := 0; i < n; i++ {
		from2 := data[i*2+1]
		to2 := data[i*2+2]
		if from2 < from1 {
			if to2 < from1 {
				continue
			}
			from1 = to2 + 1
		} else if from2 <= to1 {
			if to2 < to1 {
				if from1 <= from2-1 {
					pbuf = AddToBuffer(pbuf, env, from1, from2-1, true)
				}
				from1 = to2 + 1
			} else {
				to1 = from2 - 1
			}
		} else {
			from1 = from2
		}
		if from1 > to1 {
			break
		}
	}

	if from1 <= to1 {
		pbuf = AddToBuffer(pbuf, env, from1, to1, true)
	}

	return pbuf
}

func AndBuffer(bbuf1 *Buffer, not1 bool,
	bbuf2 *Buffer, not2 bool, env goni.ScanEnvironment) (pbuf *Buffer) {
	if bbuf1 == nil {
		if not1 && bbuf2 != nil {
			pbuf = bbuf2.Copy()
		}
		return
	}
	if bbuf2 == nil {
		if not2 {
			pbuf = bbuf1.Copy()
		}
		return
	}

	if not1 {
		// swap
		tnot := not1
		not1 = not2
		not2 = tnot
		tbuf := bbuf1
		bbuf1 = bbuf2
		bbuf2 = tbuf
	}

	p1 := bbuf1.p
	n1 := p1[0]
	p2 := bbuf2.p
	n2 := p2[0]

	if !not2 && !not1 { /* 1 AND 2 */
		for i := 0; i < n1; i++ {
			from1 := p1[i*2+1]
			to1 := p1[i*2+2]

			for j := 0; j < n2; j++ {
				from2 := p2[j*2+1]
				to2 := p2[j*2+2]

				if from2 > to1 {
					break
				}
				if to2 < from1 {
					continue
				}
				from := from2
				if from1 > from2 {
					from = from1
				}
				to := to2
				if to1 < to2 {
					to = to1
				}
				pbuf = AddToBuffer(pbuf, env, from, to, true)
			}
		}
	} else if !not1 { /* 1 AND (not 2) */
		for i := 0; i < n1; i++ {
			from1 := p1[i*2+1]
			to1 := p1[i*2+2]
			pbuf = andCodeRange1(pbuf, env, from1, to1, p2, n2)
		}
	}
	return
}

func IsInCodeRange(p []int, code int) bool {
	return IsInCodeRange2(p, 0, code)
}

func IsInCodeRange2(p []int, offset, code int) bool {
	low := 0
	n := p[offset]
	high := n

	for low < high {
		x := (low + high) >> 1
		if code > p[(x<<1)+2+offset] {
			low = x + 1
		} else {
			high = x
		}
	}
	return low < n && code >= p[(low<<1)+1+offset]
}
