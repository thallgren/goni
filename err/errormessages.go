package err

import (
	"github.com/lyraproj/goni/config"
	"github.com/lyraproj/issue/issue"
)

const (
	/* internal error */
	ParserBug            = issue.Code(`PARSER_BUG`)
	UndefinedBytecode    = issue.Code(`UNDEFINED_BYTECODE`)
	UnexpectedBytecode   = issue.Code(`UNEXPECTED_BYTECODE`)
	TooManyCaptureGroups = issue.Code(`TOO_MANY_CAPTURE_GROUPS`)

	/* general error */
	InvalidArgument = issue.Code(`INVALID_ARGUMENT`)

	/* syntax error */
	EndPatternAtLeftBrace              = issue.Code(`END_PATTERN_AT_LEFT_BRACE`)
	EndPatternAtLeftBracket            = issue.Code(`END_PATTERN_AT_LEFT_BRACKET`)
	EmptyCharClass                     = issue.Code(`EMPTY_CHAR_CLASS`)
	PrematureEndOfCharClass            = issue.Code(`PREMATURE_END_OF_CHAR_CLASS`)
	EndPatternAtEscape                 = issue.Code(`END_PATTERN_AT_ESCAPE`)
	EndPatternAtMeta                   = issue.Code(`END_PATTERN_AT_META`)
	EndPatternAtControl                = issue.Code(`END_PATTERN_AT_CONTROL`)
	MetaCodeSyntax                     = issue.Code(`META_CODE_SYNTAX`)
	ControlCodeSyntax                  = issue.Code(`CONTROL_CODE_SYNTAX`)
	CharClassValueAtEndOfRange         = issue.Code(`CHAR_CLASS_VALUE_AT_END_OF_RANGE`)
	CharClassValueAtStartOfRange       = issue.Code(`CHAR_CLASS_VALUE_AT_START_OF_RANGE`)
	UnmatchedRangeSpecifierInCharClass = issue.Code(`UNMATCHED_RANGE_SPECIFIER_IN_CHAR_CLASS`)
	TargetOfRepeatOperatorNotSpecified = issue.Code(`TARGET_OF_REPEAT_OPERATOR_NOT_SPECIFIED`)
	TargetOfRepeatOperatorInvalid      = issue.Code(`TARGET_OF_REPEAT_OPERATOR_INVALID`)
	NestedRepeatNotAllowed             = issue.Code(`NESTED_REPEAT_NOT_ALLOWED`)
	NestedRepeatOperator               = issue.Code(`NESTED_REPEAT_OPERATOR`)
	UnmatchedCloseParenthesis          = issue.Code(`UNMATCHED_CLOSE_PARENTHESIS`)
	EndPatternWithUnmatchedParenthesis = issue.Code(`END_PATTERN_WITH_UNMATCHED_PARENTHESIS`)
	EndPatternInGroup                  = issue.Code(`END_PATTERN_IN_GROUP`)
	UndefinedGroupOption               = issue.Code(`UNDEFINED_GROUP_OPTION`)
	InvalidPosixBracketType            = issue.Code(`INVALID_POSIX_BRACKET_TYPE`)
	InvalidLookBehindPattern           = issue.Code(`INVALID_LOOK_BEHIND_PATTERN`)
	InvalidRepeatRangePattern          = issue.Code(`INVALID_REPEAT_RANGE_PATTERN`)
	InvalidConditionPattern            = issue.Code(`INVALID_CONDITION_PATTERN`)

	/* values error (syntax error) */
	TooBigNumber                       = issue.Code(`TOO_BIG_NUMBER`)
	TooBigNumberForRepeatRange         = issue.Code(`TOO_BIG_NUMBER_FOR_REPEAT_RANGE`)
	UpperSmallerThanLowerInRepeatRange = issue.Code(`UPPER_SMALLER_THAN_LOWER_IN_REPEAT_RANGE`)
	EmptyRangeInCharClass              = issue.Code(`EMPTY_RANGE_IN_CHAR_CLASS`)
	MismatchCodeLengthInClassRange     = issue.Code(`MISMATCH_CODE_LENGTH_IN_CLASS_RANGE`)
	TooManyMultiByteRanges             = issue.Code(`TOO_MANY_MULTI_BYTE_RANGES`)
	TooShortMultiByteString            = issue.Code(`TOO_SHORT_MULTI_BYTE_STRING`)
	TooBigBackrefNumber                = issue.Code(`TOO_BIG_BACKREF_NUMBER`)
	InvalidBackref                     = issue.Code(`INVALID_BACKREF`)
	NumberedBackrefOrCallNotAllowed    = issue.Code(`NUMBERED_BACKREF_OR_CALL_NOT_ALLOWED`)
	TooShortDigits                     = issue.Code(`TOO_SHORT_DIGITS`)
	InvalidWideCharValue               = issue.Code(`INVALID_WIDE_CHAR_VALUE`)
	EmptyGroupName                     = issue.Code(`EMPTY_GROUP_NAME`)
	InvalidGroupName                   = issue.Code(`INVALID_GROUP_NAME`)
	InvalidCharInGroupName             = issue.Code(`INVALID_CHAR_IN_GROUP_NAME`)
	UndefinedNameReference             = issue.Code(`UNDEFINED_NAME_REFERENCE`)
	UndefinedGroupReference            = issue.Code(`UNDEFINED_GROUP_REFERENCE`)
	MultiplexDefinedName               = issue.Code(`MULTIPLEX_DEFINED_NAME`)
	MultiplexDefinitionNameCall        = issue.Code(`MULTIPLEX_DEFINITION_NAME_CALL`)
	NeverEndingRecursion               = issue.Code(`NEVER_ENDING_RECURSION`)
	GroupNumberOverForCaptureHistory   = issue.Code(`GROUP_NUMBER_OVER_FOR_CAPTURE_HISTORY`)
	NotSupportedEncodingCombination    = issue.Code(`NOT_SUPPORTED_ENCODING_COMBINATION`)
	InvalidCombinationOfOptions        = issue.Code(`INVALID_COMBINATION_OF_OPTIONS`)
	OverThreadPassLimitCount           = issue.Code(`OVER_THREAD_PASS_LIMIT_COUNT`)
	TooBigSbCharValue                  = issue.Code(`TOO_BIG_SB_CHAR_VALUE`)

	CCTypeBug = issue.Code(`CCTYPE_BUG`)

	CCTooBigWideCharValue  = issue.Code(`CCTOO_BIG_WIDE_CHAR_VALUE`)
	CCTooLongWideCharValue = issue.Code(`CCTOO_LONG_WIDE_CHAR_VALUE`)

	CCInvalidCharPropertyName = issue.Code(`CCINVALID_CHAR_PROPERTY_NAME`)
	CCInvalidCodePointValue   = issue.Code(`CCINVALID_CODE_POINT_VALUE`)

	CCEncodingClassDefNotFound = issue.Code(`CCENCODING_CLASS_DEF_NOT_FOUND`)
	CCEncodingLoadError        = issue.Code(`CCENCODING_LOAD_ERROR`)

	CCIllegalCharacter = issue.Code(`CCILLEGAL_CHARACTER`)

	CCEncodingAlreadyRegistered        = issue.Code(`CCENCODING_ALREADY_REGISTERED`)
	CCEncodingAliasAlreadyRegistered   = issue.Code(`CCENCODING_ALIAS_ALREADY_REGISTERED`)
	CCEncodingReplicaAlreadyRegistered = issue.Code(`CCENCODING_REPLICA_ALREADY_REGISTERED`)
	CCNoSuchEncodng                    = issue.Code(`CCNO_SUCH_ENCODNG`)
	CCCouldNotReplicate                = issue.Code(`CCCOULD_NOT_REPLICATE`)

	// transcoder messages
	CCTranscoderAlreadyRegistered = issue.Code(`CCTRANSCODER_ALREADY_REGISTERED`)
	CCTranscoderClassDefNotFound  = issue.Code(`CCTRANSCODER_CLASS_DEF_NOT_FOUND`)
	CCTranscoderLoadError         = issue.Code(`CCTRANSCODER_LOAD_ERROR`)
)

