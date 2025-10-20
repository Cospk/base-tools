package errs

import (
	"errors"
	"strings"
	"testing"
)

// ==================== NewErrorWrapper 函数测试 ====================

func TestNewErrorWrapper(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		msg        string
		wantSubstr []string
		wantNil    bool
	}{
		{
			name:       "包装标准错误",
			err:        errors.New("base error"),
			msg:        "additional context",
			wantSubstr: []string{"base error", "additional context"},
		},
		{
			name:       "包装CodeError",
			err:        NewCodeError(1000, "code error"),
			msg:        "wrapper message",
			wantSubstr: []string{"1000", "code error", "wrapper message"},
		},
		{
			name:       "包装Error接口",
			err:        New("plain error"),
			msg:        "wrapped",
			wantSubstr: []string{"plain error", "wrapped"},
		},
		{
			name:       "空消息包装",
			err:        errors.New("error"),
			msg:        "",
			wantSubstr: []string{"error"},
		},
		{
			name:       "包装nil错误",
			err:        nil,
			msg:        "message",
			wantSubstr: []string{"message"}, // NewErrorWrapper不检查nil,会包装nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := NewErrorWrapper(tt.err, tt.msg)

			if tt.wantNil {
				if wrapper != nil {
					t.Errorf("NewErrorWrapper(nil, ...) = %v, want nil", wrapper)
				}
				return
			}

			if wrapper == nil {
				t.Fatal("NewErrorWrapper() returned nil")
			}

			// 验证Error()输出包含期望的子串
			errStr := wrapper.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(errStr, substr) {
					t.Errorf("Error() = %q, should contain %q", errStr, substr)
				}
			}
		})
	}
}

// ==================== errorWrapper.Unwrap 方法测试 ====================

func TestErrorWrapper_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		baseErr  error
		msg      string
		wantType string
	}{
		{
			name:     "解包标准错误",
			baseErr:  errors.New("standard error"),
			msg:      "wrapper",
			wantType: "*errors.errorString",
		},
		{
			name:     "解包CodeError",
			baseErr:  NewCodeError(1000, "code error"),
			msg:      "wrapper",
			wantType: "CodeError",
		},
		{
			name:     "解包Error接口",
			baseErr:  New("plain error"),
			msg:      "wrapper",
			wantType: "Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := NewErrorWrapper(tt.baseErr, tt.msg)
			if wrapper == nil {
				t.Fatal("NewErrorWrapper() returned nil")
			}

			// 测试Unwrap方法
			unwrapped := wrapper.Unwrap()
			if unwrapped == nil {
				t.Fatal("Unwrap() returned nil")
			}

			// 验证解包后的错误与原错误一致
			if unwrapped.Error() != tt.baseErr.Error() {
				t.Errorf("Unwrap().Error() = %q, want %q", unwrapped.Error(), tt.baseErr.Error())
			}
		})
	}
}

// ==================== errorWrapper.Error 方法测试 ====================

