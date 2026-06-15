package errors

// ErrorCode indicates what type of error return to user.
//
// [CodeDirect] indiactes that error can be returned directly.
type ErrorCode int

const (
	CodeDirect ErrorCode = iota
	CodeUnfixable
	CodeInternal
	CodeNotFound
	CodeAlreadyExists
	CodeMarshalError
	CodeUnmarshalError
	CodeCacheExpired
	CodeInvalidArgument
)

var CodeMap = map[ErrorCode]string{
	1: "unfixable",
	2: "internal error",
	3: "not found",
	4: "already exists",
	5: "marshal error",
	6: "unmarshal error",
	7: "cache expired",
	8: "invalid argument",
}
