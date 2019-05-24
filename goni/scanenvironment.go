package goni

type ScanEnvironment interface {
	UnknownEscapeWarning(s string)

	CloseBracketWithoutEscapeWarning(s string)

	CcDuplicateWarning()

	MemNodes() []Node

	NumMem() int
}
