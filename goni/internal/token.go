package internal

type TokenType int

const (
	TkEOT = TokenType(iota) /* end of token */
	TkRawByte
	TkChar
	TkString
	TkCodePoint
	TkAnyChar
	TkCharType
	TkBackRef
	TkCall
	TkAnchor
	TkOpRepeat
	TkInterval
	TkAnycharAnytime  /* SQL '%' == .* */
	TkAlt
	TkSubexpOpen
	TkSubexpClose
	TkCcOpen
	TkQuoteOpen
	TkCharProperty  /* \p{...}, \P{...} */
	TkLineBreak
	TkExtendedGraphemeCluster
	TkKeep
	/* in cc */
	TkCcClose
	TkCcRange
	TkPosixBracketOpen
	TkCcAnd     /* && */
	TkCcCcOpen  /* [ */
)