func init() {
	/* internal error */
	issue.Hard(ParserBug, `internal parser error (bug)`)
	issue.Hard(UndefinedBytecode, `undefined bytecode (bug)`)
	issue.Hard(UnexpectedBytecode, `unexpected bytecode (bug)`)
	issue.Hard(TooManyCaptureGroups, `too many capture groups are specified`)

	/* general error */
	issue.Hard(InvalidArgument, `invalid argument`)

	/* syntax error */
	issue.Hard(EndPatternAtLeftBrace, `end pattern at left brace`)
	issue.Hard(EndPatternAtLeftBracket, `end pattern at left bracket`)
	issue.Hard(EmptyCharClass, `empty char-class`)
	issue.Hard(PrematureEndOfCharClass, `premature end of char-class`)
	issue.Hard(EndPatternAtEscape, `end pattern at escape`)
	issue.Hard(EndPatternAtMeta, `end pattern at meta`)
	issue.Hard(EndPatternAtControl, `end pattern at control`)
	issue.Hard(MetaCodeSyntax, `invalid meta-code syntax`)
	issue.Hard(ControlCodeSyntax, `invalid control-code syntax`)
	issue.Hard(CharClassValueAtEndOfRange, `char-class value at end of range`)
	issue.Hard(CharClassValueAtStartOfRange, `char-class value at start of range`)
	issue.Hard(UnmatchedRangeSpecifierInCharClass, `unmatched range specifier in char-class`)
	issue.Hard(TargetOfRepeatOperatorNotSpecified, `target of repeat operator is not specified`)
	issue.Hard(TargetOfRepeatOperatorInvalid, `target of repeat operator is invalid`)
	issue.Hard(NestedRepeatNotAllowed, `nested repeat is not allowed`)
	issue.Hard(NestedRepeatOperator, `nested repeat operator`)
	issue.Hard(UnmatchedCloseParenthesis, `unmatched close parenthesis`)
	issue.Hard(EndPatternWithUnmatchedParenthesis, `end pattern with unmatched parenthesis`)
	issue.Hard(EndPatternInGroup, `end pattern in group`)
	issue.Hard(UndefinedGroupOption, `undefined group option`)
	issue.Hard(InvalidPosixBracketType, `invalid POSIX bracket type`)
	issue.Hard(InvalidLookBehindPattern, `invalid pattern in look-behind`)
	issue.Hard(InvalidRepeatRangePattern, `invalid repeat range {lower,upper}`)
	issue.Hard(InvalidConditionPattern, `invalid conditional pattern`)

	/* values error (syntax error) */
	issue.Hard(TooBigNumber, `too big number`)
	issue.Hard(TooBigNumberForRepeatRange, `too big number for repeat range`)
	issue.Hard(UpperSmallerThanLowerInRepeatRange, `upper is smaller than lower in repeat range`)
	issue.Hard(EmptyRangeInCharClass, `empty range in char class`)
	issue.Hard(MismatchCodeLengthInClassRange, `mismatch multibyte code length in char-class range`)
	issue.Hard(TooManyMultiByteRanges, `too many multibyte code ranges are specified`)
	issue.Hard(TooShortMultiByteString, `too short multibyte code string`)
	issue.Hard(TooBigBackrefNumber, `too big backref number`)
	if //noinspection GoBoolExpressions
	config.UseNamedGroup {
		issue.Hard(InvalidBackref, `invalid backref number/name`)
		issue.Hard(InvalidCharInGroupName, `invalid char in group name <%{n}>`)
	} else {
		issue.Hard(InvalidBackref, `invalid backref number`)
		issue.Hard(InvalidCharInGroupName, `invalid char in group number <%{n}>`)
	}
	issue.Hard(NumberedBackrefOrCallNotAllowed, `numbered backref/call is not allowed. (use name)`)
	issue.Hard(TooShortDigits, `too short digits`)
	issue.Hard(InvalidWideCharValue, `invalid wide-char value`)
	issue.Hard(EmptyGroupName, `group name is empty`)
	issue.Hard(InvalidGroupName, `invalid group name <%n>`)
	issue.Hard(UndefinedNameReference, `undefined name <%n> reference`)
	issue.Hard(UndefinedGroupReference, `undefined group <%n> reference`)
	issue.Hard(MultiplexDefinedName, `multiplex defined name <%n>`)
	issue.Hard(MultiplexDefinitionNameCall, `multiplex definition name <%{n}> call`)
	issue.Hard(NeverEndingRecursion, `never ending recursion`)
	issue.Hard(GroupNumberOverForCaptureHistory, `group number is too big for capture history`)
	issue.Hard(NotSupportedEncodingCombination, `not supported encoding combination`)
	issue.Hard(InvalidCombinationOfOptions, `invalid combination of options`)
	issue.Hard(OverThreadPassLimitCount, `over thread pass limit count`)
	issue.Hard(TooBigSbCharValue, `too big singlebyte char value`)

	issue.Hard(CCTypeBug, `undefined type (bug)`)

	issue.Hard(CCTooBigWideCharValue, `too big wide-char value`)
	issue.Hard(CCTooLongWideCharValue, `too long wide-char value`)

	issue.Hard(CCInvalidCharPropertyName, `invalid character property name <%n>`)
	issue.Hard(CCInvalidCodePointValue, `invalid code point value`)

	issue.Hard(CCEncodingClassDefNotFound, `encoding class <%n> not found`)
	issue.Hard(CCEncodingLoadError, `problem loading encoding <%n>`)

	issue.Hard(CCIllegalCharacter, `illegal character`)

	issue.Hard(CCEncodingAlreadyRegistered, `encoding already registerd <%n>`)
	issue.Hard(CCEncodingAliasAlreadyRegistered, `encoding alias already registerd <%n>`)
	issue.Hard(CCEncodingReplicaAlreadyRegistered, `encoding replica already registerd <%n>`)
	issue.Hard(CCNoSuchEncodng, `no such encoding <%n>`)
	issue.Hard(CCCouldNotReplicate, `could not replicate <%n> encoding`)

	// transcoder messages
	issue.Hard(CCTranscoderAlreadyRegistered, `transcoder from <%n> has been already registered`)
	issue.Hard(CCTranscoderClassDefNotFound, `transcoder class <%n> not found`)
	issue.Hard(CCTranscoderLoadError, `problem loading transcoder <%n>`)

}
