package errs

import (
	"errors"
	"strings"
	"testing"
)

// ==================== NewCodeError 函数测试 ====================

func TestNewCodeError(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		msg     string
		wantErr string
	}{
		{
			name:    "简单错误码",
			code:    1000,
			msg:     "user not found",
			wantErr: "1000 user not found",
		},
		{
			name:    "负数错误码",
			code:    -1,
			msg:     "invalid code",
			wantErr: "-1 invalid code",
		},
		{
			name:    "零错误码",
			code:    0,
			msg:     "success",
			wantErr: "0 success",
		},
		{
			name:    "空消息",
			code:    100,
			msg:     "",
			wantErr: "100 ",
		},
		{
			name:    "包含特殊字符的消息",
			code:    200,
			msg:     "错误：数据库连接失败",
			wantErr: "200 错误：数据库连接失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCodeError(tt.code, tt.msg)

			// 验证返回非nil
			if err == nil {
				t.Fatal("NewCodeError() returned nil")
			}

			// 验证错误码
			if got := err.Code(); got != tt.code {
				t.Errorf("Code() = %d, want %d", got, tt.code)
			}

			// 验证消息
			if got := err.Msg(); got != tt.msg {
				t.Errorf("Msg() = %q, want %q", got, tt.msg)
			}

			// 验证Error()方法
			if got := err.Error(); got != tt.wantErr {
				t.Errorf("Error() = %q, want %q", got, tt.wantErr)
			}

			// 验证初始Detail为空
			if got := err.Detail(); got != "" {
				t.Errorf("Detail() = %q, want empty string", got)
			}
		})
	}
}

// ==================== CodeError.WithDetail 方法测试 ====================

func TestCodeError_WithDetail(t *testing.T) {
	tests := []struct {
		name       string
		baseErr    CodeError
		details    []string
		wantDetail string
		wantErr    string
	}{
		{
			name:       "单次添加详细信息",
			baseErr:    NewCodeError(1001, "database error"),
			details:    []string{"connection timeout"},
			wantDetail: "connection timeout",
			wantErr:    "1001 database error connection timeout",
		},
		{
			name:       "多次添加详细信息",
			baseErr:    NewCodeError(1002, "query failed"),
			details:    []string{"table not found", "user: admin"},
			wantDetail: "table not found, user: admin",
			wantErr:    "1002 query failed table not found, user: admin",
		},
		{
			name:       "添加空详细信息",
			baseErr:    NewCodeError(1003, "test"),
			details:    []string{""},
			wantDetail: "",
			wantErr:    "1003 test",
		},
		{
			name:       "链式调用多次",
			baseErr:    NewCodeError(1004, "operation failed"),
			details:    []string{"step1", "step2", "step3"},
			wantDetail: "step1, step2, step3",
			wantErr:    "1004 operation failed step1, step2, step3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err CodeError = tt.baseErr
			for _, detail := range tt.details {
				err = err.WithDetail(detail)
			}

			// 验证详细信息
			if got := err.Detail(); got != tt.wantDetail {
				t.Errorf("Detail() = %q, want %q", got, tt.wantDetail)
			}

			// 验证Error()包含详细信息
			if got := err.Error(); got != tt.wantErr {
				t.Errorf("Error() = %q, want %q", got, tt.wantErr)
			}

			// 验证错误码和消息没有改变
			if err.Code() != tt.baseErr.Code() {
				t.Errorf("Code changed after WithDetail")
			}
			if err.Msg() != tt.baseErr.Msg() {
				t.Errorf("Msg changed after WithDetail")
			}

			// 验证原始错误没有被修改（不可变性）
			if tt.baseErr.Detail() != "" {
				t.Errorf("Original error was modified, Detail() = %q", tt.baseErr.Detail())
			}
		})
	}
}

// ==================== CodeError.Is 方法测试 ====================

