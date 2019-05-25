package goni

import (
	"github.com/lyraproj/goni/util"
)

const (
	BitsPerByte    = 8
	SingleByteSize = 1 << BitsPerByte
	BitsInRoom     = 4 * BitsPerByte
	BitSetSize     = SingleByteSize / BitsInRoom
)

var RoomShift = log2(BitsInRoom)

func log2(n int) uint {
	i := uint(n)
	log := uint(0)
	for i > 0 {
		log++
		i = i >> 1
	}
	return log
}

func bit(pos int) int {
	return int(1 << (uint(pos) % SingleByteSize))
}

type BitSet struct {
	bits []int
}

func NewBitSet() *BitSet {
	return &BitSet{make([]int, BitSetSize)}
}

const bitsToStringWrap = 4

func (b *BitSet) AppendTo(w *util.Indenter) {
	w.Append(`BitSet`)
	for i := 0; i < SingleByteSize; i++ {
		if (i % (SingleByteSize / bitsToStringWrap)) == 0 {
			w.NewLine()
		}
		c := '0'
		if b.At(i) {
			c = '1'
		}
		w.AppendRune(c)
	}
}

func (b *BitSet) At(pos int) bool {
	return (b.bits[uint(pos)>>RoomShift] & bit(pos)) != 0
}

func (b *BitSet) RoomAt(pos int) int {
	return b.bits[pos]
}

func (b *BitSet) CheckedSet(env ScanEnvironment, pos int) {
	if b.At(pos) {
		env.CCDuplicateWarning()
	}
	b.Set(pos)
}

func (b *BitSet) CheckedSetRange(env ScanEnvironment, from, to int) {
	for i := from; i <= to && i < SingleByteSize; i++ {
		b.CheckedSet(env, i)
	}
}

func (b *BitSet) Set(pos int) {
	b.bits[pos>>RoomShift] |= bit(pos)
}

func (b *BitSet) SetAll() {
	for i := range b.bits {
		b.bits[i] = ^0
	}
}

func (b *BitSet) SetRange(from, to int) {
	for i := from; i <= to && i < SingleByteSize; i++ {
		b.Set(i)
	}
}

func (b *BitSet) Clear(pos int) {
	b.bits[pos>>RoomShift] &= ^bit(pos)
}

func (b *BitSet) Invert(pos int) {
	b.bits[pos>>RoomShift] ^= bit(pos)
}

func (b *BitSet) InvertAll() {
	for i := range b.bits {
		b.bits[i] = ^b.bits[i]
	}
}

func (b *BitSet) InvertTo(to *BitSet) {
	for i := range b.bits {
		to.bits[i] = ^b.bits[i]
	}
}

func (b *BitSet) And(other *BitSet) {
	for i := range b.bits {
		b.bits[i] &= other.bits[i]
	}
}

func (b *BitSet) Or(other *BitSet) {
	for i := range b.bits {
		b.bits[i] |= other.bits[i]
	}
}

func (b *BitSet) Copy(other *BitSet) {
	for i := range b.bits {
		b.bits[i] = other.bits[i]
	}
}

func (b *BitSet) ClearAll() {
	for i := range b.bits {
		b.bits[i] = 0
	}
}

func (b *BitSet) NumOn() int {
	count := 0
	for i := 0; i < SingleByteSize; i++ {
		if b.At(i) {
			count++
		}
	}
	return count
}

func (b *BitSet) IsEmpty() bool {
	for i := range b.bits {
		if b.bits[i] != 0 {
			return false
		}
	}
	return true
}

// BitCount returns the number of one-bits in the two's complement binary representation of
// the specified value.  This function is sometimes referred to as the <i>population count</i>.
func BitCount(n int) int {
	i := uint(n)
	// HD, Figure 5-2
	i = i - ((i >> 1) & 0x55555555)
	i = (i & 0x33333333) + ((i >> 2) & 0x33333333)
	i = (i + (i >> 4)) & 0x0f0f0f0f
	i = i + (i >> 8)
	i = i + (i >> 16)
	return int(i & 0x3f)
}
