package errs

import (
	"strings"
	"testing"
)

// ==================== 错误码常量测试 ====================

func TestErrorCodeConstants(t *testing.T) {
	tests := []struct {
		name string
		code int
		want int
	}{
		{"ServerInternalError", ServerInternalError, 500},
		{"ArgsError", ArgsError, 1001},
		{"NoPermissionError", NoPermissionError, 1002},
		{"DuplicateKeyError", DuplicateKeyError, 1003},
		{"RecordNotFoundError", RecordNotFoundError, 1004},
		{"TokenExpiredError", TokenExpiredError, 1501},
		{"TokenInvalidError", TokenInvalidError, 1502},
		{"TokenMalformedError", TokenMalformedError, 1503},
		{"TokenNotValidYetError", TokenNotValidYetError, 1504},
		{"TokenUnknownError", TokenUnknownError, 1505},
		{"TokenKickedError", TokenKickedError, 1506},
		{"TokenNotExistError", TokenNotExistError, 1507},
		{"OrgUserNoPermissionError", OrgUserNoPermissionError, 1520},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, tt.code, tt.want)
			}
		})
	}
}

// ==================== 预定义错误实例测试 ====================

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      CodeError
		wantCode int
		wantMsg  string
	}{
		{
			name:     "ErrArgs",
			err:      ErrArgs,
			wantCode: ArgsError,
			wantMsg:  "ArgsError",
		},
		{
			name:     "ErrNoPermission",
			err:      ErrNoPermission,
			wantCode: NoPermissionError,
			wantMsg:  "NoPermissionError",
		},
		{
			name:     "ErrInternalServer",
			err:      ErrInternalServer,
			wantCode: ServerInternalError,
			wantMsg:  "ServerInternalError",
		},
		{
			name:     "ErrRecordNotFound",
			err:      ErrRecordNotFound,
			wantCode: RecordNotFoundError,
			wantMsg:  "RecordNotFoundError",
		},
		{
			name:     "ErrDuplicateKey",
			err:      ErrDuplicateKey,
			wantCode: DuplicateKeyError,
			wantMsg:  "DuplicateKeyError",
		},
		{
			name:     "ErrTokenExpired",
			err:      ErrTokenExpired,
			wantCode: TokenExpiredError,
			wantMsg:  "TokenExpiredError",
		},
		{
			name:     "ErrTokenInvalid",
			err:      ErrTokenInvalid,
			wantCode: TokenInvalidError,
			wantMsg:  "TokenInvalidError",
		},
		{
			name:     "ErrTokenMalformed",
			err:      ErrTokenMalformed,
			wantCode: TokenMalformedError,
			wantMsg:  "TokenMalformedError",
		},
		{
			name:     "ErrTokenNotValidYet",
			err:      ErrTokenNotValidYet,
			wantCode: TokenNotValidYetError,
			wantMsg:  "TokenNotValidYetError",
		},
		{
			name:     "ErrTokenUnknown",
			err:      ErrTokenUnknown,
			wantCode: TokenUnknownError,
			wantMsg:  "TokenUnknownError",
		},
		{
			name:     "ErrTokenKicked",
			err:      ErrTokenKicked,
			wantCode: TokenKickedError,
			wantMsg:  "TokenKickedError",
		},
		{
			name:     "ErrTokenNotExist",
			err:      ErrTokenNotExist,
			wantCode: TokenNotExistError,
			wantMsg:  "TokenNotExistError",
		},
		{
			name:     "ErrOrgUserNoPermissionError",
			err:      ErrOrgUserNoPermissionError,
			wantCode: OrgUserNoPermissionError,
			wantMsg:  "OrgUserNoPermissionError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证错误不为nil
			if tt.err == nil {
				t.Fatal("Predefined error is nil")
			}

			// 验证错误码
			if got := tt.err.Code(); got != tt.wantCode {
				t.Errorf("Code() = %d, want %d", got, tt.wantCode)
			}

			// 验证消息
			if got := tt.err.Msg(); got != tt.wantMsg {
				t.Errorf("Msg() = %q, want %q", got, tt.wantMsg)
			}

			// 验证初始Detail为空
			if got := tt.err.Detail(); got != "" {
				t.Errorf("Detail() = %q, want empty string", got)
			}

			// 验证Error()输出格式
			errStr := tt.err.Error()
			expectedSubstrs := []string{
				strings.TrimSpace(strings.ReplaceAll(tt.wantMsg, "\n", " ")),
			}
			for _, substr := range expectedSubstrs {
				if !strings.Contains(errStr, substr) {
					t.Errorf("Error() = %q, should contain %q", errStr, substr)
				}
			}
		})
	}
}

