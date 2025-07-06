package logkit

import "context"

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
	Fatal LogLevel = "fatal"
)

type Logger interface {
	// Base returns the underlying logger
	Base() any

	// Debug writes a debug log
	// args is a list of key-value pairs where key should always be a string.
	Debug(msg string, args ...any)

	// Info writes an info log
	// args is a list of key-value pairs where key should always be a string.
	Info(msg string, args ...any)

	// Warn writes a warn log
	// args is a list of key-value pairs where key should always be a string.
	Warn(msg string, args ...any)

	// Error writes an error log. If err is not nil, the error will be logged as well
	// args is a list of key-value pairs where key should always be a string.
	Error(msg string, err error, args ...any)

	// ErrorStack writes an error log. It will add the stack trace to the error
	// args is a list of key-value pairs where key should always be a string.
	ErrorStack(msg string, err error, args ...any)

	// Error writes an error log. If err is not nil, the error will be logged as well
	// args is a list of key-value pairs where key should always be a string.
	Fatal(msg string, err error, args ...any)

	// Error writes an error log. It will add the stack trace to the error
	// args is a list of key-value pairs where key should always be a string.
	FatalStack(msg string, err error, args ...any)

	// With returns a new logger with contextual fields
	// args is a list of key-value pairs where key should always be a string.
	With(args ...any) Logger

	// WithContext returns a new logger with contextual fields and context.Context
	// This is useful for adding trace id and span id
	// args is a list of key-value pairs where key should always be a string.
	WithContext(ctx context.Context, args ...any) Logger
}
