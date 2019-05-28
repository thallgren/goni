package goni

import "github.com/lyraproj/goni/goni/option"

type ScanEnvironment interface {
	Option() option.Type

	Encoding() Encoding

	Syntax() *Syntax

	NumMem() int

	MemNodes() []Node

	Warnings() WarnCallback
}
