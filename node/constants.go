package node

type Type uint

const (
	BeginBuf      = 1 << 0
	BeginLine     = 1 << 1
	BeginPosition = 1 << 2
	EndBuf        = 1 << 3
	SemiEndBuf    = 1 << 4
	EndLine       = 1 << 5

	WordBound    = 1 << 6
	NotWordBound = 1 << 7
	WordBegin    = 1 << 8
	WordEnd      = 1 << 9
	PrecRead     = 1 << 10
	PrecReadNot  = 1 << 11
	LookBehind   = 1 << 12
	LookBeindNot = 1 << 13

	AnycharStar   = 1 << 14
	AnycharStarMl = 1 << 15

	AnycharStartMask = AnycharStar | AnycharStarMl
	EndBufMask       = EndBuf | SemiEndBuf
	Keep             = 1 << 16

	AllowedInLb = LookBehind |
		LookBeindNot |
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

	TypeStr     = Type(0)
	TypeCClass  = Type(1)
	TypeCType   = Type(2)
	TypeCAny    = Type(3)
	TypeBref  = Type(4)
	TypeQtfr    = Type(5)
	TypeEnclose = Type(6)
	TypeAnchor  = Type(7)
	TypeList    = Type(8)
	TypeAlt     = Type(9)
	TypeCall    = Type(10)

	TypeTop = Type(15)

	BitStr     = 1 << TypeStr
	BitCClass  = 1 << TypeCClass
	BitCType   = 1 << TypeCType
	BitCAny    = 1 << TypeCAny
	BitBref    = 1 << TypeBref
	BitQtfr    = 1 << TypeQtfr
	BitEnclose = 1 << TypeEnclose
	BitAnchor  = 1 << TypeAnchor
	BitList    = 1 << TypeList
	BitAlt     = 1 << TypeAlt
	BitCall    = 1 << TypeCall

	BitAllowedInLb = BitList |
		BitAlt |
		BitStr |
		BitCClass |
		BitCType |
		BitCAny |
		BitAnchor |
		BitEnclose |
		BitQtfr |
		BitCall

	BitSimple = BitStr |
		BitCClass |
		BitCType |
		BitCAny |
		BitBref

	StateMinFixed           = 1 << 0
	StateMaxFixed           = 1 << 1
	StateCLenFixed          = 1 << 2
	StateMark1              = 1 << 3
	StateMark2              = 1 << 4
	StateMemBackrefed       = 1 << 5
	StateStopBtSimpleRepeat = 1 << 6
	StateRecursion          = 1 << 7
	StateCalled             = 1 << 8
	StateAddrFixed          = 1 << 9
	StateNamedGroup         = 1 << 10
	StateNameRef            = 1 << 11
	StateInRepeat           = 1 << 12 /* STK_REPEAT is nested in stack. */
	StateNestLevel          = 1 << 13
	StateByNumber           = 1 << 14 /* {n,m} */
)
