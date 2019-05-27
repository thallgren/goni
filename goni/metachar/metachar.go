package metachar

const (
	Escape = iota
	AnyChar
	AnyTime
	ZeroOrOneTime
	OneOrMoreTime
	AnyCharAnyTime

	InnefectiveMetaChar = 0
)

type Table struct {
	Esc int
	AnyChar int
	AnyTime int
	ZeroOrOneTime int
	OneOrMoreTime int
	AnyCharAnyTime int
}

