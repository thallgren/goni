package internal

import "github.com/lyraproj/goni/goni"

type Lexer struct {
	scannerSupport
	regex *Regex
	env goni.ScanEnvironment
	syntax goni.Syntax

}