func TestCodeError_Is(t *testing.T) {
	err1 := NewCodeError(1000, "error 1")
	err2 := NewCodeError(2000, "error 2")
	err3 := NewCodeError(1000, "error 3") // 错误码与err1相同

	tests := []struct {
		name   string
		err    CodeError
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
			name:   "相同错误码的不同实例",
			err:    err1,
			target: err3,
			want:   true,
		},
		{
			name:   "不同错误码",
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
			target: errors.New("standard error"),
			want:   false,
		},
		{
			name:   "与普通Error比较",
			err:    err1,
			target: New("plain error"),
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

// ==================== CodeError.Wrap 和 WrapMsg 方法测试 ====================

func TestCodeError_Wrap(t *testing.T) {
	err := NewCodeError(1000, "test error")

	wrapped := err.Wrap()

	// 验证不是nil
	if wrapped == nil {
		t.Fatal("Wrap() returned nil")
	}

	// 验证错误消息包含原始消息
	if !strings.Contains(wrapped.Error(), "1000 test error") {
		t.Errorf("Wrap() error = %q, should contain '1000 test error'", wrapped.Error())
	}

	// 验证Unwrap能找到原始CodeError
	unwrapped := Unwrap(wrapped)
	if unwrapped == nil {
		t.Fatal("Unwrap() returned nil")
	}

	codeErr, ok := unwrapped.(CodeError)
	if !ok {
		t.Fatal("Unwrap() did not return CodeError")
	}

	if codeErr.Code() != err.Code() {
		t.Errorf("Unwrap() code = %d, want %d", codeErr.Code(), err.Code())
	}
}

func TestCodeError_WrapMsg(t *testing.T) {
	tests := []struct {
		name       string
		baseErr    CodeError
		msg        string
		kv         []any
		wantSubstr []string
	}{
		{
			name:       "添加简单消息",
			baseErr:    NewCodeError(1001, "base error"),
			msg:        "additional context",
			kv:         nil,
			wantSubstr: []string{"1001", "base error", "additional context"},
		},
		{
			name:       "添加消息和键值对",
			baseErr:    NewCodeError(1002, "query error"),
			msg:        "failed",
			kv:         []any{"table", "users", "id", 123},
			wantSubstr: []string{"1002", "query error", "failed", "table=users", "id=123"},
		},
		{
			name:       "只添加键值对",
			baseErr:    NewCodeError(1003, "error"),
			msg:        "",
			kv:         []any{"reason", "timeout"},
			wantSubstr: []string{"1003", "error", "reason=timeout"},
		},
		{
			name:       "在已有detail上添加",
			baseErr:    NewCodeError(1004, "test").WithDetail("existing detail"),
			msg:        "new info",
			kv:         nil,
			wantSubstr: []string{"1004", "test", "existing detail", "new info"},
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

			// 验证原始错误没有被修改
			if tt.baseErr.Detail() != "" && tt.baseErr.Detail() != "existing detail" {
				t.Errorf("Original error was modified")
			}
		})
	}
}

// ==================== Unwrap 函数测试 ====================

func TestUnwrap(t *testing.T) {
	tests := []struct {
		name       string
		setupErr   func() error
		wantCode   int
		wantNil    bool
		wantString string
	}{
		{
			name: "解包CodeError",
			setupErr: func() error {
				return NewCodeError(1000, "test").Wrap()
			},
			wantCode: 1000,
		},
		{
			name: "解包多层包装的CodeError",
			setupErr: func() error {
				return NewCodeError(1001, "test").WrapMsg("extra context")
			},
			wantCode: 1001,
		},
		{
			name: "解包普通Error",
			setupErr: func() error {
				return New("plain error").Wrap()
			},
			wantString: "plain error",
		},
		{
			name: "nil错误",
			setupErr: func() error {
				return nil
			},
			wantNil: true,
		},
		{
			name: "标准错误（不可解包）",
			setupErr: func() error {
				return errors.New("standard error")
			},
			wantString: "standard error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupErr()
			unwrapped := Unwrap(err)

			if tt.wantNil {
				if unwrapped != nil {
					t.Errorf("Unwrap() = %v, want nil", unwrapped)
				}
				return
			}

			if unwrapped == nil {
				t.Fatal("Unwrap() returned nil")
			}

			if tt.wantCode != 0 {
				codeErr, ok := unwrapped.(CodeError)
				if !ok {
					t.Fatal("Unwrap() did not return CodeError")
				}
				if codeErr.Code() != tt.wantCode {
					t.Errorf("Code() = %d, want %d", codeErr.Code(), tt.wantCode)
				}
			}

			if tt.wantString != "" {
				if !strings.Contains(unwrapped.Error(), tt.wantString) {
					t.Errorf("Error() = %q, should contain %q", unwrapped.Error(), tt.wantString)
				}
			}
		})
	}
}

// ==================== Wrap 和 WrapMsg 函数测试 ====================

func TestWrap(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantNil bool
	}{
		{
			name:    "包装标准错误",
			err:     errors.New("standard error"),
			wantNil: false,
		},
		{
			name:    "包装CodeError",
			err:     NewCodeError(1000, "code error"),
			wantNil: false,
		},
		{
			name:    "包装nil",
			err:     nil,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := Wrap(tt.err)

			if tt.wantNil {
				if wrapped != nil {
					t.Errorf("Wrap(nil) = %v, want nil", wrapped)
				}
				return
			}

			if wrapped == nil {
				t.Fatal("Wrap() returned nil")
			}

			// 验证能解包回原始错误
			unwrapped := Unwrap(wrapped)
			if unwrapped == nil {
				t.Fatal("Unwrap() returned nil")
			}
		})
	}
}

