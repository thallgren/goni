package character

type Type uint

const (
	Newline = Type(iota)
	Alpha
	Blank
	Cntrl
	Digit
	Graph
	Lower
	Print
	Punct
	Space
	Upper
	Xdigit
	Word
	Alnum
	Ascii

	MaxStdCtype = Ascii

	BitNewline = Type(1 << Newline)
	BitAlpha   = Type(1 << Alpha)
	BitBlank   = Type(1 << Blank)
	BitCntrl   = Type(1 << Cntrl)
	BitDigit   = Type(1 << Digit)
	BitGraph   = Type(1 << Graph)
	BitLower   = Type(1 << Lower)
	BitPrint   = Type(1 << Print)
	BitPunct   = Type(1 << Punct)
	BitSpace   = Type(1 << Space)
	BitUpper   = Type(1 << Upper)
	BitXdigit  = Type(1 << Xdigit)
	BitWord    = Type(1 << Word)
	BitAlnum   = Type(1 << Alnum)
	BitAscii   = Type(1 << Ascii)
)
