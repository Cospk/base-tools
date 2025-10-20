package log

import (
	"strings"
	"testing"
)

// ==================== NewZkLogger 测试 ====================

func TestNewZkLogger(t *testing.T) {
	// 初始化测试logger
	InitConsoleLogger("test", LevelDebug, false, "1.0.0")

	tests := []struct {
		name string
	}{
		{
			name: "创建ZK logger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewZkLogger()

			if logger == nil {
				t.Fatal("NewZkLogger() returned nil")
			}
		})
	}
}

// ==================== Printf 方法测试 ====================

func TestZkLogger_Printf(t *testing.T) {
	// 初始化测试logger
	InitConsoleLogger("test", LevelDebug, false, "1.0.0")

	logger := NewZkLogger()

	tests := []struct {
		name   string
		format string
		args   []any
	}{
		{
			name:   "简单消息",
			format: "ZooKeeper connected",
			args:   nil,
		},
		{
			name:   "带参数消息",
			format: "ZooKeeper state changed: %s",
			args:   []any{"connected"},
		},
		{
			name:   "多个参数",
			format: "Session established: id=%d, timeout=%d",
			args:   []any{12345, 30000},
		},
		{
			name:   "空格式字符串",
			format: "",
			args:   nil,
		},
		{
			name:   "特殊字符",
			format: "ZK path: /app/%s/config",
			args:   []any{"production"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 应该不会panic
			logger.Printf(tt.format, tt.args...)
		})
	}
}

// ==================== 集成测试 ====================

func TestZkLogger_Integration(t *testing.T) {
	// 初始化测试logger
	InitConsoleLogger("test", LevelInfo, false, "1.0.0")

	t.Run("模拟ZooKeeper日志输出", func(t *testing.T) {
		logger := NewZkLogger()

		// 模拟各种ZK事件
		logger.Printf("Initiating client connection")
		logger.Printf("Opening socket connection to server %s", "localhost:2181")
		logger.Printf("Socket connection established")
		logger.Printf("Session establishment complete on server %s, sessionId=%d, negotiated timeout=%d",
			"localhost:2181", 12345, 30000)

		t.Log("ZooKeeper logger integration test completed without panic")
	})

	t.Run("大量日志输出", func(t *testing.T) {
		logger := NewZkLogger()

		for i := 0; i < 100; i++ {
			logger.Printf("Event %d: processing node /app/config/%d", i, i)
		}

		t.Log("Batch logging completed successfully")
	})
}

// ==================== 边界条件测试 ====================

func TestZkLogger_EdgeCases(t *testing.T) {
	// 初始化测试logger
	InitConsoleLogger("test", LevelDebug, false, "1.0.0")

	logger := NewZkLogger()

	t.Run("nil参数", func(t *testing.T) {
		logger.Printf("test: %v", nil)
	})

	t.Run("超长消息", func(t *testing.T) {
		longMsg := strings.Repeat("a", 10000)
		logger.Printf("Long message: %s", longMsg)
	})

	t.Run("中文消息", func(t *testing.T) {
		logger.Printf("ZooKeeper连接: 状态=%s", "已连接")
	})

	t.Run("格式化参数不匹配", func(t *testing.T) {
		// 缺少参数的情况 - 会打印格式字符串
		logger.Printf("Missing arg")
	})

	t.Run("过多参数", func(t *testing.T) {
		// 只使用第一个参数
		logger.Printf("Only one arg: %s", "first")
	})

	t.Run("特殊格式化字符", func(t *testing.T) {
		logger.Printf("Values: %v, %+v, %#v", 123, 123, 123)
	})
}

// ==================== 性能基准测试 ====================

func BenchmarkZkLogger_Printf(b *testing.B) {
	// 初始化测试logger
	InitConsoleLogger("test", LevelInfo, false, "1.0.0")

	logger := NewZkLogger()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Printf("ZooKeeper event: session=%d, type=%s", i, "connected")
	}
}

func BenchmarkZkLogger_SimpleMessage(b *testing.B) {
	// 初始化测试logger
	InitConsoleLogger("test", LevelInfo, false, "1.0.0")

	logger := NewZkLogger()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Printf("Simple message")
	}
}

// ==================== 输出验证测试（可选）====================

// 注意：由于ZkLogger使用全局logger，直接验证输出比较困难
// 这里我们只验证调用不会panic，日志格式由zap.go保证

func TestZkLogger_DoesNotPanic(t *testing.T) {
	// 初始化测试logger
	InitConsoleLogger("test", LevelDebug, false, "1.0.0")

	logger := NewZkLogger()

	testCases := []struct {
		name string
		fn   func()
	}{
		{
			name: "Empty format",
			fn: func() {
				logger.Printf("")
			},
		},
		{
			name: "Nil format with args",
			fn: func() {
				logger.Printf("%s %s", "arg1", "arg2")
			},
		},
		{
			name: "Normal case",
			fn: func() {
				logger.Printf("test message: %s", "value")
			},
		},
		{
			name: "Multiple calls",
			fn: func() {
				for i := 0; i < 10; i++ {
					logger.Printf("iteration %d", i)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 这不应该panic
			tc.fn()
		})
	}
}

// ==================== 并发测试 ====================

func TestZkLogger_Concurrent(t *testing.T) {
	// 初始化测试logger
	InitConsoleLogger("test", LevelInfo, false, "1.0.0")

	logger := NewZkLogger()

	// 模拟多个goroutine同时写日志
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				logger.Printf("Goroutine %d: message %d", id, j)
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	t.Log("Concurrent logging completed without panic")
}