func TestWrapMsg(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		msg        string
		kv         []any
		wantNil    bool
		wantSubstr []string
	}{
		{
			name:       "包装标准错误并添加消息",
			err:        errors.New("base error"),
			msg:        "context",
			kv:         []any{"key", "value"},
			wantSubstr: []string{"base error", "context", "key=value"},
		},
		{
			name:       "包装CodeError并添加消息",
			err:        NewCodeError(1000, "code error"),
			msg:        "extra info",
			kv:         nil,
			wantSubstr: []string{"1000", "code error", "extra info"},
		},
		{
			name:    "包装nil错误",
			err:     nil,
			msg:     "message",
			kv:      []any{"key", "value"},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := WrapMsg(tt.err, tt.msg, tt.kv...)

			if tt.wantNil {
				if wrapped != nil {
					t.Errorf("WrapMsg(nil, ...) = %v, want nil", wrapped)
				}
				return
			}

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

// ==================== CodeRelation 测试 ====================

func TestCodeRelation_Add(t *testing.T) {
	tests := []struct {
		name      string
		codes     []int
		wantError bool
	}{
		{
			name:      "添加简单关系链",
			codes:     []int{1000, 1001, 1002},
			wantError: false,
		},
		{
			name:      "添加两个错误码",
			codes:     []int{2000, 2001},
			wantError: false,
		},
		{
			name:      "只有一个错误码",
			codes:     []int{3000},
			wantError: true,
		},
		{
			name:      "空切片",
			codes:     []int{},
			wantError: true,
		},
		{
			name:      "长关系链",
			codes:     []int{4000, 4001, 4002, 4003, 4004},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := newCodeRelation()
			err := cr.Add(tt.codes...)

			if (err != nil) != tt.wantError {
				t.Errorf("Add() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestCodeRelation_Is(t *testing.T) {
	cr := newCodeRelation()

	// 设置错误码关系：1000 -> 1001 -> 1002 -> 1003
	cr.Add(1000, 1001, 1002, 1003)
	// 设置错误码关系：2000 -> 2001
	cr.Add(2000, 2001)

	tests := []struct {
		name   string
		parent int
		child  int
		want   bool
	}{
		{
			name:   "直接父子关系",
			parent: 1000,
			child:  1001,
			want:   true,
		},
		{
			name:   "间接父子关系",
			parent: 1000,
			child:  1002,
			want:   true,
		},
		{
			name:   "更深层父子关系",
			parent: 1000,
			child:  1003,
			want:   true,
		},
		{
			name:   "相同错误码",
			parent: 1000,
			child:  1000,
			want:   true,
		},
		{
			name:   "无关系的错误码",
			parent: 1000,
			child:  2000,
			want:   false,
		},
		{
			name:   "反向关系",
			parent: 1001,
			child:  1000,
			want:   false,
		},
		{
			name:   "不存在的错误码",
			parent: 9999,
			child:  1000,
			want:   false,
		},
		{
			name:   "另一个关系链",
			parent: 2000,
			child:  2001,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cr.Is(tt.parent, tt.child); got != tt.want {
				t.Errorf("Is(%d, %d) = %v, want %v", tt.parent, tt.child, got, tt.want)
			}
		})
	}
}

func TestCodeRelation_IsWithDefaultRelation(t *testing.T) {
	// 测试使用全局DefaultCodeRelation
	// 注意：这会影响全局状态，所以使用唯一的错误码
	testParent := 9000
	testChild := 9001

	err := DefaultCodeRelation.Add(testParent, testChild)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// 创建两个CodeError并测试Is方法
	parentErr := NewCodeError(testParent, "parent error")
	childErr := NewCodeError(testChild, "child error")

	// 子错误应该能匹配父错误
	if !parentErr.Is(childErr) {
		t.Errorf("parentErr.Is(childErr) = false, want true")
	}

	// 父错误不应该能匹配子错误
	if childErr.Is(parentErr) {
		t.Errorf("childErr.Is(parentErr) = true, want false")
	}

	// 相同错误码应该匹配
	sameErr := NewCodeError(testParent, "same code")
	if !parentErr.Is(sameErr) {
		t.Errorf("parentErr.Is(sameErr) = false, want true")
	}
}

// ==================== 边界条件测试 ====================

func TestCodeError_EdgeCases(t *testing.T) {
	t.Run("非常大的错误码", func(t *testing.T) {
		err := NewCodeError(999999999, "large code")
		if err.Code() != 999999999 {
			t.Errorf("Code() = %d, want 999999999", err.Code())
		}
	})

	t.Run("多次WithDetail调用", func(t *testing.T) {
		err := NewCodeError(1000, "test")
		for i := 0; i < 100; i++ {
			err = err.WithDetail("detail")
		}
		// 确保不会panic或崩溃
		if err.Detail() == "" {
			t.Error("Detail should not be empty after 100 calls")
		}
	})

	t.Run("特殊字符处理", func(t *testing.T) {
		err := NewCodeError(1000, "错误：\n换行\t制表符")
		if !strings.Contains(err.Error(), "错误：\n换行\t制表符") {
			t.Error("Should handle special characters")
		}
	})

	t.Run("Detail包含中文", func(t *testing.T) {
		err := NewCodeError(1000, "test").WithDetail("详细信息：数据库连接失败")
		if !strings.Contains(err.Detail(), "详细信息") {
			t.Error("Should handle Chinese characters in detail")
		}
	})
}
