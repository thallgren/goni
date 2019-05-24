package anchor

import (
	"github.com/lyraproj/goni/util"
	"io"
)

type Type int

const (
	BeginBuf      = Type(1 << 0)
	BeginLine     = Type(1 << 1)
	BeginPosition = Type(1 << 2)
	EndBuf        = Type(1 << 3)
	SemiEndBuf    = Type(1 << 4)
	EndLine       = Type(1 << 5)

	WordBound     = Type(1 << 6)
	NotWordBound  = Type(1 << 7)
	WordBegin     = Type(1 << 8)
	WordEnd       = Type(1 << 9)
	PrecRead      = Type(1 << 10)
	PrecReadNot   = Type(1 << 11)
	LookBehind    = Type(1 << 12)
	LookBehindNot = Type(1 << 13)

	AnyCharStar   = Type(1 << 14)
	AnyCharStarMl = Type(1 << 15)

	AnyCharStartMask = AnyCharStar | AnyCharStarMl
	EndBufMask       = EndBuf | SemiEndBuf
	Keep             = Type(1 << 16)

	AllowedInLb = LookBehind |
		LookBehindNot |
		BeginLine |
		EndLine |
		BeginBuf |
		BeginPosition |
		Keep |
		WordBound |
		NotWordBound |
		WordBegin |
		WordEnd

	AllowedInLbNot = AllowedInLb
)

func (a Type) AppendString(w io.Writer) {
	if a.IsType(BeginBuf) {
		util.WriteString(w, `BEGIN_BUF `)
	}
	if a.IsType(BeginLine) {
		util.WriteString(w, `BEGIN_LINE `)
	}
	if a.IsType(BeginPosition) {
		util.WriteString(w, `BEGIN_POSITION `)
	}
	if a.IsType(EndBuf) {
		util.WriteString(w, `END_BUF `)
	}
	if a.IsType(SemiEndBuf) {
		util.WriteString(w, `SEMI_END_BUF `)
	}
	if a.IsType(EndLine) {
		util.WriteString(w, `END_LINE `)
	}
	if a.IsType(WordBound) {
		util.WriteString(w, `WORD_BOUND `)
	}
	if a.IsType(NotWordBound) {
		util.WriteString(w, `NOT_WORD_BOUND `)
	}
	if a.IsType(WordBegin) {
		util.WriteString(w, `WORD_BEGIN `)
	}
	if a.IsType(WordEnd) {
		util.WriteString(w, `WORD_END `)
	}
	if a.IsType(PrecRead) {
		util.WriteString(w, `PREC_READ `)
	}
	if a.IsType(PrecReadNot) {
		util.WriteString(w, `PREC_READ_NOT `)
	}
	if a.IsType(LookBehind) {
		util.WriteString(w, `LOOK_BEHIND `)
	}
	if a.IsType(LookBehindNot) {
		util.WriteString(w, `LOOK_BEHIND_NOT `)
	}
	if a.IsType(AnyCharStar) {
		util.WriteString(w, `ANYCHAR_STAR `)
	}
	if a.IsType(AnyCharStarMl) {
		util.WriteString(w, `ANYCHAR_STAR_ML `)
	}
}

func (t Type) IsType(ot Type) bool {
	return (t & ot) != 0
}
