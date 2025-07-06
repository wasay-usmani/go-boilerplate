package errorx

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode is a type for internal error codes
// You can expand this as needed
// Example: ErrNotFound, ErrValidation, etc.
type ErrorCode string

const (
	Unknown         ErrorCode = "unknown"
	NotFound        ErrorCode = "not_found"
	Unauthorized    ErrorCode = "unauthorized"
	Internal        ErrorCode = "internal"
	AlreadyExists   ErrorCode = "already_exists"
	Forbidden       ErrorCode = "forbidden"
	Conflict        ErrorCode = "conflict"
	Timeout         ErrorCode = "timeout"
	TooManyRequests ErrorCode = "too_many_requests"
	BadRequest      ErrorCode = "bad_request"
)

// Error represents an internal error with code, message, internal code, and additional fields
// Fields can be used for validation errors, etc.
type Error struct {
	Code         ErrorCode      `json:"code"`
	Message      string         `json:"message"`
	InternalCode string         `json:"internal_code,omitempty"`
	Fields       map[string]any `json:"fields,omitempty"`
}

// NewValidation creates a validation error with fields
func NewValidation(message string, fields map[string]any) *Error {
	return New(BadRequest, message, WithFields(fields))
}

// Error implements the error interface
func (e *Error) Error() string {
	if len(e.Fields) > 0 {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Fields)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Option is a functional option for Error
// Example: errorx.New(errorx.BadRequest, "msg", errorx.WithInternalCode("123"), errorx.WithFields(fields))
type Option func(*Error)

// WithInternalCode sets the InternalCode field
func WithInternalCode(code string) Option {
	return func(e *Error) {
		e.InternalCode = code
	}
}

// WithFields sets the Fields field
func WithFields(fields map[string]any) Option {
	return func(e *Error) {
		e.Fields = fields
	}
}

// New creates a new Error with options
func New(code ErrorCode, message string, opts ...Option) *Error {
	err := &Error{
		Code:    code,
		Message: message,
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

// Is checks if the error matches the given ErrorCode
func Is(err error, code ErrorCode) bool {
	e, ok := err.(*Error)
	return ok && e.Code == code
}

// ToHTTPStatus converts an internal error code to an HTTP status code
func ToHTTPStatus(code ErrorCode) int {
	switch code {
	case NotFound:
		return http.StatusNotFound
	case BadRequest:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case AlreadyExists:
		return http.StatusConflict
	case Conflict:
		return http.StatusConflict
	case Timeout:
		return http.StatusRequestTimeout
	case TooManyRequests:
		return http.StatusTooManyRequests
	case Internal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ToGRPCCode converts an internal error code to a gRPC code
func ToGRPCCode(code ErrorCode) codes.Code {
	switch code {
	case NotFound:
		return codes.NotFound
	case BadRequest:
		return codes.InvalidArgument
	case Unauthorized:
		return codes.Unauthenticated
	case Forbidden:
		return codes.PermissionDenied
	case AlreadyExists:
		return codes.AlreadyExists
	case Conflict:
		return codes.Aborted
	case Timeout:
		return codes.DeadlineExceeded
	case TooManyRequests:
		return codes.ResourceExhausted
	case Internal:
		return codes.Internal
	default:
		return codes.Unknown
	}
}

// ToGRPCStatus converts an Error to a gRPC status error
func (e *Error) ToGRPCStatus() error {
	return status.Error(ToGRPCCode(e.Code), e.Error())
}

// FromGRPCStatus converts a gRPC status error to *Error
func FromGRPCStatus(err error) *Error {
	if err == nil {
		return nil
	}

	// Check if it's already our Error type
	if e, ok := err.(*Error); ok {
		return e
	}

	// Try to convert from gRPC status
	if st, ok := status.FromError(err); ok {
		return FromGRPCCode(st.Code(), st.Message())
	}

	// Fallback to unknown error
	return New(Unknown, err.Error())
}

// parseGRPCMessage extracts the actual message from a gRPC status message
// The message format from ToGRPCStatus is: "code: message (fields)"
// We want to extract just the "message" part
func parseGRPCMessage(message string) string {
	if message == "" {
		return message
	}

	// Look for the pattern "code: message" or "code: message (fields)"
	for i := 0; i < len(message); i++ {
		if message[i] != ':' {
			continue
		}
		// Found the colon, skip it and any spaces
		start := i + 1
		for start < len(message) && message[start] == ' ' {
			start++
		}

		// Check if there are fields at the end "(fields)"
		end := len(message)
		if start < len(message) && message[len(message)-1] == ')' {
			// Look for the opening parenthesis
			for j := len(message) - 2; j >= start; j-- {
				if message[j] == '(' {
					end = j
					break
				}
			}
			// Remove trailing space before the parenthesis
			for end > start && message[end-1] == ' ' {
				end--
			}
		}

		if start < end {
			return message[start:end]
		}
		break
	}

	return message
}

// FromGRPCCode converts a gRPC code and message to *Error
func FromGRPCCode(code codes.Code, message string) *Error {
	errorCode := FromGRPCCodeToErrorCode(code)
	actualMessage := parseGRPCMessage(message)
	return New(errorCode, actualMessage)
}

// FromGRPCCodeToErrorCode converts a gRPC code to ErrorCode
func FromGRPCCodeToErrorCode(code codes.Code) ErrorCode {
	switch code {
	case codes.NotFound:
		return NotFound
	case codes.InvalidArgument:
		return BadRequest
	case codes.Unauthenticated:
		return Unauthorized
	case codes.PermissionDenied:
		return Forbidden
	case codes.AlreadyExists:
		return AlreadyExists
	case codes.Aborted:
		return Conflict
	case codes.DeadlineExceeded:
		return Timeout
	case codes.ResourceExhausted:
		return TooManyRequests
	case codes.Internal:
		return Internal
	case codes.Unknown:
		return Unknown
	default:
		return Internal
	}
}

// FromError attempts to cast a generic error to *Error
func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return e
	}

	return New(Unknown, err.Error(), nil)
}

// WithMessage creates a new Error with the given message
func (c ErrorCode) WithMessage(msg string) *Error {
	return &Error{
		Code:    c,
		Message: msg,
		Fields:  nil,
	}
}

// WithMessagef creates a new Error with a formatted message
func (c ErrorCode) WithMessagef(format string, args ...any) *Error {
	return &Error{
		Code:    c,
		Message: fmt.Sprintf(format, args...),
		Fields:  nil,
	}
}

// FromGRPCWithDetails converts a gRPC error with details to *Error
// This is useful when the gRPC error contains additional structured data
func FromGRPCWithDetails(err error, details map[string]any) *Error {
	if err == nil {
		return nil
	}

	// Convert the base error
	baseError := FromGRPCStatus(err)

	// Add the details as fields
	if len(details) > 0 {
		baseError.Fields = details
	}

	return baseError
}
