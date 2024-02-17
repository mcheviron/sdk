// Package logger provides a wrapper around slog that emulates zap by avoiding heap allocations.
package logger

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

// Logger represents a logger instance.
type Logger struct {
	innerLogger *slog.Logger
	context     context.Context
}

func New(env string) (Logger, error) {
	ctx := context.Background()

	switch env {
	case "prod":
		return Logger{
			innerLogger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.SourceKey {
						source := a.Value.Any().(*slog.Source)
						source.File = filepath.Base(source.File)
					}

					// used in AWS or GCP, they'll add the timestamp themselves
					if a.Key == slog.TimeKey {
						return slog.Attr{}
					}
					return a
				},
			})),
			context: ctx,
		}, nil
	case "dev":
		return Logger{
			innerLogger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.SourceKey {
						source := a.Value.Any().(*slog.Source)
						source.File = filepath.Base(source.File)
					}

					if a.Key == slog.TimeKey {
						return slog.Attr{}
					}
					return a
				},
			})),
			context: ctx,
		}, nil
	default:
		return Logger{}, errors.New("unknown environment")
	}
}

func (l Logger) Info(msg string, attrs ...slog.Attr) {
	l.innerLogger.LogAttrs(l.context, slog.LevelInfo, msg, attrs...)
}

func (l Logger) Error(msg string, attrs ...slog.Attr) {
	l.innerLogger.LogAttrs(l.context, slog.LevelError, msg, attrs...)
}

func (l Logger) Debug(msg string, attrs ...slog.Attr) {
	l.innerLogger.LogAttrs(l.context, slog.LevelDebug, msg, attrs...)
}

func (l Logger) Warn(msg string, attrs ...slog.Attr) {
	l.innerLogger.LogAttrs(l.context, slog.LevelWarn, msg, attrs...)
}
