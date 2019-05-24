package goni

import (
	"github.com/lyraproj/goni/util"
)

const (
	BitsPerByte    = uint(8)
	SingleByteSize = uint(1 << BitsPerByte)
	BitsInRoom     = uint(4 * BitsPerByte)
	BitSetSize     = uint(SingleByteSize / BitsInRoom)
)

var RoomShift = log2(BitsInRoom)

func log2(n uint) uint {
	log := uint(0)
	for n > 0 {
		log++
		n = n >> 1
	}
	return log
}

func bit(pos uint) uint {
	return 1 << (pos % SingleByteSize)
}

type BitSet struct {
	bits []uint
}

func NewBitSet() *BitSet {
	return &BitSet{make([]uint, BitSetSize)}
}

const bitsToStringWrap = 4

func (b *BitSet) AppendTo(w *util.Indenter) {
	w.Append(`BitSet`)
	for i := uint(0); i < SingleByteSize; i++ {
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

func (b *BitSet) At(pos uint) bool {
	return (b.bits[pos >> RoomShift] & bit(pos)) != 0;
}

func (b *BitSet) CheckedSet(env ScanEnvironment, pos uint) {
	if b.At(pos) {
		env.CcDuplicateWarning()
	}
	b.Set(pos)
}

func (b *BitSet) CheckedSetRange(env ScanEnvironment, from, to uint) {
	for i := from; i <= to && i < SingleByteSize; i++ {
		b.CheckedSet(env, i)
	}
}

func (b *BitSet) Set(pos uint) {
	b.bits[pos >> RoomShift] |= bit(pos)
}

func (b *BitSet) SetAll() {
	for i := range b.bits {
		b.bits[i] = ^0
	}
}

func (b *BitSet) SetRange(from, to uint) {
	for i := from; i <= to && i < SingleByteSize; i++ {
		b.Set(i)
	}
}

func (b *BitSet) Clear(pos uint) {
	b.bits[pos >> RoomShift] &= ^bit(pos)
}

func (b *BitSet) Invert(pos uint) {
	b.bits[pos >> RoomShift] ^= bit(pos)
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
	for i := uint(0); i < SingleByteSize; i++ {
		if b.At(i) {
			count++
		}
	}
	return count
}

func (b *BitSet) Empty() bool {
	for i := range b.bits {
		if b.bits[i] != 0 {
			return false
		}
	}
	return true
}
