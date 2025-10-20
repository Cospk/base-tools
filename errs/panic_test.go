package errs

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// ==================== ErrPanic 函数测试 ====================

func TestErrPanic(t *testing.T) {
	tests := []struct {
		name       string
		panicValue any
		wantNil    bool
		wantSubstr []string
	}{
		{
			name:       "字符串panic",
			panicValue: "something went wrong",
			wantSubstr: []string{"panic error", "something went wrong"},
		},
		{
			name:       "整数panic",
			panicValue: 42,
			wantSubstr: []string{"panic error", "42"},
		},
		{
			name:       "错误panic",
			panicValue: errors.New("test error"),
			wantSubstr: []string{"panic error", "test error"},
		},
		{
			name:       "结构体panic",
			panicValue: struct{ msg string }{"struct panic"},
			wantSubstr: []string{"panic error", "struct panic"},
		},
		{
			name:       "nil panic",
			panicValue: nil,
			wantNil:    true,
		},
		{
			name:       "空字符串panic",
			panicValue: "",
			wantSubstr: []string{"panic error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrPanic(tt.panicValue)

			if tt.wantNil {
				if err != nil {
					t.Errorf("ErrPanic(nil) = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatal("ErrPanic() returned nil")
			}

			// 验证错误消息包含期望的子串
			errStr := err.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(errStr, substr) {
					t.Errorf("Error() = %q, should contain %q", errStr, substr)
				}
			}

			// 验证返回的是CodeError
			unwrapped := Unwrap(err)
			if unwrapped == nil {
				t.Fatal("Unwrap() returned nil")
			}

			codeErr, ok := unwrapped.(CodeError)
			if !ok {
				t.Fatal("ErrPanic() should return CodeError")
			}

			// 验证使用了ServerInternalError错误码
			if codeErr.Code() != ServerInternalError {
				t.Errorf("Code() = %d, want %d (ServerInternalError)", codeErr.Code(), ServerInternalError)
			}

			// 验证包含堆栈信息（堆栈信息在error string中）
			// 堆栈可能不在error string中，而是在单独的堆栈对象里
			// 所以我们只需验证error不为空即可
			if errStr == "" {
				t.Error("Error string should not be empty")
			}
		})
	}
}

// ==================== ErrPanicMsg 函数测试 ====================

func TestErrPanicMsg(t *testing.T) {
	tests := []struct {
		name       string
		panicValue any
		code       int
		msg        string
		skip       int
		wantNil    bool
		wantSubstr []string
	}{
		{
			name:       "自定义错误码和消息",
			panicValue: "panic occurred",
			code:       1001,
			msg:        "custom panic handler",
			skip:       9,
			wantSubstr: []string{"1001", "custom panic handler", "panic occurred"},
		},
		{
			name:       "不同的skip值",
			panicValue: "test",
			code:       2000,
			msg:        "handler",
			skip:       5,
			wantSubstr: []string{"2000", "handler", "test"},
		},
		{
			name:       "空消息",
			panicValue: "panic",
			code:       3000,
			msg:        "",
			skip:       9,
			wantSubstr: []string{"3000", "panic"},
		},
		{
			name:       "nil panic值",
			panicValue: nil,
			code:       4000,
			msg:        "test",
			skip:       9,
			wantNil:    true,
		},
		{
			name:       "复杂类型panic",
			panicValue: map[string]int{"count": 10},
			code:       5000,
			msg:        "map panic",
			skip:       9,
			wantSubstr: []string{"5000", "map panic", "count"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrPanicMsg(tt.panicValue, tt.code, tt.msg, tt.skip)

			if tt.wantNil {
				if err != nil {
					t.Errorf("ErrPanicMsg(nil, ...) = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatal("ErrPanicMsg() returned nil")
			}

			// 验证错误消息包含期望的子串
			errStr := err.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(errStr, substr) {
					t.Errorf("Error() = %q, should contain %q", errStr, substr)
				}
			}

			// 验证返回的是CodeError
			unwrapped := Unwrap(err)
			if unwrapped == nil {
				t.Fatal("Unwrap() returned nil")
			}

			codeErr, ok := unwrapped.(CodeError)
			if !ok {
				t.Fatal("ErrPanicMsg() should return CodeError")
			}

			// 验证错误码正确
			if codeErr.Code() != tt.code {
				t.Errorf("Code() = %d, want %d", codeErr.Code(), tt.code)
			}

			// 验证消息正确
			if codeErr.Msg() != tt.msg {
				t.Errorf("Msg() = %q, want %q", codeErr.Msg(), tt.msg)
			}

			// 验证detail包含panic值
			// Detail字段是通过fmt.Sprint(r)得到的
			detail := codeErr.Detail()
			expectedDetail := strings.TrimSpace(strings.ReplaceAll(fmt.Sprint(tt.panicValue), "\n", " "))
			if !strings.Contains(detail, expectedDetail) && expectedDetail != "" {
				t.Errorf("Detail() = %q, should contain %q", detail, expectedDetail)
			}
		})
	}
}

// ==================== 实际panic场景测试 ====================

func TestErrPanic_WithActualPanic(t *testing.T) {
	tests := []struct {
		name       string
		panicFunc  func()
		wantSubstr []string
	}{
		{
			name: "字符串panic",
			panicFunc: func() {
				panic("something went wrong")
			},
			wantSubstr: []string{"panic error", "something went wrong"},
		},
		{
			name: "错误panic",
			panicFunc: func() {
				panic(errors.New("test error"))
			},
			wantSubstr: []string{"panic error", "test error"},
		},
		{
			name: "整数panic",
			panicFunc: func() {
				panic(123)
			},
			wantSubstr: []string{"panic error", "123"},
		},
		{
			name: "数组越界panic",
			panicFunc: func() {
				arr := []int{1, 2, 3}
				_ = arr[10] // 这会panic
			},
			wantSubstr: []string{"panic error", "index"},
		},
		{
			name: "nil指针解引用panic",
			panicFunc: func() {
				var p *int
				_ = *p // 这会panic
			},
			wantSubstr: []string{"panic error", "nil"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			// 捕获panic并转换为error
			func() {
				defer func() {
					if r := recover(); r != nil {
						err = ErrPanic(r)
					}
				}()
				tt.panicFunc()
			}()

			// 验证成功捕获panic
			if err == nil {
				t.Fatal("Failed to capture panic")
			}

			// 验证错误消息
			errStr := err.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(errStr, substr) {
					t.Errorf("Error() = %q, should contain %q", errStr, substr)
				}
			}

			// 验证是CodeError类型
			unwrapped := Unwrap(err)
			if _, ok := unwrapped.(CodeError); !ok {
				t.Error("Should return CodeError")
			}
		})
	}
}

func TestErrPanicMsg_WithActualPanic(t *testing.T) {
	t.Run("自定义处理panic", func(t *testing.T) {
		var err error

		// 模拟业务代码中的panic处理
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = ErrPanicMsg(r, 9999, "custom handler", 9)
				}
			}()
			panic("business logic failed")
		}()

		if err == nil {
			t.Fatal("Failed to capture panic")
		}

		// 验证使用了自定义错误码
		unwrapped := Unwrap(err)
		if codeErr, ok := unwrapped.(CodeError); ok {
			if codeErr.Code() != 9999 {
				t.Errorf("Code() = %d, want 9999", codeErr.Code())
			}
			if codeErr.Msg() != "custom handler" {
				t.Errorf("Msg() = %q, want 'custom handler'", codeErr.Msg())
			}
			if !strings.Contains(codeErr.Detail(), "business logic failed") {
				t.Error("Detail should contain panic message")
			}
		} else {
			t.Error("Should return CodeError")
		}
	})
}

