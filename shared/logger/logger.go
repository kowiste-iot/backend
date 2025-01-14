package logger

import "context"

type Logger interface {
	Info(ctx context.Context, msg string, fields any)
	Error(ctx context.Context, err error, msg string, fields any)
	Debug(ctx context.Context, msg string, fields any)
	Warn(ctx context.Context, msg string, fields any)
}
