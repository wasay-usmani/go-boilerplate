package logkit

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

func TestNewZeroLogger(t *testing.T) {
	logger := NewLogger(Info, "gotils-test")
	if logger == nil {
		t.Error("Expected non-nil logger from NewZeroLogger")
	}
}

func TestLogger_Base(t *testing.T) {
	logger := NewLogger(Info, "gotils-test", WithStackTrace(), WithTraceHook("span-id", "trace-id"))
	if logger.Base() == nil {
		t.Errorf("Expected non-nil Base() logger")
	}
}

func TestLogger_Debug(t *testing.T) {
	buf := &bytes.Buffer{}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := NewLogger(Debug, "gotils-test", WithOutput(buf))

	logger.Debug("gotils log initialized", "component", "logger", "status", "active")

	out := buf.String()
	if !strings.Contains(out, "gotils log initialized") ||
		!strings.Contains(out, `"component":"logger"`) ||
		!strings.Contains(out, `"status":"active"`) {
		t.Errorf("Expected debug log with message and fields, got: %s", out)
	}
}

func TestLogger_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(Info, "gotils-test", WithOutput(buf))

	logger.Info("info message", "version", "1.0.0")

	out := buf.String()
	if !strings.Contains(out, "info message") || !strings.Contains(out, `"version":"1.0.0"`) {
		t.Errorf("Expected info log with message and field, got: %s", out)
	}
}

func TestLogger_Warn(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(Warn, "gotils-test", WithOutput(buf))

	logger.Warn("warning issued", "threshold", "90%")

	out := buf.String()
	if !strings.Contains(out, "warning issued") || !strings.Contains(out, `"threshold":"90%"`) {
		t.Errorf("Expected warn log with message and field, got: %s", out)
	}
}

func TestLogger_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(Error, "gotils-test", WithOutput(buf))

	err := errors.New("network down")
	logger.Error("error occurred", err, "attempt", 3)

	out := buf.String()
	if !strings.Contains(out, "error occurred") || !strings.Contains(out, "network down") || !strings.Contains(out, `"attempt":3`) {
		t.Errorf("Expected error log with error and fields, got: %s", out)
	}
}

func TestLogger_With(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(Info, "gotils-test", WithOutput(buf))

	l := logger.With("request_id", "abc123")
	l.Info("with", "status", "ok")

	out := buf.String()
	if !strings.Contains(out, `"request_id":"abc123"`) || !strings.Contains(out, `"status":"ok"`) {
		t.Errorf("Expected With and Info log, got: %s", out)
	}
}

func TestLogger_WithContext(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(Info, "gotils-test", WithOutput(buf), WithTraceHook(spanFieldName, traceFieldName))

	ctx := context.Background()
	ctx = context.WithValue(ctx, spanFieldName, "span-123")   //nolint:staticcheck
	ctx = context.WithValue(ctx, traceFieldName, "trace-456") //nolint:staticcheck

	loggerWithCtx := logger.WithContext(ctx, "user", "anand", "env", "staging")
	loggerWithCtx.Info("contextual log")

	out := buf.String()

	if !strings.Contains(out, `"user":"anand"`) || !strings.Contains(out, `"env":"staging"`) {
		t.Errorf("Expected fields from args in log output, got: %s", out)
	}

	if !strings.Contains(out, `"span-id":"span-123"`) || !strings.Contains(out, `"trace-id":"trace-456"`) {
		t.Errorf("Expected context values in log output, got: %s", out)
	}

	if !strings.Contains(out, "contextual log") {
		t.Errorf("Expected log message to be present, got: %s", out)
	}
}
