package config

//noinspection GoBoolExpressions
const (
	// JCodings config
	/* work size */
	EncCodeToMbcMaxlen   = 7
	EncMbcCaseFoldMaxlen = 18

	EncMaxCompCaseFoldCodeLen = 3
	EncGetCaseFoldCodesMaxNum = 13 /* 13 => Unicode:0x1ffc */

	UseUnicodeCaseFoldTurkishAzeri = false
	UseUnicodeAllLineTerminators   = false
	UseCrnlAsLineTerminator        = false

	UseUnicodeProperties = true

	CodePointMaskWidth = 3
	CodePointMask      = (1 << CodePointMaskWidth) - 1

	SpecialIndexShift = 3
	SpecialIndexWidth = 10
	SpecialIndexMask  = ((1 << SpecialIndexWidth) - 1) << SpecialIndexShift

	SpecialsLengthOffset = 25

	CaseUpcase        = 1 << 13 /* has/needs uppercase mapping */
	CaseDowncase      = 1 << 14 /* has/needs lowercase mapping */
	CaseTitlecase     = 1 << 15 /* has/needs (special) titlecase mapping */
	CaseSpecialOffset = 3       /* offset in bits from ONIGENC_CASE to ONIGENC_CASE_SPECIAL */
	CaseUpSpecial     = 1 << 16 /* has special upcase mapping */
	CaseDownSpecial   = 1 << 17 /* has special downcase mapping */
	CaseModified      = 1 << 18 /* data has been modified */
	CaseFold          = 1 << 19 /* has/needs case folding */

	CaseFoldTurkishAzeri = 1 << 20 /* needs mapping specific to Turkic languages; better not change original value! */

	CaseFoldLithuanian = 1 << 21 /* needs Lithuanian-specific mapping */
	CaseAsciiOnly      = 1 << 22 /* only modify ASCII range */
	CaseIsTitlecase    = 1 << 23 /* character itself is already titlecase */
	CaseSpecials       = CaseTitlecase | CaseIsTitlecase | CaseUpSpecial | CaseDownSpecial

	InternalEncCaseFoldMultiChar = 1 << 30 /* better not change original value! */
	EncCaseFoldMin               = InternalEncCaseFoldMultiChar
	EncCaseFoldDefault           = EncCaseFoldMin

	// Joni config
	CharTableSize          = 256
	UseNoInvalidQuantifier = true
	ScanEnvMemNodesSize    = 8

	UseNamedGroup       = true
	UseSubExpCall       = true
	UsePerlSubExpCall   = true
	UseBackrefWithLevel = true /* \k<name+n>, \k<name-n> */

	UseMonomaniacCheckCapturesInEndlessRepeat = true /* /(?:()|())*\2/ */
	UseNewlineAtEndOfStringHasEmptyLine       = true /* /\n$/ =~ "\n" */
	UseWarningRedundantNestedRepeatOperator   = true

	CaseFoldIsAppliedInsideNegativeCClass = true

	UseMatchRangeMustBeInsideOfSpecifiedRange = false
	UseCaptureHistory                         = false
	UseVariableMetaChars                      = true
	UseWordBeginEnd                           = true /* "\<": word-begin, "\>": word-end */
	UseFindLongestSearchAllOfRange            = true
	UseSundayQuickSearch                      = true
	UseCec                                    = false
	UseDynamicOption                          = false
	// TODO: USE_BYTE_MAP = OptExactInfo.OPT_EXACT_MAXLEN <= CharTableSize
	UseIntMapBackward = false

	NRegion               = 10
	MaxBackrefNum         = 1000
	MaxCaptureGroupNum    = 32767
	MaxRepeatNum          = 100000
	MaxMultiByteRangesNum = 10000

	// internal config
	UseOpPushOrJumpExact = true
	UseQTFRPeekNext      = true

	InitMatchStackSize = 64

	DontOptimize = false

	UseStringTemplates = true // use embedded string templates in Regex object as byte arrays instead of compiling them into int byte-code array

	MaxCaptureHistoryGroup = 31

	CheckStringThresholdLen = 7
	CheckBuffMaxSize        = 0x4000

	DebugAll = false

	Debug                    = DebugAll
	DebugParseTree           = DebugAll
	DebugParseTreeRaw        = true
	DebugCompile             = DebugAll
	DebugCompileByteCodeInfo = DebugAll
	DebugSearch              = DebugAll
	DebugMatch               = DebugAll
)
