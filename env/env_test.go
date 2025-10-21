package env

import (
	"fmt"
	"os"
	"testing"
)

// TestGetString 测试字符串类型环境变量获取
func TestGetString(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		setEnv       bool
		defaultValue string
		want         string
	}{
		{
			name:         "环境变量存在",
			key:          "TEST_STRING",
			envValue:     "test_value",
			setEnv:       true,
			defaultValue: "default",
			want:         "test_value",
		},
		{
			name:         "环境变量不存在_返回默认值",
			key:          "TEST_STRING_NOT_EXIST",
			setEnv:       false,
			defaultValue: "default",
			want:         "default",
		},
		{
			name:         "环境变量为空字符串",
			key:          "TEST_STRING_EMPTY",
			envValue:     "",
			setEnv:       true,
			defaultValue: "default",
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理环境变量
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			got := GetString(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetInt 测试整数类型环境变量获取
func TestGetInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		setEnv       bool
		defaultValue int
		want         int
		wantErr      bool
	}{
		{
			name:         "环境变量存在_正整数",
			key:          "TEST_INT",
			envValue:     "123",
			setEnv:       true,
			defaultValue: 0,
			want:         123,
			wantErr:      false,
		},
		{
			name:         "环境变量存在_负整数",
			key:          "TEST_INT_NEGATIVE",
			envValue:     "-456",
			setEnv:       true,
			defaultValue: 0,
			want:         -456,
			wantErr:      false,
		},
		{
			name:         "环境变量不存在_返回默认值",
			key:          "TEST_INT_NOT_EXIST",
			setEnv:       false,
			defaultValue: 999,
			want:         999,
			wantErr:      false,
		},
		{
			name:         "环境变量格式错误_返回默认值和错误",
			key:          "TEST_INT_INVALID",
			envValue:     "not_a_number",
			setEnv:       true,
			defaultValue: 100,
			want:         100,
			wantErr:      true,
		},
		{
			name:         "环境变量为零",
			key:          "TEST_INT_ZERO",
			envValue:     "0",
			setEnv:       true,
			defaultValue: 100,
			want:         0,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理环境变量
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			got, err := GetInt(tt.key, tt.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetFloat64 测试浮点数类型环境变量获取
func TestGetFloat64(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		setEnv       bool
		defaultValue float64
		want         float64
		wantErr      bool
	}{
		{
			name:         "环境变量存在_正浮点数",
			key:          "TEST_FLOAT",
			envValue:     "123.456",
			setEnv:       true,
			defaultValue: 0.0,
			want:         123.456,
			wantErr:      false,
		},
		{
			name:         "环境变量存在_负浮点数",
			key:          "TEST_FLOAT_NEGATIVE",
			envValue:     "-789.012",
			setEnv:       true,
			defaultValue: 0.0,
			want:         -789.012,
			wantErr:      false,
		},
		{
			name:         "环境变量存在_科学计数法",
			key:          "TEST_FLOAT_SCIENTIFIC",
			envValue:     "1.23e10",
			setEnv:       true,
			defaultValue: 0.0,
			want:         1.23e10,
			wantErr:      false,
		},
		{
			name:         "环境变量不存在_返回默认值",
			key:          "TEST_FLOAT_NOT_EXIST",
			setEnv:       false,
			defaultValue: 99.99,
			want:         99.99,
			wantErr:      false,
		},
		{
			name:         "环境变量格式错误_返回默认值和错误",
			key:          "TEST_FLOAT_INVALID",
			envValue:     "not_a_float",
			setEnv:       true,
			defaultValue: 10.5,
			want:         10.5,
			wantErr:      true,
		},
		{
			name:         "环境变量为零",
			key:          "TEST_FLOAT_ZERO",
			envValue:     "0.0",
			setEnv:       true,
			defaultValue: 10.5,
			want:         0.0,
			wantErr:      false,
		},
		{
			name:         "环境变量为整数_自动转换",
			key:          "TEST_FLOAT_INT",
			envValue:     "42",
			setEnv:       true,
			defaultValue: 0.0,
			want:         42.0,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理环境变量
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			got, err := GetFloat64(tt.key, tt.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetBool 测试布尔类型环境变量获取
func TestGetBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		setEnv       bool
		defaultValue bool
		want         bool
		wantErr      bool
	}{
		{
			name:         "环境变量为true",
			key:          "TEST_BOOL_TRUE",
			envValue:     "true",
			setEnv:       true,
			defaultValue: false,
			want:         true,
			wantErr:      false,
		},
		{
			name:         "环境变量为false",
			key:          "TEST_BOOL_FALSE",
			envValue:     "false",
			setEnv:       true,
			defaultValue: true,
			want:         false,
			wantErr:      false,
		},
		{
			name:         "环境变量为1_表示true",
			key:          "TEST_BOOL_ONE",
			envValue:     "1",
			setEnv:       true,
			defaultValue: false,
			want:         true,
			wantErr:      false,
		},
		{
			name:         "环境变量为0_表示false",
			key:          "TEST_BOOL_ZERO",
			envValue:     "0",
			setEnv:       true,
			defaultValue: true,
			want:         false,
			wantErr:      false,
		},
		{
			name:         "环境变量为t_表示true",
			key:          "TEST_BOOL_T",
			envValue:     "t",
			setEnv:       true,
			defaultValue: false,
			want:         true,
			wantErr:      false,
		},
		{
			name:         "环境变量为f_表示false",
			key:          "TEST_BOOL_F",
			envValue:     "f",
			setEnv:       true,
			defaultValue: true,
			want:         false,
			wantErr:      false,
		},
		{
			name:         "环境变量为TRUE_大写",
			key:          "TEST_BOOL_UPPER",
			envValue:     "TRUE",
			setEnv:       true,
			defaultValue: false,
			want:         true,
			wantErr:      false,
		},
		{
			name:         "环境变量不存在_返回默认值",
			key:          "TEST_BOOL_NOT_EXIST",
			setEnv:       false,
			defaultValue: true,
			want:         true,
			wantErr:      false,
		},
		{
			name:         "环境变量格式错误_返回默认值和错误",
			key:          "TEST_BOOL_INVALID",
			envValue:     "not_a_bool",
			setEnv:       true,
			defaultValue: false,
			want:         false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理环境变量
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			got, err := GetBool(tt.key, tt.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

// BenchmarkGetString 基准测试：字符串获取性能
func BenchmarkGetString(b *testing.B) {
	os.Setenv("BENCH_STRING", "test_value")
	defer os.Unsetenv("BENCH_STRING")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetString("BENCH_STRING", "default")
	}
}

// BenchmarkGetInt 基准测试：整数获取性能
func BenchmarkGetInt(b *testing.B) {
	os.Setenv("BENCH_INT", "12345")
	defer os.Unsetenv("BENCH_INT")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := GetInt("BENCH_INT", 0); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetFloat64 基准测试：浮点数获取性能
func BenchmarkGetFloat64(b *testing.B) {
	os.Setenv("BENCH_FLOAT", "123.456")
	defer os.Unsetenv("BENCH_FLOAT")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := GetFloat64("BENCH_FLOAT", 0.0); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetBool 基准测试：布尔值获取性能
func BenchmarkGetBool(b *testing.B) {
	os.Setenv("BENCH_BOOL", "true")
	defer os.Unsetenv("BENCH_BOOL")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := GetBool("BENCH_BOOL", false); err != nil {
			b.Fatal(err)
		}
	}
}

// ExampleGetString 示例：获取字符串环境变量
func ExampleGetString() {
	os.Setenv("APP_NAME", "MyApp")
	defer os.Unsetenv("APP_NAME")

	appName := GetString("APP_NAME", "DefaultApp")
	fmt.Println(appName)
	// Output: MyApp
}

// ExampleGetInt 示例：获取整数环境变量
func ExampleGetInt() {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	port, err := GetInt("PORT", 3000)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println(port)
	// Output: 8080
}

// ExampleGetFloat64 示例：获取浮点数环境变量
func ExampleGetFloat64() {
	os.Setenv("RATE", "0.95")
	defer os.Unsetenv("RATE")

	rate, err := GetFloat64("RATE", 1.0)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println(rate)
	// Output: 0.95
}

// ExampleGetBool 示例：获取布尔环境变量
func ExampleGetBool() {
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	debug, err := GetBool("DEBUG", false)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println(debug)
	// Output: true
}