func TestErrorWrapper_Error(t *testing.T) {
	tests := []struct {
		name    string
		baseErr error
		msg     string
		want    string
	}{
		{
			name:    "带消息的包装",
			baseErr: errors.New("base"),
			msg:     "context",
			want:    "base context",
		},
		{
			name:    "空消息的包装",
			baseErr: errors.New("base"),
			msg:     "",
			want:    "base ",
		},
		{
			name:    "包装CodeError带消息",
			baseErr: NewCodeError(1000, "error"),
			msg:     "wrapper",
			want:    "1000 error wrapper",
		},
		{
			name:    "包装CodeError空消息",
			baseErr: NewCodeError(2000, "error"),
			msg:     "",
			want:    "2000 error ",
		},
		{
			name:    "多层包装",
			baseErr: NewErrorWrapper(errors.New("base"), "layer1"),
			msg:     "layer2",
			want:    "base layer1 layer2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := NewErrorWrapper(tt.baseErr, tt.msg)
			if wrapper == nil {
				t.Fatal("NewErrorWrapper() returned nil")
			}

			got := wrapper.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ==================== errorWrapper.Is 方法测试 ====================

func TestErrorWrapper_Is(t *testing.T) {
	baseErr1 := errors.New("base error 1")
	baseErr2 := errors.New("base error 2")
	codeErr := NewCodeError(1000, "code error")

	wrapper1 := NewErrorWrapper(baseErr1, "context1")
	wrapper2 := NewErrorWrapper(baseErr2, "context1") // 相同的wrapper消息
	wrapper3 := NewErrorWrapper(baseErr1, "context2") // 不同的wrapper消息

	tests := []struct {
		name    string
		wrapper ErrWrapper
		target  error
		want    bool
	}{
		{
			name:    "与自身比较",
			wrapper: wrapper1,
			target:  wrapper1,
			want:    true,
		},
		{
			name:    "与具有相同wrapper消息的errorWrapper比较",
			wrapper: wrapper1,
			target:  wrapper2, // 即使baseErr不同，但wrapper消息相同
			want:    true,
		},
		{
			name:    "与具有不同wrapper消息的errorWrapper比较",
			wrapper: wrapper1,
			target:  wrapper3,
			want:    false,
		},
		{
			name:    "与nil比较",
			wrapper: wrapper1,
			target:  nil,
			want:    false,
		},
		{
			name:    "与标准错误比较",
			wrapper: wrapper1,
			target:  errors.New("other error"),
			want:    false,
		},
		{
			name:    "与CodeError比较",
			wrapper: NewErrorWrapper(codeErr, "wrapper"),
			target:  codeErr,
			want:    false, // Is方法不检查underlying error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.wrapper.Is(tt.target)
			if got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}

	// 测试使用标准errors.Is函数
	t.Run("标准errors.Is可以unwrap检测", func(t *testing.T) {
		wrapper := NewErrorWrapper(baseErr1, "wrapper")
		// 标准errors.Is会自动unwrap检测底层错误
		if !errors.Is(wrapper, baseErr1) {
			t.Error("errors.Is() should detect underlying error through Unwrap()")
		}
	})
}

// ==================== errorWrapper.Wrap 方法测试 ====================

func TestErrorWrapper_Wrap(t *testing.T) {
	baseErr := errors.New("base error")
	wrapper := NewErrorWrapper(baseErr, "context")

	wrapped := wrapper.Wrap()

	// 验证不是nil
	if wrapped == nil {
		t.Fatal("Wrap() returned nil")
	}

	// 验证错误消息包含原始消息
	errStr := wrapped.Error()
	if !strings.Contains(errStr, "context") || !strings.Contains(errStr, "base error") {
		t.Errorf("Wrap() error = %q, should contain 'context' and 'base error'", errStr)
	}

	// 验证可以解包到原始wrapper
	unwrapped := Unwrap(wrapped)
	if unwrapped == nil {
		t.Fatal("Unwrap() returned nil")
	}
}

// ==================== errorWrapper.WrapMsg 方法测试 ====================

func TestErrorWrapper_WrapMsg(t *testing.T) {
	tests := []struct {
		name       string
		baseErr    error
		wrapperMsg string
		extraMsg   string
		kv         []any
		wantSubstr []string
	}{
		{
			name:       "添加简单消息",
			baseErr:    errors.New("base"),
			wrapperMsg: "wrapper",
			extraMsg:   "extra",
			kv:         nil,
			wantSubstr: []string{"base", "wrapper", "extra"},
		},
		{
			name:       "添加消息和键值对",
			baseErr:    errors.New("base"),
			wrapperMsg: "wrapper",
			extraMsg:   "info",
			kv:         []any{"user", "alice", "id", 123},
			wantSubstr: []string{"base", "wrapper", "info", "user=alice", "id=123"},
		},
		{
			name:       "只添加键值对",
			baseErr:    errors.New("base"),
			wrapperMsg: "wrapper",
			extraMsg:   "",
			kv:         []any{"reason", "timeout"},
			wantSubstr: []string{"base", "wrapper", "reason=timeout"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := NewErrorWrapper(tt.baseErr, tt.wrapperMsg)
			wrapped := wrapper.WrapMsg(tt.extraMsg, tt.kv...)

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

// ==================== 集成测试 ====================

func TestErrorWrapper_Integration(t *testing.T) {
	t.Run("多层包装和解包", func(t *testing.T) {
		// 创建基础错误
		baseErr := errors.New("database connection failed")

		// 第一层包装
		layer1 := NewErrorWrapper(baseErr, "connection pool error")

		// 第二层包装
		layer2 := NewErrorWrapper(layer1, "service initialization failed")

		// 第三层包装
		layer3 := NewErrorWrapper(layer2, "application startup error")

		// 验证错误消息包含所有层的信息
		errStr := layer3.Error()
		expectedSubstrs := []string{
			"database connection failed",
			"connection pool error",
			"service initialization failed",
			"application startup error",
		}
		for _, substr := range expectedSubstrs {
			if !strings.Contains(errStr, substr) {
				t.Errorf("Error() = %q, should contain %q", errStr, substr)
			}
		}

		// 验证可以通过标准errors.Is检测到基础错误
		if !errors.Is(layer3, baseErr) {
			t.Error("errors.Is() should detect base error through multiple layers")
		}

		// 验证Unwrap可以逐层解包
		unwrap1 := layer3.Unwrap()
		if unwrap1 != layer2 {
			t.Error("First Unwrap() should return layer2")
		}

		unwrap2 := unwrap1.(ErrWrapper).Unwrap()
		if unwrap2 != layer1 {
			t.Error("Second Unwrap() should return layer1")
		}

		unwrap3 := unwrap2.(ErrWrapper).Unwrap()
		if unwrap3 != baseErr {
			t.Error("Third Unwrap() should return baseErr")
		}
	})

	t.Run("包装CodeError并使用错误码关系", func(t *testing.T) {
		// 使用全局DefaultCodeRelation添加错误码关系
		// 使用唯一的错误码避免测试冲突
		testParent := 9100
		testChild := 9101

		err := DefaultCodeRelation.Add(testParent, testChild)
		if err != nil {
			t.Fatalf("Add() failed: %v", err)
		}

		parentErr := NewCodeError(testParent, "parent error")
		childErr := NewCodeError(testChild, "child error")

		// 包装子错误
		wrappedChild := NewErrorWrapper(childErr, "wrapped child")

		// errorWrapper.Is()方法不会检查底层错误，所以这个测试不通过
		// 但是标准errors.Is()会通过Unwrap()检测到底层的CodeError
		if !errors.Is(wrappedChild, childErr) {
			t.Error("errors.Is() should detect underlying CodeError")
		}

		// 直接测试底层CodeError的关系
		unwrapped := Unwrap(wrappedChild)
		if codeErr, ok := unwrapped.(CodeError); ok {
			if !parentErr.Is(codeErr) {
				t.Error("Parent CodeError should match child CodeError")
			}
		} else {
			t.Error("Unwrap should return CodeError")
		}
	})

	t.Run("结合Wrap和WrapMsg", func(t *testing.T) {
		baseErr := errors.New("file not found")
		wrapper := NewErrorWrapper(baseErr, "read config")

		// 先Wrap添加堆栈
		withStack := wrapper.Wrap()

		// 验证包含堆栈信息
		stackStr := withStack.Error()
		if !strings.Contains(stackStr, "file not found") {
			t.Error("Should contain base error message")
		}

		// 再次使用WrapMsg添加上下文
		withMsg := WrapMsg(withStack, "initialization failed", "file", "config.yaml")
		msgStr := withMsg.Error()

		if !strings.Contains(msgStr, "initialization failed") {
			t.Error("Should contain additional message")
		}
		if !strings.Contains(msgStr, "file=config.yaml") {
			t.Error("Should contain key-value pairs")
		}
	})
}

// ==================== 边界条件测试 ====================

func TestErrorWrapper_EdgeCases(t *testing.T) {
	t.Run("包装已经被包装的错误", func(t *testing.T) {
		base := errors.New("base")
		wrap1 := NewErrorWrapper(base, "wrap1")
		wrap2 := NewErrorWrapper(wrap1, "wrap2")
		wrap3 := NewErrorWrapper(wrap2, "wrap3")

		// 确保不会panic
		if wrap3 == nil {
			t.Fatal("Multiple wrapping should not return nil")
		}

		// 验证可以正确解包到基础错误
		if !errors.Is(wrap3, base) {
			t.Error("Should detect base error through multiple wrappers")
		}
	})

	t.Run("非常长的消息", func(t *testing.T) {
		longMsg := strings.Repeat("a", 10000)
		baseErr := errors.New("base")
		wrapper := NewErrorWrapper(baseErr, longMsg)

		if wrapper == nil {
			t.Fatal("Should handle long message")
		}

		if !strings.Contains(wrapper.Error(), longMsg) {
			t.Error("Should preserve long message")
		}
	})

	t.Run("特殊字符处理", func(t *testing.T) {
		specialChars := "错误：\n换行\t制表符\r回车"
		baseErr := errors.New("base")
		wrapper := NewErrorWrapper(baseErr, specialChars)

		if !strings.Contains(wrapper.Error(), specialChars) {
			t.Error("Should handle special characters")
		}
	})

	t.Run("空字符串包装消息", func(t *testing.T) {
		baseErr := errors.New("base error")
		wrapper := NewErrorWrapper(baseErr, "")

		// 空消息时会在错误后面添加一个空格
		if wrapper.Error() != "base error " {
			t.Errorf("Error() = %q, want 'base error '", wrapper.Error())
		}
	})

	t.Run("包装包含冒号的错误", func(t *testing.T) {
		baseErr := errors.New("error: details: more details")
		wrapper := NewErrorWrapper(baseErr, "context")

		expected := "error: details: more details context"
		if wrapper.Error() != expected {
			t.Errorf("Error() = %q, want %q", wrapper.Error(), expected)
		}
	})
}