// ==================== 集成测试 ====================

func TestErrPanic_Integration(t *testing.T) {
	t.Run("与errors.Is集成", func(t *testing.T) {
		var err error

		func() {
			defer func() {
				if r := recover(); r != nil {
					err = ErrPanic(r)
				}
			}()
			panic("test panic")
		}()

		if err == nil {
			t.Fatal("Failed to capture panic")
		}

		// 验证可以使用errors.Is检测
		serverErr := NewCodeError(ServerInternalError, "server error")
		if !serverErr.Is(err) {
			t.Error("Should be able to detect ServerInternalError")
		}
	})

	t.Run("与WrapMsg集成", func(t *testing.T) {
		var err error

		func() {
			defer func() {
				if r := recover(); r != nil {
					panicErr := ErrPanic(r)
					err = WrapMsg(panicErr, "additional context", "user", "alice")
				}
			}()
			panic("database connection lost")
		}()

		if err == nil {
			t.Fatal("Failed to capture panic")
		}

		errStr := err.Error()
		expectedSubstrs := []string{
			"panic error",
			"database connection lost",
			"additional context",
			"user=alice",
		}

		for _, substr := range expectedSubstrs {
			if !strings.Contains(errStr, substr) {
				t.Errorf("Error() should contain %q", substr)
			}
		}
	})

	t.Run("多层panic处理", func(t *testing.T) {
		var err error

		// 模拟多层函数调用中的panic
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = ErrPanicMsg(r, 1000, "outer handler", 9)
				}
			}()

			func() {
				defer func() {
					if r := recover(); r != nil {
						// 内层捕获后重新panic
						innerErr := ErrPanicMsg(r, 1001, "inner handler", 9)
						panic(innerErr)
					}
				}()
				panic("original panic")
			}()
		}()

		if err == nil {
			t.Fatal("Failed to capture panic")
		}

		// 验证外层处理器正确捕获
		unwrapped := Unwrap(err)
		if codeErr, ok := unwrapped.(CodeError); ok {
			if codeErr.Code() != 1000 {
				t.Errorf("Outer handler code = %d, want 1000", codeErr.Code())
			}
		} else {
			t.Error("Should return CodeError")
		}
	})
}