// ==================== 预定义错误的WithDetail测试 ====================

func TestPredefinedErrors_WithDetail(t *testing.T) {
	tests := []struct {
		name       string
		baseErr    CodeError
		detail     string
		wantSubstr []string
	}{
		{
			name:       "ErrArgs with detail",
			baseErr:    ErrArgs,
			detail:     "username is required",
			wantSubstr: []string{"ArgsError", "username is required"},
		},
		{
			name:       "ErrNoPermission with detail",
			baseErr:    ErrNoPermission,
			detail:     "user is not admin",
			wantSubstr: []string{"NoPermissionError", "user is not admin"},
		},
		{
			name:       "ErrRecordNotFound with detail",
			baseErr:    ErrRecordNotFound,
			detail:     "user id: 123",
			wantSubstr: []string{"RecordNotFoundError", "user id: 123"},
		},
		{
			name:       "ErrTokenExpired with detail",
			baseErr:    ErrTokenExpired,
			detail:     "expired at 2024-01-01",
			wantSubstr: []string{"TokenExpiredError", "expired at 2024-01-01"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.baseErr.WithDetail(tt.detail)

			if err == nil {
				t.Fatal("WithDetail() returned nil")
			}

			// 验证错误消息包含所有期望的子串
			errStr := err.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(errStr, substr) {
					t.Errorf("Error() = %q, should contain %q", errStr, substr)
				}
			}

			// 验证Detail字段
			if !strings.Contains(err.Detail(), tt.detail) {
				t.Errorf("Detail() = %q, should contain %q", err.Detail(), tt.detail)
			}

			// 验证原始错误没有被修改
			if tt.baseErr.Detail() != "" {
				t.Errorf("Original error was modified, Detail() = %q", tt.baseErr.Detail())
			}
		})
	}
}

// ==================== 预定义错误的Is方法测试 ====================

func TestPredefinedErrors_Is(t *testing.T) {
	tests := []struct {
		name    string
		err1    CodeError
		err2    CodeError
		want    bool
		wantMsg string
	}{
		{
			name:    "相同错误",
			err1:    ErrArgs,
			err2:    ErrArgs,
			want:    true,
			wantMsg: "Same predefined error should match",
		},
		{
			name:    "不同错误",
			err1:    ErrArgs,
			err2:    ErrNoPermission,
			want:    false,
			wantMsg: "Different predefined errors should not match",
		},
		{
			name:    "相同错误码的新实例",
			err1:    ErrArgs,
			err2:    NewCodeError(ArgsError, "custom message"),
			want:    true,
			wantMsg: "Same error code should match",
		},
		{
			name:    "Token错误组",
			err1:    ErrTokenExpired,
			err2:    ErrTokenInvalid,
			want:    false,
			wantMsg: "Different token errors should not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err1.Is(tt.err2)
			if got != tt.want {
				t.Errorf("Is() = %v, want %v: %s", got, tt.want, tt.wantMsg)
			}
		})
	}
}

// ==================== 预定义错误与WrapMsg集成测试 ====================

