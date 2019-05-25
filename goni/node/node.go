package node

type Type uint

const (
	Str     = Type(iota)
	CClass
	CType
	CAny
	Bref
	QTFR
	Enclose
	Anchor
	List
	Alt
	Call

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
