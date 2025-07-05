package logkit

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	spanFieldName    = "span-id"
	traceFieldName   = "trace-id"
	serviceFieldName = "service"
)

type LogOpt func(*zeroLogger)
type zeroLogger struct {
	spanKey, traceKey string

	l zerolog.Logger
}

// NewLogger returns a new logger, with the specified log level
// and optional options (WithStackTrace, WithTraceHook, WithOutput)
func NewLogger(logLevel LogLevel, svc string, options ...LogOpt) Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerSkipFrameCount = 3 // avoids logging gotils reference
	zerolog.SetGlobalLevel(logLevelMap[logLevel])
	zlog := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Str(serviceFieldName, svc).
		Logger()

	z := &zeroLogger{l: zlog}
	for _, opt := range options {
		opt(z)
	}

	return z
}

// WithStackTrace adds the stack trace to the error logs.
func WithStackTrace() LogOpt {
	return func(z *zeroLogger) {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}
}

// WithTraceHook adds the trace id and span id to the logs
func WithTraceHook(spanKey, traceKey string) LogOpt {
	return func(z *zeroLogger) {
		z.traceKey = traceKey
		z.spanKey = spanKey
		z.l = z.l.Hook(z)
	}
}

// WithOutput sets the output of the logger
func WithOutput(w io.Writer) LogOpt {
	return func(z *zeroLogger) {
		z.l = z.l.Output(w)
	}
}

func (z *zeroLogger) Base() any {
	return z.l
}

func (z *zeroLogger) Debug(msg string, args ...any) {
	z.l.Debug().Fields(args).Msg(msg)
}

func (z *zeroLogger) Info(msg string, args ...any) {
	z.l.Info().Fields(args).Msg(msg)
}

func (z *zeroLogger) Warn(msg string, args ...any) {
	z.l.Warn().Fields(args).Msg(msg)
}

func (z *zeroLogger) Error(msg string, err error, args ...any) {
	z.l.Error().Err(err).Fields(args).Msg(msg)
}

func (z *zeroLogger) ErrorStack(msg string, err error, args ...any) {
	z.l.Error().Stack().Err(errors.Wrap(err, msg)).Fields(args).Msg(msg)
}

func (z *zeroLogger) Fatal(msg string, err error, args ...any) {
	z.l.Fatal().Err(err).Fields(args).Msg(msg)
}

func (z *zeroLogger) FatalStack(msg string, err error, args ...any) {
	z.l.Fatal().Stack().Err(errors.Wrap(err, msg)).Fields(args).Msg(msg)
}

func (z *zeroLogger) With(args ...any) Logger {
	if len(args) == 0 {
		return z
	}

	c := z.clone()
	c.l = z.l.With().Fields(args).Logger()
	return c
}

func (z *zeroLogger) WithContext(ctx context.Context, args ...any) Logger {
	return &zeroLogger{
		traceKey: z.traceKey,
		spanKey:  z.spanKey,
		l:        z.l.With().Ctx(ctx).Fields(args).Logger(),
	}
}

var logLevelMap = map[LogLevel]zerolog.Level{
	Debug: zerolog.DebugLevel,
	Info:  zerolog.InfoLevel,
	Warn:  zerolog.WarnLevel,
	Error: zerolog.ErrorLevel,
	Fatal: zerolog.FatalLevel,
}

// Hook to add trace id and span id to logs
func (z *zeroLogger) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	spanID := ctx.Value(z.spanKey)
	s, ok := spanID.(string)
	if ok {
		e.Str(spanFieldName, s)
	}

	traceID := ctx.Value(z.traceKey)
	t, ok := traceID.(string)
	if ok {
		e.Str(traceFieldName, t)
	}
}

func (z *zeroLogger) clone() *zeroLogger {
	zl := *z
	return &zl
}