func TestPredefinedErrors_WrapMsg(t *testing.T) {
	tests := []struct {
		name       string
		baseErr    CodeError
		msg        string
		kv         []any
		wantSubstr []string
	}{
		{
			name:       "ErrArgs with context",
			baseErr:    ErrArgs,
			msg:        "validation failed",
			kv:         []any{"field", "email", "value", "invalid@"},
			wantSubstr: []string{"ArgsError", "validation failed", "field=email", "value=invalid@"},
		},
		{
			name:       "ErrRecordNotFound with entity info",
			baseErr:    ErrRecordNotFound,
			msg:        "user not found",
			kv:         []any{"userID", 12345, "database", "users"},
			wantSubstr: []string{"RecordNotFoundError", "user not found", "userID=12345", "database=users"},
		},
		{
			name:       "ErrInternalServer with stack trace",
			baseErr:    ErrInternalServer,
			msg:        "database connection lost",
			kv:         nil,
			wantSubstr: []string{"ServerInternalError", "database connection lost"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.baseErr.WrapMsg(tt.msg, tt.kv...)

			if err == nil {
				t.Fatal("WrapMsg() returned nil")
			}

			// 验证错误消息包含所有期望的子串
			errStr := err.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(errStr, substr) {
					t.Errorf("Error() = %q, should contain %q", errStr, substr)
				}
			}

			// 验证可以解包到原始CodeError
			unwrapped := Unwrap(err)
			if unwrapped == nil {
				t.Fatal("Unwrap() returned nil")
			}

			if codeErr, ok := unwrapped.(CodeError); ok {
				if codeErr.Code() != tt.baseErr.Code() {
					t.Errorf("Unwrapped code = %d, want %d", codeErr.Code(), tt.baseErr.Code())
				}
			} else {
				t.Error("Unwrap() should return CodeError")
			}
		})
	}
}

// ==================== 错误码唯一性测试 ====================

func TestErrorCodeUniqueness(t *testing.T) {
	// 收集所有错误码
	codes := map[int]string{
		ServerInternalError:      "ServerInternalError",
		ArgsError:                "ArgsError",
		NoPermissionError:        "NoPermissionError",
		DuplicateKeyError:        "DuplicateKeyError",
		RecordNotFoundError:      "RecordNotFoundError",
		TokenExpiredError:        "TokenExpiredError",
		TokenInvalidError:        "TokenInvalidError",
		TokenMalformedError:      "TokenMalformedError",
		TokenNotValidYetError:    "TokenNotValidYetError",
		TokenUnknownError:        "TokenUnknownError",
		TokenKickedError:         "TokenKickedError",
		TokenNotExistError:       "TokenNotExistError",
		OrgUserNoPermissionError: "OrgUserNoPermissionError",
	}

	// 验证没有重复的错误码
	if len(codes) != 13 {
		t.Errorf("Expected 13 unique error codes, got %d", len(codes))
	}

	// 验证错误码在合理范围内
	for code, name := range codes {
		if code < 0 {
			t.Errorf("%s has negative error code: %d", name, code)
		}
		if code > 10000 {
			t.Errorf("%s has too large error code: %d", name, code)
		}
	}
}

// ==================== 实际使用场景测试 ====================

