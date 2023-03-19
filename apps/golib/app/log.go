package app

import (
	"context"
)

// LogDebug logs the given message at debug level
func LogDebug(ctx context.Context, msg string, fields ...interface{}) {
	l := zapFromContext(ctx)
	if l == nil {
		return
	}

	// TODO: Add OTEL stuff

	l.Debugw(msg, fields...)
}

// LogInfo logs the given message at info level
func LogInfo(ctx context.Context, msg string, fields ...interface{}) {
	l := zapFromContext(ctx)
	if l == nil {
		return
	}

	// TODO: Add OTEL stuff

	l.Infow(msg, fields...)
}

// LogWarn logs the given message at warn level
func LogWarn(ctx context.Context, msg string, fields ...interface{}) {
	l := zapFromContext(ctx)
	if l == nil {
		return
	}

	// TODO: Add OTEL stuff

	l.Warnw(msg, fields...)
}

// LogError logs the given message at error level and also reports the error
func LogError(ctx context.Context, err error, fields ...interface{}) {
	l := zapFromContext(ctx)
	if l == nil {
		return
	}

	// TODO: Add OTEL stuff

	// TODO: Add error tracking stuff

	l.Errorw(err.Error(), fields...)
}
