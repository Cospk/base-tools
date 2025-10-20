package errs

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// ==================== New 函数测试 ====================

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		kv      []any
		wantMsg string
	}{
		{
			name:    "简单错误消息",
			msg:     "user not found",
			kv:      nil,
			wantMsg: "user not found",
		},
		{
			name:    "带单个键值对",
			msg:     "database error",
			kv:      []any{"table", "users"},
			wantMsg: "database error, table=users",
		},
		{
			name:    "带多个键值对",
			msg:     "query failed",
			kv:      []any{"table", "users", "id", 123},
			wantMsg: "query failed, table=users, id=123",
		},
		{
			name:    "空消息带键值对",
			msg:     "",
			kv:      []any{"key", "value"},
			wantMsg: "key=value",
		},
		{
			name:    "奇数个键值对-缺失值",
			msg:     "odd params",
			kv:      []any{"key1", "value1", "key2"},
			wantMsg: "odd params, key1=value1, key2=MISSING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.msg, tt.kv...)

			// 验证错误消息
			if got := err.Error(); got != tt.wantMsg {
				t.Errorf("Error() = %q, want %q", got, tt.wantMsg)
			}

			// 验证返回的是Error接口
			if err == nil {
				t.Fatal("New() returned nil")
			}
		})
	}
}

// ==================== Error.Is 方法测试 ====================

func TestError_Is(t *testing.T) {
	err1 := New("error 1")
	err2 := New("error 2")
	err3 := New("error 1") // 与err1消息相同

	tests := []struct {
		name   string
		err    Error
		target error
		want   bool
	}{
		{
			name:   "相同错误实例",
			err:    err1,
			target: err1,
			want:   true,
		},
		{
			name:   "相同消息的不同实例",
			err:    err1,
			target: err3,
			want:   true,
		},
		{
			name:   "不同错误",
			err:    err1,
			target: err2,
			want:   false,
		},
		{
			name:   "与nil比较",
			err:    err1,
			target: nil,
			want:   false,
		},
		{
			name:   "与标准错误比较",
			err:    err1,
			target: errors.New("error 1"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Is(tt.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ==================== Error.Wrap 方法测试 ====================

func TestError_Wrap(t *testing.T) {
	err := New("test error")

	wrapped := err.Wrap()

	// 验证不是nil
	if wrapped == nil {
		t.Fatal("Wrap() returned nil")
	}

	// 验证错误消息包含原始消息
	if !strings.Contains(wrapped.Error(), "test error") {
		t.Errorf("Wrap() error = %q, should contain 'test error'", wrapped.Error())
	}

	// 验证Unwrap能找到原始错误
	unwrapped := Unwrap(wrapped)
	if unwrapped == nil {
		t.Fatal("Unwrap() returned nil")
	}

	if unwrapped.Error() != err.Error() {
		t.Errorf("Unwrap() = %q, want %q", unwrapped.Error(), err.Error())
	}
}

// ==================== Error.WrapMsg 方法测试 ====================

func TestError_WrapMsg(t *testing.T) {
	tests := []struct {
		name       string
		baseErr    Error
		msg        string
		kv         []any
		wantSubstr []string
	}{
		{
			name:       "添加简单消息",
			baseErr:    New("base error"),
			msg:        "additional context",
			kv:         nil,
			wantSubstr: []string{"base error", "additional context"},
		},
		{
			name:       "添加消息和键值对",
			baseErr:    New("base error"),
			msg:        "context",
			kv:         []any{"user", "alice"},
			wantSubstr: []string{"base error", "context", "user=alice"},
		},
		{
			name:       "只添加键值对",
			baseErr:    New("base error"),
			msg:        "",
			kv:         []any{"id", 123},
			wantSubstr: []string{"base error", "id=123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := tt.baseErr.WrapMsg(tt.msg, tt.kv...)

			if wrapped == nil {
				t.Fatal("WrapMsg() returned nil")
			}

			errStr := wrapped.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(errStr, substr) {
					t.Errorf("WrapMsg() error = %q, should contain %q", errStr, substr)
				}
			}
		})
	}
}

// ==================== toString 辅助函数测试 ====================

func TestToString(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		kv   []any
		want string
	}{
		{
			name: "无键值对",
			msg:  "simple message",
			kv:   nil,
			want: "simple message",
		},
		{
			name: "单个键值对",
			msg:  "error",
			kv:   []any{"key", "value"},
			want: "error, key=value",
		},
		{
			name: "多个键值对",
			msg:  "error",
			kv:   []any{"k1", "v1", "k2", "v2"},
			want: "error, k1=v1, k2=v2",
		},
		{
			name: "奇数个参数",
			msg:  "error",
			kv:   []any{"k1", "v1", "k2"},
			want: "error, k1=v1, k2=MISSING",
		},
		{
			name: "空消息",
			msg:  "",
			kv:   []any{"key", "value"},
			want: "key=value",
		},
		{
			name: "数值类型",
			msg:  "error",
			kv:   []any{"count", 42, "rate", 3.14},
			want: "error, count=42, rate=3.14",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toString(tt.msg, tt.kv)
			if got != tt.want {
				t.Errorf("toString() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ==================== 边界条件测试 ====================

func TestError_EdgeCases(t *testing.T) {
	t.Run("空字符串消息", func(t *testing.T) {
		err := New("")
		if err.Error() != "" {
			t.Errorf("New(\"\") = %q, want \"\"", err.Error())
		}
	})

	t.Run("非常长的消息", func(t *testing.T) {
		longMsg := strings.Repeat("a", 10000)
		err := New(longMsg)
		if len(err.Error()) != 10000 {
			t.Errorf("Long message length = %d, want 10000", len(err.Error()))
		}
	})

	t.Run("大量键值对", func(t *testing.T) {
		kv := make([]any, 100)
		for i := 0; i < 100; i += 2 {
			kv[i] = fmt.Sprintf("key%d", i)
			kv[i+1] = i + 1
		}
		err := New("test", kv...)
		if err == nil {
			t.Fatal("New() with many kv returned nil")
		}
		// 验证包含一些键值对
		if !strings.Contains(err.Error(), "key0=1") {
			t.Error("Should contain key0=1")
		}
	})

	t.Run("nil值在键值对中", func(t *testing.T) {
		err := New("test", "key", nil)
		if !strings.Contains(err.Error(), "key=<nil>") {
			t.Errorf("Error() = %q, should handle nil value", err.Error())
		}
	})

	t.Run("特殊字符", func(t *testing.T) {
		err := New("error: 中文错误", "键", "值")
		if !strings.Contains(err.Error(), "中文错误") {
			t.Error("Should handle Chinese characters")
		}
		if !strings.Contains(err.Error(), "键=值") {
			t.Error("Should handle Chinese key-value pairs")
		}
	})
}
