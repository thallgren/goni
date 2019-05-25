package goni

import "github.com/lyraproj/goni/goni/character"

type Encoding interface {
	// CodeToMbcLength returns character length given a code point
	CodeToMbcLength(code int) int

	// CTypeCodeRange returns code range for a given character type
	CTypeCodeRange(ctype character.Type, sbOut *int) []int

	// Length returns character length given stream, character position and stream end returns 1 for
	// singlebyte encodings or performs sanity validations for multibyte ones and returns the character
	// length, missing characters in the stream otherwise
	Length(bytes []byte, p, end int) int

	// MbcToCode returns code point for a character
	MbcToCode(bytes []byte, p, end int) (code, len int)

	// MinLength returns minimum character byte length that can appear in an encoding
	MinLength() int

	// IsCodeCType performs a check whether given code is of given character type (e.g. used by
	// isWord(someByte) and similar methods)
	IsCodeCType(code int, ctype character.Type) bool

	IsDigit(code int) bool
	IsXDigit(code int) bool
	IsSbWord(code int) bool
	IsSingleByte() bool
	IsWord(code int) bool
	IsUpper(code int) bool
}
