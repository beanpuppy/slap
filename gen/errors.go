//soppo:generated v1
// Error types for slap CLI parser.
// ErrorKind represents the type of parse error.
package slap

/*soppo:enum
ErrorKind {
    MissingRequired
    InvalidValue
    UnknownFlag
    UnknownArg
    TooManyArgs
    ValidationFailed
    HelpRequested
    VersionRequested
    UnknownSubcommand
}
*/
type ErrorKind interface {
	isErrorKind()
}

type ErrorKind_MissingRequired struct {}
func (ErrorKind_MissingRequired) isErrorKind() {}
func (ErrorKind_MissingRequired) String() string { return "MissingRequired" }

type ErrorKind_InvalidValue struct {}
func (ErrorKind_InvalidValue) isErrorKind() {}
func (ErrorKind_InvalidValue) String() string { return "InvalidValue" }

type ErrorKind_UnknownFlag struct {}
func (ErrorKind_UnknownFlag) isErrorKind() {}
func (ErrorKind_UnknownFlag) String() string { return "UnknownFlag" }

type ErrorKind_UnknownArg struct {}
func (ErrorKind_UnknownArg) isErrorKind() {}
func (ErrorKind_UnknownArg) String() string { return "UnknownArg" }

type ErrorKind_TooManyArgs struct {}
func (ErrorKind_TooManyArgs) isErrorKind() {}
func (ErrorKind_TooManyArgs) String() string { return "TooManyArgs" }

type ErrorKind_ValidationFailed struct {}
func (ErrorKind_ValidationFailed) isErrorKind() {}
func (ErrorKind_ValidationFailed) String() string { return "ValidationFailed" }

type ErrorKind_HelpRequested struct {}
func (ErrorKind_HelpRequested) isErrorKind() {}
func (ErrorKind_HelpRequested) String() string { return "HelpRequested" }

type ErrorKind_VersionRequested struct {}
func (ErrorKind_VersionRequested) isErrorKind() {}
func (ErrorKind_VersionRequested) String() string { return "VersionRequested" }

type ErrorKind_UnknownSubcommand struct {}
func (ErrorKind_UnknownSubcommand) isErrorKind() {}
func (ErrorKind_UnknownSubcommand) String() string { return "UnknownSubcommand" }

var (
	ErrorKindMissingRequired ErrorKind = ErrorKind_MissingRequired{}
	ErrorKindInvalidValue ErrorKind = ErrorKind_InvalidValue{}
	ErrorKindUnknownFlag ErrorKind = ErrorKind_UnknownFlag{}
	ErrorKindUnknownArg ErrorKind = ErrorKind_UnknownArg{}
	ErrorKindTooManyArgs ErrorKind = ErrorKind_TooManyArgs{}
	ErrorKindValidationFailed ErrorKind = ErrorKind_ValidationFailed{}
	ErrorKindHelpRequested ErrorKind = ErrorKind_HelpRequested{}
	ErrorKindVersionRequested ErrorKind = ErrorKind_VersionRequested{}
	ErrorKindUnknownSubcommand ErrorKind = ErrorKind_UnknownSubcommand{}
)

// ParseError represents a CLI parsing error.
type ParseError struct {
	Kind ErrorKind
	Field string
	Value string
	Message string
}

func (e ParseError) Error() string {
	return e.Message
}

// errorAs is a helper since we can't use errors.As directly with enums
func errorAs(err error, target *ParseError) bool {
	if pe, ok := err.(ParseError); ok {
		(*target) = pe
		return true
	}
	return false
}

