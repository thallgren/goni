package goni

type ScanEnvironment interface {
	UnknownEscapeWarning(s string)

	CloseBracketWithoutEscapeWarning(s string)

	CCDuplicateWarning()

	Encoding() Encoding

	MemNodes() []Node

	NumMem() int

	Syntax() *Syntax
}