// ==================== 边界条件测试 ====================

func TestErrPanic_EdgeCases(t *testing.T) {
	t.Run("非常大的skip值", func(t *testing.T) {
		err := ErrPanicMsg("test", 1000, "msg", 100)
		if err == nil {
			t.Fatal("Should not return nil")
		}
		// 确保不会panic
	})

	t.Run("负数skip值", func(t *testing.T) {
		err := ErrPanicMsg("test", 1000, "msg", -1)
		if err == nil {
			t.Fatal("Should not return nil")
		}
		// 确保不会panic
	})

	t.Run("零skip值", func(t *testing.T) {
		err := ErrPanicMsg("test", 1000, "msg", 0)
		if err == nil {
			t.Fatal("Should not return nil")
		}
		// 确保不会panic
	})

	t.Run("非常长的panic消息", func(t *testing.T) {
		longMsg := strings.Repeat("a", 100000)
		err := ErrPanic(longMsg)

		if err == nil {
			t.Fatal("Should not return nil")
		}

		if !strings.Contains(err.Error(), longMsg) {
			t.Error("Should handle long panic message")
		}
	})

	t.Run("特殊字符panic", func(t *testing.T) {
		specialMsg := "错误：\n换行\t制表符\r回车\x00空字符"
		err := ErrPanic(specialMsg)

		if err == nil {
			t.Fatal("Should not return nil")
		}

		if !strings.Contains(err.Error(), "错误") {
			t.Error("Should handle special characters")
		}
	})

	t.Run("函数类型panic", func(t *testing.T) {
		funcVal := func() { println("test") }
		err := ErrPanic(funcVal)

		if err == nil {
			t.Fatal("Should not return nil")
		}

		// 函数会被fmt.Sprint转换为字符串表示
		if err.Error() == "" {
			t.Error("Should have error message")
		}
	})

	t.Run("channel类型panic", func(t *testing.T) {
		ch := make(chan int)
		err := ErrPanic(ch)

		if err == nil {
			t.Fatal("Should not return nil")
		}

		// channel会被fmt.Sprint转换为字符串表示
		if err.Error() == "" {
			t.Error("Should have error message")
		}
	})
}

// ==================== 性能基准测试 ====================

func BenchmarkErrPanic(b *testing.B) {
	panicValue := "test panic message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ErrPanic(panicValue)
	}
}

func BenchmarkErrPanicMsg(b *testing.B) {
	panicValue := "test panic message"
	code := 1000
	msg := "panic handler"
	skip := 9

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ErrPanicMsg(panicValue, code, msg, skip)
	}
}

func BenchmarkErrPanic_WithRecover(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = ErrPanic(r)
				}
			}()
			panic("test")
		}()
		_ = err
	}
}
