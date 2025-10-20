package log

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestSDKLog tests the SDKLog function for proper log output including custom [file:line] information
func TestSDKLog(t *testing.T) {
	sdkType := "TestSDK"
	platformName := "testPlatform"

	err := InitLoggerFromConfig(
		"testLogger", // loggerPrefixName
		"testModule", // moduleName
		sdkType,      // sdkType
		platformName, // platformName
		5,            // logLevel (Debug)
		true,         // isStdout
		false,        // isJson
		// "./logs",     // logLocation
		".",     // logLocation
		5,       // rotateCount
		24,      // rotationTime
		"1.0.0", // moduleVersion
		false,   // isSimplify
	)
	assert.NoError(t, err)

	// var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := zap.NewExample()
	defer logger.Sync()

	ZDebug(context.Background(), "hello")
	SDKLog(context.Background(), 5, "cmd/abc.go", 666, "This is a test message", nil, []any{"key", "value"})
	SDKLog(context.Background(), 4, "cmd/abc.go", 666, "This is a test message", nil, []any{"key", "value", "key", "key", 1})
	SDKLog(context.Background(), 3, "cmd/abc.go", 666, "This is a test message", nil, []any{"key", "value"})
	SDKLog(context.Background(), 2, "cmd/abc.go", 666, "This is a test message", nil, []any{"key", "value"})
	ZWarn(context.TODO(), "msg", nil)
	ZInfo(context.TODO(), "msg", nil)
	ZDebug(context.TODO(), "msg")

	w.Close()
	out, _ := os.ReadFile(r.Name())
	os.Stdout = stdout

	_ = string(out)
	// assert.Contains(t, output, "This is a test message")
	// assert.Contains(t, output, "[TestSDK/TestPlatform]")
	// assert.Contains(t, output, "[test_file.go:123]")
	// assert.Contains(t, output, "key")
	// assert.Contains(t, output, "value")
}

func TestDefaultLog(t *testing.T) {
	er := errors.New("error")
	ZInfo(context.Background(), "Are you OK?")
	ZDebug(context.Background(), "Hello")
	ZWarn(context.Background(), "3Q", er)
	ZError(context.Background(), "3Q very much", er)

	sdkType := "TestSDK"
	platformName := "testPlatform"

	err := InitLoggerFromConfig(
		"testLogger", // loggerPrefixName
		"testModule", // moduleName
		sdkType,      // sdkType
		platformName, // platformName
		int(5),       // logLevel (Debug)
		true,         // isStdout
		false,        // isJson
		"./logs",     // logLocation
		uint(5),      // rotateCount
		uint(24),     // rotationTime
		"1.0.0",      // moduleVersion
		false,        // isSimplify
	)
	assert.NoError(t, err)

	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := zap.NewExample()
	defer logger.Sync()

	ZDebug(context.Background(), "hello")
	SDKLog(context.Background(), 5, "cmd/abc.go", 666, "This is a test message", nil, []any{"key", "value"})
	SDKLog(context.Background(), 4, "cmd/abc.go", 666, "This is a test message", nil, []any{"key", "value", "key", "key", 1})
	SDKLog(context.Background(), 3, "cmd/abc.go", 666, "This is a test message", nil, []any{"key", "value"})
	SDKLog(context.Background(), 2, "cmd/abc.go", 666, "This is a test message", nil, []any{"key", "value"})
	ZWarn(context.TODO(), "msg", nil)
	ZInfo(context.TODO(), "msg", nil)
	ZDebug(context.TODO(), "msg")

	w.Close()
	out, _ := os.ReadFile(r.Name())
	os.Stdout = stdout

	_ = string(out)
	// assert.Contains(t, output, "This is a test message")
	// assert.Contains(t, output, "[TestSDK/TestPlatform]")
	// assert.Contains(t, output, "[test_file.go:123]")
	// assert.Contains(t, output, "key")
	// assert.Contains(t, output, "value")
}

// ==================== ZPanic 函数测试 ====================

func TestZPanic(t *testing.T) {
	// 初始化测试logger
	err := InitConsoleLogger("test", LevelDebug, false, "1.0.0")
	assert.NoError(t, err)

	tests := []struct {
		name          string
		msg           string
		err           error
		keysAndValues []any
	}{
		{
			name:          "简单panic日志",
			msg:           "panic occurred",
			err:           nil,
			keysAndValues: []any{},
		},
		{
			name:          "带错误的panic日志",
			msg:           "critical error",
			err:           errors.New("system failure"),
			keysAndValues: []any{},
		},
		{
			name:          "带键值对的panic日志",
			msg:           "panic with context",
			err:           errors.New("database error"),
			keysAndValues: []any{"userID", 12345, "action", "update"},
		},
		{
			name:          "复杂panic场景",
			msg:           "multiple issues",
			err:           errors.New("connection lost"),
			keysAndValues: []any{"server", "db1", "retry", 3, "timeout", "30s"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ZPanic应该不会panic，只是记录日志
			// 注意：ZPanic实际调用的是Error级别，不会真正panic
			ZPanic(context.Background(), tt.msg, tt.err, tt.keysAndValues...)
		})
	}
}

