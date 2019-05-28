package goni

import "github.com/lyraproj/goni/goni/character"

type Encoding interface {
	// CodeToMbc extracts code point into it's multibyte representation and appends
	// it to the given slice. The result of the append is returned
	CodeToMbc(code int, bytes []byte) []byte

	// CodeToMbcLength returns character length given a code point
	CodeToMbcLength(code int) int

	// CTypeCodeRange returns code range for a given character type
	CTypeCodeRange(ctype character.Type, sbOut *int) []int

	// Length returns character length given stream, character position and stream end returns 1 for
	// singlebyte encodings or performs sanity validations for multibyte ones and returns the character
	// length, missing characters in the stream otherwise
	Length(bytes []byte, p, end int) int

	StrNCmp(bytes []byte, p, end int, ascii []byte, asciiP, n int) int

	StrLength(bytes []byte, p, end int) int

	Step(bytes []byte, p, end, n int) int

	// MbcToCode returns code point for a character
	MbcToCode(bytes []byte, p, end int) (code, len int)

	// MaxLength Returns maximum character byte length that can appear in an encoding
	MaxLength() int

	// MinLength returns minimum character byte length that can appear in an encoding
	MinLength() int

	PrevCharHead(bytes []byte, p, s, end int) int

	PropertyNameToCType(name []byte, p, end int) character.Type

	// IsCodeCType performs a check whether given code is of given character type (e.g. used by
	// isWord(someByte) and similar methods)
	IsCodeCType(code int, ctype character.Type) bool

	IsDigit(code int) bool
	IsXDigit(code int) bool
	IsSbWord(code int) bool
	IsSingleByte() bool
	IsWord(code int) bool
	IsUpper(code int) bool
	IsNewLine(code int) bool
	IsUnicode() bool
}
