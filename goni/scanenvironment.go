package goni

import (
	"github.com/lyraproj/goni/goni/option"
)

type ScanEnvironment interface {
	AddMemEntry() int

	UnknownEscapeWarning(s string)

	CloseBracketWithoutEscapeWarning(s string)

	CCDuplicateWarning()

	CCEscWarn(msg string)

	ConvertBackslashValue(c int) int

	Encoding() Encoding

	MemNodes() []Node

	NumMem() int

	Option() option.Type

	SetOption(option.Type)

	CurrentPrecReadNotNode() Node

	PushPrecReadNotNode(node Node)

	PopPrecReadNotNode(node Node)

	Syntax() *Syntax

	Warnings() WarnCallback

	SetMemNode(i int, node Node)

}
