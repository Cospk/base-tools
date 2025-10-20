package log

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/stretchr/testify/assert"
)

// TestContextPropagation 测试上下文信息传播
func TestContextPropagation(t *testing.T) {
	// 初始化日志到临时目录
	tmpDir := t.TempDir()
	err := InitLoggerFromConfig(
		"testLogger",
		"testModule",
		"",
		"",
		LevelDebug,
		false, // 不输出到控制台，只输出到文件
		true,  // JSON格式便于解析
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	assert.NoError(t, err)
	defer Flush()

	// 创建带有上下文信息的 context
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.OperationID, "test-operation-123")
	ctx = context.WithValue(ctx, constant.OpUserID, "user-456")
	ctx = context.WithValue(ctx, constant.ConnID, "conn-789")

	// 记录日志
	ZInfo(ctx, "test message with context", "key", "value")

	// 等待日志写入
	time.Sleep(100 * time.Millisecond)
	Flush()

	// 读取日志文件
	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, files)

	logFile := files[0].Name()
	content, err := os.ReadFile(tmpDir + "/" + logFile)
	assert.NoError(t, err)

	logContent := string(content)
	t.Logf("Log content: %s", logContent)

	// 验证上下文信息是否被记录
	assert.Contains(t, logContent, "test-operation-123", "Should contain operationID")
	assert.Contains(t, logContent, "user-456", "Should contain opUserID")
	assert.Contains(t, logContent, "conn-789", "Should contain connID")
	assert.Contains(t, logContent, "test message with context", "Should contain message")
	assert.Contains(t, logContent, "key", "Should contain custom key")
	assert.Contains(t, logContent, "value", "Should contain custom value")
}

// TestLogLevelFilter 测试日志级别过滤
func TestLogLevelFilter(t *testing.T) {
	tmpDir := t.TempDir()

	// 设置日志级别为 Info
	err := InitLoggerFromConfig(
		"testLogger",
		"testModule",
		"",
		"",
		LevelInfo, // 只记录 Info 及以上级别
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	assert.NoError(t, err)
	defer Flush()

	ctx := context.Background()

	// 记录不同级别的日志
	ZDebug(ctx, "debug message should not appear")
	ZInfo(ctx, "info message should appear")
	ZWarn(ctx, "warn message should appear", nil)
	ZError(ctx, "error message should appear", errors.New("test error"))

	time.Sleep(100 * time.Millisecond)
	Flush()

	// 读取日志文件
	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)

	if len(files) > 0 {
		content, err := os.ReadFile(tmpDir + "/" + files[0].Name())
		assert.NoError(t, err)
		logContent := string(content)

		// Debug 日志应该被过滤掉
		assert.NotContains(t, logContent, "debug message should not appear")

		// Info, Warn, Error 日志应该存在
		assert.Contains(t, logContent, "info message should appear")
		assert.Contains(t, logContent, "warn message should appear")
		assert.Contains(t, logContent, "error message should appear")
	}
}

