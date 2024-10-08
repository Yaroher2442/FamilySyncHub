package logger

import (
	"context"

	"go.uber.org/zap"
)

func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}

func WithFields(fields ...zap.Field) *zap.Logger {
	return globalLogger.With(fields...)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	globalLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	globalLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func DebugContext(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, WithCtxFields(ctx, fields...)...)
}

func InfoContext(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(msg, WithCtxFields(ctx, fields...)...)
}

func WarnContext(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, WithCtxFields(ctx, fields...)...)
}

func ErrorContext(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(msg, WithCtxFields(ctx, fields...)...)
}

func DPanicContext(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.DPanic(msg, WithCtxFields(ctx, fields...)...)
}

func PanicContext(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Panic(msg, WithCtxFields(ctx, fields...)...)
}

func FatalContext(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, WithCtxFields(ctx, fields...)...)
}
