package node

type Type uint

const (
	Str     = Type(0)
	CClass  = Type(1)
	CType   = Type(2)
	CAny    = Type(3)
	Bref    = Type(4)
	QTFR    = Type(5)
	Enclose = Type(6)
	Anchor  = Type(7)
	List    = Type(8)
	Alt     = Type(9)
	Call    = Type(10)

	Top = Type(15)

	BitStr     = 1 << Str
	BitCClass  = 1 << CClass
	BitCType   = 1 << CType
	BitCAny    = 1 << CAny
	BitBref    = 1 << Bref
	BitQTFR    = 1 << QTFR
	BitEnclose = 1 << Enclose
	BitAnchor  = 1 << Anchor
	BitList    = 1 << List
	BitAlt     = 1 << Alt
	BitCall    = 1 << Call

	AllowedInLb = BitList |
			BitAlt |
			BitStr |
			BitCClass |
			BitCType |
			BitCAny |
			BitAnchor |
			BitEnclose |
			BitQTFR |
			BitCall

	Simple = BitStr |
			BitCClass |
			BitCType |
			BitCAny |
			BitBref
)
