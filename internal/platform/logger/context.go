package logger

import (
	"context"
	"strings"

	"go.uber.org/zap"
)

type loggerContextKey struct{}
type requestIDContextKey struct{}

func WithContext(ctx context.Context, log *zap.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if log == nil {
		return ctx
	}

	return context.WithValue(ctx, loggerContextKey{}, log)
}

func FromContext(ctx context.Context) *zap.Logger {
	log, ok := contextValue[*zap.Logger](ctx, loggerContextKey{})
	if ok && log != nil {
		return log
	}

	return zap.L()
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return ctx
	}

	ctx = context.WithValue(ctx, requestIDContextKey{}, requestID)
	return WithContext(ctx, FromContext(ctx).With(zap.String("request_id", requestID)))
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, ok := contextValue[string](ctx, requestIDContextKey{})
	if ok {
		return requestID
	}

	return ""
}

func contextValue[T any](ctx context.Context, key any) (T, bool) {
	var zero T
	if ctx == nil {
		return zero, false
	}

	value, ok := ctx.Value(key).(T)
	if !ok {
		return zero, false
	}

	return value, true
}
