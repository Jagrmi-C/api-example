package errorx

import "fmt"

// ErrorsRespo defines model for errorRespo.
type ErrorsRespo struct {
	Data   *any          `json:"data"`
	Status int           `json:"-"`
	Errors []ErrorEntity `json:"errors"`
}

func (e *ErrorsRespo) Error() string {
	var msg string

	for _, err := range e.Errors {
		msg += fmt.Sprintf("%s\n", err.Detail)
	}

	return msg
}

func (e *ErrorsRespo) GetStatus() int {
	return e.Status
}

// ErrorEntity defines model for errorEntity.
type ErrorEntity struct {
	ID     string  `json:"id"`
	Code   string  `json:"code,omitempty"`
	Detail string  `json:"detail,omitempty"`
	Source *Source `json:"source,omitempty"`
	Value  any     `json:"value,omitempty"`
	Status int     `json:"status"`
	Title  string  `json:"title"`
}

func (e *ErrorEntity) Error() string {
	return e.Title
}

// Source represents the source of the error entity.
type Source struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
	Header    string `json:"header,omitempty"`
	Path      string `json:"path,omitempty"`
}

type ErrorCategory int

const (
	InvalidInput ErrorCategory = iota
	Unexpected
	BusinessLogicError
	Unauthorized
	DatabaseError
	HTTPError
	InternalError
)

func (c ErrorCategory) String() string {
	switch c {
	case InvalidInput:
		return "Invalid input."
	case BusinessLogicError:
		return "Business logic error."
	case Unauthorized:
		return "Unauthorized access."
	case DatabaseError:
		return "Database error."
	case HTTPError:
		return "HTTP error."
	default:
		return "internal error category."
	}
}

// DomainError represents a custom domain error-entity that
// should be caught detailed information about issues.
type DomainError interface {
	error
	Category() ErrorCategory
	Code() string
	Message() string
	Unwrap() error
	Details() string
	Debug() any
}

// type StatusError interface {
// 	GetStatus() int
// 	Error() string
// }

// Error represents an error that could be wrapping another error, it includes a code for determining what
// triggered the error.
type Error struct {
	orig     error
	msg      string
	code     ErrorCode
	category ErrorCategory
	details  string
	source   string
	debug    any
}

// Error returns the message, when wrapping errors the wrapped error is returned.
func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.orig)
	}

	return e.msg
}

// WrapMsg allows to add additional information to the message.
// usecase: get retailers handler: get retailers usecase: impossible to retrieve data from database.
func (e *Error) WrapMsg(msg string) {
	if e.msg == "" {
		e.msg = msg
		return
	}

	e.msg = fmt.Sprintf("%s: %s", msg, e.msg)
}

// Unwrap returns the wrapped error, if any.
func (e *Error) Unwrap() error {
	return e.orig
}

// Code returns the code representing this error.
func (e *Error) Code() ErrorCode {
	return e.code
}

// Message returns the message associated with this error.
func (e *Error) Message() string {
	return e.msg
}

// Category returns the code representing an error's category.
func (e *Error) Category() string {
	return e.category.String()
}

// Source is a getter.
func (e *Error) Source() string {
	return e.source
}

// Details is a getter.
func (e *Error) Details() string {
	return fmt.Sprintf("%v", e.details)
}

// Debug returns a debug data for the logging context.
func (e *Error) Debug() any {
	return e.debug
}

// NewErrorf instantiates a new error without the original error.
func NewErrorf(code ErrorCode, format string, a ...any) error {
	return WrapErrorf(nil, code, format, a...)
}

// WrapErrorf returns a wrapped error around original with context message.
func WrapErrorf(orig error, code ErrorCode, format string, a ...any) error {
	return &Error{
		code: code,
		orig: orig,
		msg:  fmt.Sprintf(format, a...),
	}
}

// WrapErrorfWithDebug returns a wrapped error around original with debug data.
func WrapErrorfWithDebug(orig error, debug any, code ErrorCode, format string, a ...any) error {
	return &Error{
		code:  code,
		orig:  orig,
		msg:   fmt.Sprintf(format, a...),
		debug: debug,
	}
}

// WrapDetailerErrorf returns a wrapped error around original with source of error.
func WrapDetailerErrorf(orig error, source, details string, code ErrorCode, format string, a ...any) error {
	return &Error{
		code:    code,
		orig:    orig,
		msg:     fmt.Sprintf(format, a...),
		source:  source,
		details: details,
	}
}
