package goni

const (
CharTableSize          = 256
UseNoInvalidQuantifier = true
ScanenvMemnodesSize    = 8

UseNamedGroup       = true
UseSubexpCall       = true
UsePerlSubexpCall   = true
UseBackrefWithLevel = true /* \k<name+n>, \k<name-n> */

UseMonomaniacCheckCapturesInEndlessRepeat = true /* /(?:()|())*\2/ */
UseNewlineAtEndOfStringHasEmptyLine       = true /* /\n$/ =~ "\n" */
UseWarningRedundantNestedRepeatOperator   = true

CaseFoldIsAppliedInsideNegativeCclass = true

UseMatchRangeMustBeInsideOfSpecifiedRange = false
UseCaptureHistory                         = false
UseVariableMetaChars                      = true
B                                         = true /* "\<": word-begin, "\>": word-end */
UseFindLongestSearchAllOfRange            = true
UseSundayQuickSearch                      = true
UseCec                                    = false
UseDynamicOption                          = false
// TODO: USE_BYTE_MAP = OptExactInfo.OPT_EXACT_MAXLEN <= CharTableSize
UseIntMapBackward = false

Nregion               = 10
MaxBackrefNum         = 1000
MaxCaptureGroupNum    = 32767
MaxRepeatNum          = 100000
MaxMultiByteRangesNum = 10000

// internal config
UseOpPushOrJumpExact = true
UseQtfrPeekNext      = true

InitMatchStackSize = 64

DontOptimize = false

UseStringTemplates = true // use embedded string templates in Regex object as byte arrays instead of compiling them into int bytecode array


MaxCaptureHistoryGroup = 31


CheckStringThresholdLen = 7
CheckBuffMaxSize        = 0x4000

DebugAll = false

DEBUG                    = DebugAll
DebugParseTree           = DebugAll
DebugParseTreeRaw        = true
DebugCompile             = DebugAll
DebugCompileByteCodeInfo = DebugAll
DebugSearch              = DebugAll
DebugMatch               = DebugAll
)