package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps the zap SugaredLogger for structured logging.
type Logger struct {
	*zap.SugaredLogger
	base *zap.Logger
}

// New creates a new Logger instance.
// In development mode, it uses colored console output.
// In production, it uses JSON structured logging.
func New(env string, debug bool) *Logger {
	var zapLogger *zap.Logger

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if env == "development" || debug {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
		zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		core := zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
		zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
		base:          zapLogger,
	}
}

// WithRequestID returns a new logger with the request ID field.
func (l *Logger) WithRequestID(requestID string) *Logger {
	newLogger := l.base.With(zap.String("request_id", requestID))
	return &Logger{
		SugaredLogger: newLogger.Sugar(),
		base:          newLogger,
	}
}

// WithTraceID returns a new logger with the trace ID field.
func (l *Logger) WithTraceID(traceID string) *Logger {
	newLogger := l.base.With(zap.String("trace_id", traceID))
	return &Logger{
		SugaredLogger: newLogger.Sugar(),
		base:          newLogger,
	}
}

// WithFields returns a new logger with the given fields.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	newLogger := l.base.With(zapFields...)
	return &Logger{
		SugaredLogger: newLogger.Sugar(),
		base:          newLogger,
	}
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.base.Sync()
}

// Base returns the underlying zap.Logger for use with middleware.
func (l *Logger) Base() *zap.Logger {
	return l.base
}
