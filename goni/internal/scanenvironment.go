package internal

import "github.com/lyraproj/goni/goni"

type ScanEnvironment struct {

}

type WarnCallback func(string)

func NewScanEnvironment(regex *Regex, syntax *goni.Syntax, warning WarnCallback) goni.ScanEnvironment {

}

func (s *ScanEnvironment) NumMem() int {
	// TODO:
	return 0
}

func (s *ScanEnvironment) MemNodes() []goni.Node {
	return nil

}