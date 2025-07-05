package errorx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFromGRPCStatus(t *testing.T) {
	tests := []struct {
		name     string
		grpcErr  error
		expected *Error
	}{
		{
			name:     "nil error",
			grpcErr:  nil,
			expected: nil,
		},
		{
			name:    "not found error",
			grpcErr: status.Error(codes.NotFound, "resource not found"),
			expected: &Error{
				Code:    NotFound,
				Message: "resource not found",
			},
		},
		{
			name:    "invalid argument error",
			grpcErr: status.Error(codes.InvalidArgument, "invalid input"),
			expected: &Error{
				Code:    BadRequest,
				Message: "invalid input",
			},
		},
		{
			name:    "unauthenticated error",
			grpcErr: status.Error(codes.Unauthenticated, "authentication required"),
			expected: &Error{
				Code:    Unauthorized,
				Message: "authentication required",
			},
		},
		{
			name:    "permission denied error",
			grpcErr: status.Error(codes.PermissionDenied, "access denied"),
			expected: &Error{
				Code:    Forbidden,
				Message: "access denied",
			},
		},
		{
			name:    "already exists error",
			grpcErr: status.Error(codes.AlreadyExists, "resource already exists"),
			expected: &Error{
				Code:    AlreadyExists,
				Message: "resource already exists",
			},
		},
		{
			name:    "internal error",
			grpcErr: status.Error(codes.Internal, "internal server error"),
			expected: &Error{
				Code:    Internal,
				Message: "internal server error",
			},
		},
		{
			name:    "unknown error",
			grpcErr: status.Error(codes.Unknown, "unknown error"),
			expected: &Error{
				Code:    Unknown,
				Message: "unknown error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromGRPCStatus(tt.grpcErr)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected.Code, result.Code)
				assert.Equal(t, tt.expected.Message, result.Message)
			}
		})
	}
}

func TestFromGRPCCode(t *testing.T) {
	tests := []struct {
		name     string
		code     codes.Code
		message  string
		expected *Error
	}{
		{
			name:    "not found",
			code:    codes.NotFound,
			message: "not found",
			expected: &Error{
				Code:    NotFound,
				Message: "not found",
			},
		},
		{
			name:    "invalid argument",
			code:    codes.InvalidArgument,
			message: "bad request",
			expected: &Error{
				Code:    BadRequest,
				Message: "bad request",
			},
		},
		{
			name:    "timeout",
			code:    codes.DeadlineExceeded,
			message: "request timeout",
			expected: &Error{
				Code:    Timeout,
				Message: "request timeout",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromGRPCCode(tt.code, tt.message)
			assert.Equal(t, tt.expected.Code, result.Code)
			assert.Equal(t, tt.expected.Message, result.Message)
		})
	}
}

func TestFromGRPCCodeToErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		grpcCode codes.Code
		expected ErrorCode
	}{
		{"not found", codes.NotFound, NotFound},
		{"invalid argument", codes.InvalidArgument, BadRequest},
		{"unauthenticated", codes.Unauthenticated, Unauthorized},
		{"permission denied", codes.PermissionDenied, Forbidden},
		{"already exists", codes.AlreadyExists, AlreadyExists},
		{"aborted", codes.Aborted, Conflict},
		{"deadline exceeded", codes.DeadlineExceeded, Timeout},
		{"resource exhausted", codes.ResourceExhausted, TooManyRequests},
		{"internal", codes.Internal, Internal},
		{"unknown", codes.Unknown, Unknown},
		{"unimplemented", codes.Unimplemented, Internal}, // default case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromGRPCCodeToErrorCode(tt.grpcCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsGRPCError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"grpc error", status.Error(codes.NotFound, "not found"), true},
		{"regular error", assert.AnError, false},
		{"our error", New(NotFound, "not found"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsGRPCError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetGRPCCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected codes.Code
	}{
		{"nil error", nil, codes.OK},
		{"grpc not found", status.Error(codes.NotFound, "not found"), codes.NotFound},
		{"grpc invalid argument", status.Error(codes.InvalidArgument, "bad request"), codes.InvalidArgument},
		{"regular error", assert.AnError, codes.Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetGRPCCode(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetGRPCMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"nil error", nil, ""},
		{"grpc error", status.Error(codes.NotFound, "resource not found"), "resource not found"},
		{"regular error", assert.AnError, assert.AnError.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetGRPCMessage(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFromGRPCWithDetails(t *testing.T) {
	grpcErr := status.Error(codes.InvalidArgument, "validation failed")
	details := map[string]any{
		"field":  "email",
		"reason": "invalid format",
		"value":  "invalid-email",
	}

	result := FromGRPCWithDetails(grpcErr, details)

	assert.Equal(t, BadRequest, result.Code)
	assert.Equal(t, "validation failed", result.Message)
	assert.Equal(t, details, result.Fields)
}

func TestErrorToGRPCStatusRoundTrip(t *testing.T) {
	originalError := New(BadRequest, "invalid input", WithFields(map[string]any{
		"field": "email",
	}))

	// Convert to gRPC status
	grpcErr := originalError.ToGRPCStatus()

	// Convert back from gRPC status
	convertedError := FromGRPCStatus(grpcErr)

	assert.Equal(t, originalError.Code, convertedError.Code)
	assert.Equal(t, originalError.Message, convertedError.Message)
	// Note: Fields are lost in the round trip because gRPC status doesn't preserve them
	// This is expected behavior
}

func TestFromGRPCCodeMessageParsing(t *testing.T) {
	tests := []struct {
		name     string
		code     codes.Code
		message  string
		expected string
	}{
		{
			name:     "simple message",
			code:     codes.InvalidArgument,
			message:  "bad request",
			expected: "bad request",
		},
		{
			name:     "message with code prefix",
			code:     codes.InvalidArgument,
			message:  "bad_request: invalid input",
			expected: "invalid input",
		},
		{
			name:     "message with code prefix and fields",
			code:     codes.InvalidArgument,
			message:  "bad_request: invalid input (map[field:email])",
			expected: "invalid input",
		},
		{
			name:     "message with code prefix and complex fields",
			code:     codes.InvalidArgument,
			message:  "bad_request: validation failed (map[field:email reason:invalid_format])",
			expected: "validation failed",
		},
		{
			name:     "message without colon",
			code:     codes.NotFound,
			message:  "resource not found",
			expected: "resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromGRPCCode(tt.code, tt.message)
			assert.Equal(t, tt.expected, result.Message)
		})
	}
}