// ==================== ToZap 方法测试 ====================

func TestZapLogger_ToZap(t *testing.T) {
	// 初始化测试logger
	err := InitConsoleLogger("test", LevelDebug, false, "1.0.0")
	assert.NoError(t, err)

	t.Run("ToZap返回SugaredLogger", func(t *testing.T) {
		// 获取全局logger并转换为ZapLogger
		logger, ok := pkgLogger.(*ZapLogger)
		if !ok {
			t.Fatal("pkgLogger is not *ZapLogger")
		}

		// 调用ToZap获取底层zap logger
		sugaredLogger := logger.ToZap()

		// 验证返回的不是nil
		if sugaredLogger == nil {
			t.Fatal("ToZap() returned nil")
		}

		// 验证可以使用返回的logger记录日志
		sugaredLogger.Info("test message from ToZap")
		sugaredLogger.Infow("test with fields", "key1", "value1", "key2", "value2")
	})

	t.Run("ToZap返回的logger可以使用所有zap方法", func(t *testing.T) {
		logger, ok := pkgLogger.(*ZapLogger)
		if !ok {
			t.Fatal("pkgLogger is not *ZapLogger")
		}
		zapLogger := logger.ToZap()

		// 测试各种日志级别
		zapLogger.Debug("debug message")
		zapLogger.Info("info message")
		zapLogger.Warn("warn message")
		zapLogger.Error("error message")

		// 测试带字段的日志
		zapLogger.Debugw("debug", "field", "value")
		zapLogger.Infow("info", "field", "value")
		zapLogger.Warnw("warn", "field", "value")
		zapLogger.Errorw("error", "field", "value")
	})
}

// ==================== Panic 方法测试 ====================

func TestZapLogger_Panic(t *testing.T) {
	// 初始化测试logger
	err := InitConsoleLogger("test", LevelDebug, false, "1.0.0")
	assert.NoError(t, err)

	tests := []struct {
		name          string
		msg           string
		err           error
		keysAndValues []any
		shouldPanic   bool
	}{
		{
			name:          "简单panic日志",
			msg:           "panic message",
			err:           nil,
			keysAndValues: []any{},
			shouldPanic:   true,
		},
		{
			name:          "带错误的panic",
			msg:           "critical failure",
			err:           errors.New("system crash"),
			keysAndValues: []any{},
			shouldPanic:   true,
		},
		{
			name:          "带上下文的panic",
			msg:           "panic with context",
			err:           errors.New("fatal error"),
			keysAndValues: []any{"component", "database", "operation", "write"},
			shouldPanic:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				// 验证Panic方法会真正panic
				defer func() {
					if r := recover(); r == nil {
						t.Error("Panic() should panic but didn't")
					}
				}()

				logger := pkgLogger
				logger.Panic(context.Background(), tt.msg, tt.err, tt.keysAndValues...)
			}
		})
	}
}

// ==================== 集成测试 ====================

func TestZapUncoveredFunctions_Integration(t *testing.T) {
	// 初始化测试logger
	err := InitConsoleLogger("test", LevelInfo, false, "1.0.0")
	assert.NoError(t, err)

	t.Run("ZPanic与其他日志函数配合", func(t *testing.T) {
		ctx := context.Background()

		ZDebug(ctx, "debug before panic")
		ZInfo(ctx, "info before panic")
		ZWarn(ctx, "warn before panic", nil)
		ZError(ctx, "error before panic", errors.New("test error"))
		ZPanic(ctx, "panic log", errors.New("panic error"), "key", "value")
		ZInfo(ctx, "info after panic")
	})

	t.Run("ToZap与标准logger对比", func(t *testing.T) {
		// 使用标准接口记录日志
		ZInfo(context.Background(), "standard info message")

		// 使用ToZap获取底层logger记录日志
		logger, ok := pkgLogger.(*ZapLogger)
		if !ok {
			t.Fatal("pkgLogger is not *ZapLogger")
		}
		zapLogger := logger.ToZap()
		zapLogger.Info("direct zap info message")

		// 两种方式都应该正常工作
	})
}
