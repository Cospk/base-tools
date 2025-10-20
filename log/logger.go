// Package log 提供基于Zap的高性能结构化日志系统，支持日志轮转、多输出和上下文感知
package log

import "context"

// Logger 日志记录器接口，提供不同级别的日志方法和上下文管理
type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...any)
	Info(ctx context.Context, msg string, keysAndValues ...any)
	Warn(ctx context.Context, msg string, err error, keysAndValues ...any)
	Error(ctx context.Context, msg string, err error, keysAndValues ...any)
	Panic(ctx context.Context, msg string, err error, keysAndValues ...any)
	WithValues(keysAndValues ...any) Logger
	WithName(name string) Logger
	WithCallDepth(depth int) Logger
	Flush()
}
