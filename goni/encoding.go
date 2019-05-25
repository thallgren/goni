package goni

type Encoding interface {
	MinLength() int
	IsSingleByte() bool
}