func TestPredefinedErrors_UsageScenarios(t *testing.T) {
	t.Run("参数验证场景", func(t *testing.T) {
		// 模拟参数验证失败
		validateUser := func(username string) error {
			if username == "" {
				return ErrArgs.WithDetail("username cannot be empty")
			}
			if len(username) < 3 {
				return ErrArgs.WithDetail("username must be at least 3 characters")
			}
			return nil
		}

		err := validateUser("")
		if err == nil {
			t.Fatal("Expected error for empty username")
		}

		if !ErrArgs.Is(err) {
			t.Error("Should match ErrArgs")
		}

		if !strings.Contains(err.Error(), "username cannot be empty") {
			t.Error("Should contain detail message")
		}
	})

	t.Run("数据库查询场景", func(t *testing.T) {
		// 模拟数据库查询不到记录
		findUser := func(userID int) error {
			return ErrRecordNotFound.WrapMsg("user not found", "userID", userID, "table", "users")
		}

		err := findUser(123)
		if err == nil {
			t.Fatal("Expected error for user not found")
		}

		if !ErrRecordNotFound.Is(Unwrap(err)) {
			t.Error("Should match ErrRecordNotFound")
		}

		errStr := err.Error()
		if !strings.Contains(errStr, "userID=123") {
			t.Error("Should contain userID in error message")
		}
	})

	t.Run("Token验证场景", func(t *testing.T) {
		// 模拟Token验证
		validateToken := func(token string) error {
			if token == "" {
				return ErrTokenNotExist.WithDetail("token is empty")
			}
			if token == "expired" {
				return ErrTokenExpired.WithDetail("token expired at 2024-01-01")
			}
			if token == "invalid" {
				return ErrTokenInvalid.WithDetail("signature verification failed")
			}
			return nil
		}

		// 测试空token
		err := validateToken("")
		if !ErrTokenNotExist.Is(err) {
			t.Error("Should match ErrTokenNotExist")
		}

		// 测试过期token
		err = validateToken("expired")
		if !ErrTokenExpired.Is(err) {
			t.Error("Should match ErrTokenExpired")
		}

		// 测试无效token
		err = validateToken("invalid")
		if !ErrTokenInvalid.Is(err) {
			t.Error("Should match ErrTokenInvalid")
		}
	})

	t.Run("权限检查场景", func(t *testing.T) {
		// 模拟权限检查
		checkPermission := func(userRole string) error {
			if userRole != "admin" {
				return ErrNoPermission.WrapMsg("user is not admin", "role", userRole)
			}
			return nil
		}

		err := checkPermission("guest")
		if err == nil {
			t.Fatal("Expected permission error")
		}

		unwrapped := Unwrap(err)
		if !ErrNoPermission.Is(unwrapped) {
			t.Error("Should match ErrNoPermission")
		}
	})
}

// ==================== 错误分组测试 ====================

func TestErrorCodeGroups(t *testing.T) {
	t.Run("Token错误组测试", func(t *testing.T) {
		tokenErrors := []struct {
			name string
			err  CodeError
			code int
		}{
			{"ErrTokenExpired", ErrTokenExpired, TokenExpiredError},
			{"ErrTokenInvalid", ErrTokenInvalid, TokenInvalidError},
			{"ErrTokenMalformed", ErrTokenMalformed, TokenMalformedError},
			{"ErrTokenNotValidYet", ErrTokenNotValidYet, TokenNotValidYetError},
			{"ErrTokenUnknown", ErrTokenUnknown, TokenUnknownError},
			{"ErrTokenKicked", ErrTokenKicked, TokenKickedError},
			{"ErrTokenNotExist", ErrTokenNotExist, TokenNotExistError},
		}

		// 验证所有Token错误码都在1500-1599范围内
		for _, te := range tokenErrors {
			if te.code < 1500 || te.code >= 1600 {
				t.Errorf("%s code %d is outside Token error range (1500-1599)", te.name, te.code)
			}
		}

		// 验证Token错误之间不会相互匹配
		for i, te1 := range tokenErrors {
			for j, te2 := range tokenErrors {
				if i != j && te1.err.Is(te2.err) {
					t.Errorf("%s should not match %s", te1.name, te2.name)
				}
			}
		}
	})

	t.Run("通用错误组测试", func(t *testing.T) {
		generalErrors := []struct {
			name string
			err  CodeError
			code int
		}{
			{"ErrArgs", ErrArgs, ArgsError},
			{"ErrNoPermission", ErrNoPermission, NoPermissionError},
			{"ErrDuplicateKey", ErrDuplicateKey, DuplicateKeyError},
			{"ErrRecordNotFound", ErrRecordNotFound, RecordNotFoundError},
		}

		// 验证通用错误码在1000-1099范围内
		for _, ge := range generalErrors {
			if ge.code < 1000 || ge.code >= 1100 {
				t.Errorf("%s code %d is outside general error range (1000-1099)", ge.name, ge.code)
			}
		}
	})
}
