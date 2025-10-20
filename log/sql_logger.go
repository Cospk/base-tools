package log

import (
	"context"
	"fmt"
	"time"

	"errors"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	gormUtils "gorm.io/gorm/utils"
)

// 纳秒转毫秒的转换常量
const nanosecondsToMilliseconds = 1e6

// SqlLogger GORM的日志记录器实现,用于记录SQL执行日志
type SqlLogger struct {
	LogLevel                  gormLogger.LogLevel // 日志级别
	IgnoreRecordNotFoundError bool                // 是否忽略记录未找到错误
	SlowThreshold             time.Duration       // 慢SQL阈值
}

// NewSqlLogger 创建一个新的SQL日志记录器
func NewSqlLogger(logLevel gormLogger.LogLevel, ignoreRecordNotFoundError bool, slowThreshold time.Duration) *SqlLogger {
	return &SqlLogger{
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
		SlowThreshold:             slowThreshold,
	}
}

// LogMode 设置日志级别并返回新的日志记录器实例
func (l *SqlLogger) LogMode(logLevel gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = logLevel
	return &newLogger
}

// Info 记录信息级别日志
func (SqlLogger) Info(ctx context.Context, msg string, args ...any) {
	ZInfo(ctx, msg, "args", args)
}

// Warn 记录警告级别日志
func (SqlLogger) Warn(ctx context.Context, msg string, args ...any) {
	ZWarn(ctx, msg, nil, "args", args)
}

// Error 记录错误级别日志
func (SqlLogger) Error(ctx context.Context, msg string, args ...any) {
	var err error = nil
	kvList := make([]any, 0)
	v, ok := args[0].(error)
	if ok {
		err = v
		for i := 1; i < len(args); i++ {
			kvList = append(kvList, fmt.Sprintf("args[%v]", i), args[i])
		}
	} else {
		for i := 0; i < len(args); i++ {
			kvList = append(kvList, fmt.Sprintf("args[%v]", i), args[i])
		}
	}
	ZError(ctx, msg, err, kvList...)
}

// Trace 追踪SQL执行,记录SQL语句、执行时间、影响行数等信息
func (l *SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			ZError(ctx, "sql exec detail", err, "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/nanosecondsToMilliseconds), "sql", sql)
		} else {
			ZError(ctx, "sql exec detail", err, "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/nanosecondsToMilliseconds), "rows", rows, "sql", sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			ZWarn(ctx, "sql exec detail", nil, "gorm", gormUtils.FileWithLineNum(), "slow sql", slowLog, "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/nanosecondsToMilliseconds), "sql", sql)
		} else {
			ZWarn(ctx, "sql exec detail", nil, "gorm", gormUtils.FileWithLineNum(), "slow sql", slowLog, "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/nanosecondsToMilliseconds), "rows", rows, "sql", sql)
		}
	case l.LogLevel == gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			ZDebug(ctx, "sql exec detail", "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/nanosecondsToMilliseconds), "sql", sql)
		} else {
			ZDebug(ctx, "sql exec detail", "gorm", gormUtils.FileWithLineNum(), "elapsed time", fmt.Sprintf("%f(ms)", float64(elapsed.Nanoseconds())/nanosecondsToMilliseconds), "rows", rows, "sql", sql)
		}
	}
}