// TestZAdaptive 测试自适应日志级别
func TestZAdaptive(t *testing.T) {
	tmpDir := t.TempDir()
	err := InitLoggerFromConfig(
		"testLogger",
		"testModule",
		"",
		"",
		LevelDebug,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	assert.NoError(t, err)
	defer Flush()

	ctx := context.Background()

	// 配置自适应级别
	AdaptiveErrorCodeLevel[errs.ErrInternalServer.Code()] = LevelError

	// 测试不同错误的自适应级别
	t.Run("Internal error should use Error level", func(t *testing.T) {
		err := errs.ErrInternalServer.Wrap()
		ZAdaptive(ctx, "adaptive log test", err, "testKey", "testValue")
		time.Sleep(50 * time.Millisecond)
		Flush()

		// 读取日志验证级别
		files, readErr := os.ReadDir(tmpDir)
		assert.NoError(t, readErr)
		if len(files) > 0 {
			content, readErr := os.ReadFile(tmpDir + "/" + files[0].Name())
			assert.NoError(t, readErr)
			assert.Contains(t, string(content), "adaptive log test")
		}
	})
}

// TestWithValues 测试带固定字段的子Logger
func TestWithValues(t *testing.T) {
	tmpDir := t.TempDir()
	err := InitLoggerFromConfig(
		"testLogger",
		"testModule",
		"",
		"",
		LevelDebug,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	assert.NoError(t, err)
	defer Flush()

	ctx := context.Background()

	// 创建带固定字段的子Logger
	subLogger := pkgLogger.WithValues("module", "user-service", "version", "1.0")

	// 使用子Logger记录日志
	subLogger.Info(ctx, "user login", "userId", 123)

	time.Sleep(100 * time.Millisecond)
	Flush()

	// 验证固定字段被记录
	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)
	if len(files) > 0 {
		content, err := os.ReadFile(tmpDir + "/" + files[0].Name())
		assert.NoError(t, err)
		logContent := string(content)

		assert.Contains(t, logContent, "user-service")
		assert.Contains(t, logContent, "1.0")
		assert.Contains(t, logContent, "user login")
		assert.Contains(t, logContent, "123")
	}
}

// TestConcurrentWrite 测试并发写入安全性
func TestConcurrentWrite(t *testing.T) {
	tmpDir := t.TempDir()
	err := InitLoggerFromConfig(
		"testLogger",
		"testModule",
		"",
		"",
		LevelDebug,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	assert.NoError(t, err)
	defer Flush()

	ctx := context.Background()
	const goroutines = 100
	const logsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// 并发写入日志
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < logsPerGoroutine; j++ {
				ZInfo(ctx, "concurrent log", "goroutine", id, "index", j)
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(200 * time.Millisecond)
	Flush()

	// 验证没有发生panic或数据竞争
	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, files, "Log file should be created")
}

// TestLogFormatterWithSimplify 测试 LogFormatter 接口与简化功能
func TestLogFormatterWithSimplify(t *testing.T) {
	tmpDir := t.TempDir()
	err := InitLoggerFromConfig(
		"testLogger",
		"testModule",
		"",
		"",
		LevelDebug,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		true, // 启用简化
	)
	assert.NoError(t, err)
	defer Flush()

	ctx := context.Background()

	// 创建一个大切片（超过30个元素）
	largeSlice := make([]int, 100)
	for i := 0; i < 100; i++ {
		largeSlice[i] = i
	}

	// 使用 Slice 类型记录日志
	ZInfo(ctx, "large slice test", "data", Slice[int](largeSlice))

	time.Sleep(100 * time.Millisecond)
	Flush()

	// 读取日志文件
	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)
	if len(files) > 0 {
		content, err := os.ReadFile(tmpDir + "/" + files[0].Name())
		assert.NoError(t, err)
		logContent := string(content)

		// 验证数组被截断（应该只包含前30个元素）
		assert.Contains(t, logContent, "large slice test")
		// JSON中数组长度应该是30而不是100
		// 简单检查：不应该包含较大的索引值
		assert.NotContains(t, logContent, "99", "Should not contain last element")
	}
}

// TestSDKLogWithCallStack 测试 SDK 日志的调用栈信息
func TestSDKLogWithCallStack(t *testing.T) {
	tmpDir := t.TempDir()
	err := InitLoggerFromConfig(
		"testLogger",
		"testModule",
		"TestSDK",
		"TestPlatform",
		LevelDebug,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	assert.NoError(t, err)
	defer Flush()

	ctx := context.Background()

	// 记录SDK日志
	SDKLog(ctx, LevelInfo, "sdk/test.go", 100, "SDK log message", nil, []any{"key", "value"})

	time.Sleep(100 * time.Millisecond)
	Flush()

	// 读取日志文件
	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)
	if len(files) > 0 {
		content, err := os.ReadFile(tmpDir + "/" + files[0].Name())
		assert.NoError(t, err)
		logContent := string(content)

		// 验证SDK日志包含调用位置信息
		assert.Contains(t, logContent, "SDK log message")
		assert.Contains(t, logContent, "native_caller")
		assert.Contains(t, logContent, "sdk/test.go:100")
	}
}

// TestFlush 测试日志刷新功能
func TestFlush(t *testing.T) {
	tmpDir := t.TempDir()
	err := InitLoggerFromConfig(
		"testLogger",
		"testModule",
		"",
		"",
		LevelDebug,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	assert.NoError(t, err)

	ctx := context.Background()
	ZInfo(ctx, "test flush")

	// 立即刷新
	Flush()

	// 验证日志已写入
	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, files)

	content, err := os.ReadFile(tmpDir + "/" + files[0].Name())
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test flush")
}

// TestConsoleLogger 测试控制台日志
func TestConsoleLogger(t *testing.T) {
	// 捕获标准输出
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := InitConsoleLogger(
		"testModule",
		LevelDebug,
		false, // 非JSON格式
		"1.0.0",
	)
	assert.NoError(t, err)

	ctx := context.Background()
	CInfo(ctx, "console test message", "key", "value")

	// 恢复标准输出
	w.Close()
	os.Stdout = oldStdout

	// 读取输出
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// 验证控制台输出
	assert.Contains(t, output, "console test message")
}

// BenchmarkZapLogger 性能基准测试
func BenchmarkZapLogger(b *testing.B) {
	tmpDir := b.TempDir()
	err := InitLoggerFromConfig(
		"benchLogger",
		"benchModule",
		"",
		"",
		LevelInfo,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	if err != nil {
		b.Fatal(err)
	}
	defer Flush()

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ZInfo(ctx, "benchmark message", "key", "value", "index", i)
	}
}

// BenchmarkZapLoggerParallel 并发性能基准测试
func BenchmarkZapLoggerParallel(b *testing.B) {
	tmpDir := b.TempDir()
	err := InitLoggerFromConfig(
		"benchLogger",
		"benchModule",
		"",
		"",
		LevelInfo,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	if err != nil {
		b.Fatal(err)
	}
	defer Flush()

	ctx := context.Background()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ZInfo(ctx, "parallel benchmark", "key", "value", "index", i)
			i++
		}
	})
}

// BenchmarkWithContext 测试带上下文的性能
func BenchmarkWithContext(b *testing.B) {
	tmpDir := b.TempDir()
	err := InitLoggerFromConfig(
		"benchLogger",
		"benchModule",
		"",
		"",
		LevelInfo,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	if err != nil {
		b.Fatal(err)
	}
	defer Flush()

	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.OperationID, "bench-op-123")
	ctx = context.WithValue(ctx, constant.OpUserID, "bench-user-456")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ZInfo(ctx, "benchmark with context", "key", "value")
	}
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	err := InitLoggerFromConfig(
		"testLogger",
		"testModule",
		"",
		"",
		LevelDebug,
		false,
		true,
		tmpDir,
		1,
		24,
		"1.0.0",
		false,
	)
	assert.NoError(t, err)
	defer Flush()

	ctx := context.Background()

	// 测试带错误的日志
	testErr := errors.New("test error message")
	ZError(ctx, "error occurred", testErr, "context", "test")

	time.Sleep(100 * time.Millisecond)
	Flush()

	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)
	if len(files) > 0 {
		content, err := os.ReadFile(tmpDir + "/" + files[0].Name())
		assert.NoError(t, err)
		logContent := string(content)

		assert.Contains(t, logContent, "error occurred")
		assert.Contains(t, logContent, "test error message")
		assert.Contains(t, logContent, "error") // 错误键
	}
}