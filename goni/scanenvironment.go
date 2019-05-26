package goni

type ScanEnvironment interface {
	UnknownEscapeWarning(s string)

	CloseBracketWithoutEscapeWarning(s string)

	CCDuplicateWarning()

	CCEscWarn(msg string)

	ConvertBackslashValue(c int) int

	Encoding() Encoding

	MemNodes() []Node

	NumMem() int

	Syntax() *Syntax

	Warnings() WarnCallback
}
