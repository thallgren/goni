package syntax

type Op int
type Op2 int
type Op3 int
type Behavior int

const (
	/* syntax (operators); */
	OpVariableMetaCharacters = Op(1 << 0)
	OpDotAnychar             = Op(1 << 1) /* . */
	OpAsteriskZeroInf        = Op(1 << 2) /* * */
	OpEscAsteriskZeroInf     = Op(1 << 3)
	OpPlusOneInf             = Op(1 << 4) /* + */
	OpEscPlusOneInf          = Op(1 << 5)
	OpQmarkZeroOne           = Op(1 << 6) /* ? */
	OpEscQmarkZeroOne        = Op(1 << 7)
	OpBraceInterval          = Op(1 << 8)  /* {lower,upper} */
	OpEscBraceInterval       = Op(1 << 9)  /* \{lower,upper\} */
	OpVbarAlt                = Op(1 << 10) /* | */
	OpEscVbarAlt             = Op(1 << 11) /* \| */
	OpLparenSubexp           = Op(1 << 12) /* (...);   */
	OpEscLparenSubexp        = Op(1 << 13) /* \(...\); */
	OpEscAzBufAnchor         = Op(1 << 14) /* \A, \Z, \z */
	OpEscCapitalGBeginAnchor = Op(1 << 15) /* \G     */
	OpDecimalBackref         = Op(1 << 16) /* \num   */
	OpBracketCC              = Op(1 << 17) /* [...]  */
	OpEscWWord               = Op(1 << 18) /* \w, \W */
	OpEscLtgtWordBeginEnd    = Op(1 << 19) /* \<. \> */
	OpEscBWordBound          = Op(1 << 20) /* \b, \B */
	OpEscSWhiteSpace         = Op(1 << 21) /* \s, \S */
	OpEscDDigit              = Op(1 << 22) /* \d, \D */
	OpLineAnchor             = Op(1 << 23) /* ^, $   */
	OpPosixBracket           = Op(1 << 24) /* [:xxxx:] */
	OpQMarkNonGreedy         = Op(1 << 25) /* ??,*?,+?,{n,m}? */
	OpEscControlChars        = Op(1 << 26) /* \n,\r,\t,\a ... */
	OpEscCControl            = Op(1 << 27) /* \cx  */
	OpEscOctal3              = Op(1 << 28) /* \OOO */
	OpEscXHex2               = Op(1 << 29) /* \xHH */
	OpEscXBraceHex8          = Op(1 << 30) /* \x{7HHHHHHH} */
	OpEscOBraceOctal         = Op(1 << 31) /* \o{OOO} */

	Op2EscCapitalQQuote       = Op2(1 << 0)  /* \Q...\E */
	Op2QmarkGroupEffect       = Op2(1 << 1)  /* (?...); */
	Op2OptionPerl             = Op2(1 << 2)  /* (?imsxadlu), (?-imsx), (?^imsxalu) */
	Op2OptionRuby             = Op2(1 << 3)  /* (?imxadu);, (?-imx);  */
	Op2PlusPossessiveRepeat   = Op2(1 << 4)  /* ?+,*+,++ */
	Op2PlusPossessiveInterval = Op2(1 << 5)  /* {n,m}+   */
	Op2CClassSetOp            = Op2(1 << 6)  /* [...&&..[..]..] */
	Op2QmarkLtNamedGroup      = Op2(1 << 7)  /* (?<name>...); */
	Op2EscKNamedBackref       = Op2(1 << 8)  /* \k<name> */
	Op2EscGSubexpCall         = Op2(1 << 9)  /* \g<name>, \g<n> */
	Op2AtmarkCaptureHistory   = Op2(1 << 10) /* (?@..);,(?@<x>..); */
	Op2EscCapitalCBarControl  = Op2(1 << 11) /* \C-x */
	Op2EscCapitalMBarMeta     = Op2(1 << 12) /* \M-x */
	Op2EscVVtab               = Op2(1 << 13) /* \v as VTAB */
	Op2EscUHex4               = Op2(1 << 14) /* \\uHHHH */
	Op2EscGnuBufAnchor        = Op2(1 << 15) /* \`, \' */
	Op2EscPBraceCharProperty  = Op2(1 << 16) /* \p{...}, \P{...} */
	Op2EscPBraceCircumflexNot = Op2(1 << 17) /* \p{^..}, \P{^..} */
	/* Op2CharOp2PrefixIs = Op2(1<<18) */
	Op2EscHXdigit                         = Op2(1 << 19) /* \h, \H */
	Op2IneffectiveEscape                  = Op2(1 << 20) /* \ */
	Op2EscCapitalRLinebreak               = Op2(1 << 21) /* \R as (?>\x0D\x0A|[\x0A-\x0D\x{85}\x{2028}\x{2029}]) */
	Op2EscCapitalXExtendedGraphemeCluster = Op2(1 << 22) /* \X as (?:\P{M}\p{M}*) */
	Op2EscVVerticalWhitespace             = Op2(1 << 23) /* \v, \V -- Perl */
	Op2EscHHorizontalWhitespace           = Op2(1 << 24) /* \h, \H -- Perl */
	Op2EscCapitalKKeep                    = Op2(1 << 25) /* \K */
	Op2EscGBraceBackref                   = Op2(1 << 26) /* \g{name}, \g{n} */
	Op2QmarkSubexpCall                    = Op2(1 << 27) /* (?&name), (?n), (?R), (?0) */
	Op2QmarkBarBranchReset                = Op2(1 << 28) /* (?|...) */
	Op2QmarkLparenCondition               = Op2(1 << 29) /* (?(cond)yes...|no...) */
	Op2QmarkCapitalPNamedGroup            = Op2(1 << 30) /* (?P<name>...), (?P=name), (?P>name) -- Python/PCRE */
	Op2QmarkTildeAbsent                   = Op2(1 << 31) /* (?~...) */

	Op3OptionJava       = Op3(1 << 0) /* (?idmsux), (?-idmsux) */
	Op3OptionECMAScript = Op3(1 << 1) /* EcmaScript quirks */

	/* syntax (behavior); */
	ContextIndepAnchors              = Behavior(1 << 31) /* not implemented */
	ContextIndepRepeatOps            = Behavior(1 << 0)  /* ?, *, +, {n,m} */
	ContextInvalidRepeatOps          = Behavior(1 << 1)  /* error or ignore */
	AllowUnmatchedCloseSubexp        = Behavior(1 << 2)  /* ...);... */
	AllowInvalidInterval             = Behavior(1 << 3)  /* {??? */
	AllowIntervalLowAbbrev           = Behavior(1 << 4)  /* {,n} => {0,n} */
	StrictCheckBackref               = Behavior(1 << 5)  /* /(\1);/,/\1();/ ..*/
	DifferentLenAltLookBehind        = Behavior(1 << 6)  /* (?<=a|bc); */
	CaptureOnlyNamedGroup            = Behavior(1 << 7)  /* see doc/RE */
	AllowMultiplexDefinitionName     = Behavior(1 << 8)  /* (?<x>);(?<x>); */
	FixedIntervalIsGreedyOnly        = Behavior(1 << 9)  /* a{n}?=(?:a{n});? */
	AllowMultiplexDefinitionNameCall = Behavior(1 << 10) /* (?<x>)(?<x>)(?&x) */

	/* syntax (behavior); in char class [...] */
	NotNewlineInNegativeCC = Behavior(1 << 20) /* [^...] */
	BackslashEscapeInCC    = Behavior(1 << 21) /* [..\w..] etc.. */
	AllowEmptyRangeInCC    = Behavior(1 << 22)
	AllowDoubleRangeOpInCC = Behavior(1 << 23) /* [0-9-a]=[0-9\-a] */
	/* syntax (behavior); warning */
	WarnCCOpNotEscaped        = Behavior(1 << 24) /* [,-,] */
	WarnRedundantNestedRepeat = Behavior(1 << 25) /* (?:a*);+ */
	WarnCCDup                 = Behavior(1 << 26) /* [aa] */

	PosixCommonOp = OpDotAnychar | OpPosixBracket |
		OpDecimalBackref |
		OpBracketCC | OpAsteriskZeroInf |
		OpLineAnchor |
		OpEscControlChars

	GnuRegexOp = OpDotAnychar | OpBracketCC |
		OpPosixBracket | OpDecimalBackref |
		OpBraceInterval | OpLparenSubexp |
		OpVbarAlt |
		OpAsteriskZeroInf | OpPlusOneInf |
		OpQmarkZeroOne |
		OpEscAzBufAnchor | OpEscCapitalGBeginAnchor |
		OpEscWWord |
		OpEscBWordBound | OpEscLtgtWordBeginEnd |
		OpEscSWhiteSpace | OpEscDDigit |
		OpLineAnchor

	GnuRegexBv = ContextIndepAnchors | ContextIndepRepeatOps |
		ContextInvalidRepeatOps | AllowInvalidInterval |
		BackslashEscapeInCC | AllowDoubleRangeOpInCC
)

func (o Op) IsSet(op Op) bool {
	return (o & op) != 0
}

func (o Op2) IsSet(op Op2) bool {
	return (o & op) != 0
}

func (o Op3) IsSet(op Op3) bool {
	return (o & op) != 0
}

func (o Behavior) IsSet(op Behavior) bool {
	return (o & op) != 0
}
