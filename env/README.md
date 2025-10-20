# env 包

`env` 包提供了一套简洁的环境变量读取工具，支持字符串、整数、浮点数和布尔类型的环境变量获取，并提供默认值机制。

## 特性

- ✅ **类型安全**: 支持 `string`、`int`、`float64`、`bool` 四种常用类型
- ✅ **默认值机制**: 环境变量不存在时自动返回默认值
- ✅ **错误处理**: 类型转换失败时返回详细错误信息
- ✅ **零依赖**: 仅依赖标准库和项目内部错误处理包
- ✅ **高性能**: 所有操作在纳秒级别完成，零内存分配
- ✅ **完整测试**: 100% 代码覆盖率

## 安装

```bash
go get github.com/Cospk/base-tools/env
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/Cospk/base-tools/env"
)

func main() {
    // 获取字符串环境变量
    appName := env.GetString("APP_NAME", "MyApp")
    fmt.Println("应用名称:", appName)

    // 获取整数环境变量
    port, err := env.GetInt("PORT", 8080)
    if err != nil {
        fmt.Println("端口配置错误:", err)
        return
    }
    fmt.Println("监听端口:", port)

    // 获取浮点数环境变量
    rate, err := env.GetFloat64("SUCCESS_RATE", 0.95)
    if err != nil {
        fmt.Println("比率配置错误:", err)
        return
    }
    fmt.Println("成功率:", rate)

    // 获取布尔环境变量
    debug, err := env.GetBool("DEBUG", false)
    if err != nil {
        fmt.Println("调试模式配置错误:", err)
        return
    }
    fmt.Println("调试模式:", debug)
}
```

## API 文档

### GetString

获取字符串类型的环境变量，如果不存在则返回默认值。

```go
func GetString(key, defaultValue string) string
```

**参数:**
- `key`: 环境变量名称
- `defaultValue`: 环境变量不存在时的默认值

**返回值:**
- 环境变量的值或默认值

**示例:**

```go
// 设置环境变量: export APP_NAME=MyApplication
appName := env.GetString("APP_NAME", "DefaultApp")
// 返回: "MyApplication"

// 环境变量不存在
appName := env.GetString("NON_EXISTENT", "DefaultApp")
// 返回: "DefaultApp"

// 环境变量为空字符串
// export EMPTY_VAR=
emptyValue := env.GetString("EMPTY_VAR", "default")
// 返回: "" (空字符串，而不是默认值)
```

---

### GetInt

获取整数类型的环境变量，如果不存在则返回默认值。

```go
func GetInt(key string, defaultValue int) (int, error)
```

**参数:**
- `key`: 环境变量名称
- `defaultValue`: 环境变量不存在或转换失败时的默认值

**返回值:**
- `int`: 环境变量的整数值或默认值
- `error`: 类型转换失败时返回错误，否则为 `nil`

**示例:**

```go
// 设置环境变量: export PORT=8080
port, err := env.GetInt("PORT", 3000)
if err != nil {
    log.Fatal(err)
}
// 返回: 8080, nil

// 负数
// export OFFSET=-100
offset, err := env.GetInt("OFFSET", 0)
// 返回: -100, nil

// 环境变量不存在
port, err := env.GetInt("NON_EXISTENT", 3000)
// 返回: 3000, nil

// 格式错误
// export INVALID_PORT=abc123
port, err := env.GetInt("INVALID_PORT", 3000)
// 返回: 3000, error("Atoi failed")
```

---

### GetFloat64

获取浮点数类型的环境变量，如果不存在则返回默认值。

```go
func GetFloat64(key string, defaultValue float64) (float64, error)
```

**参数:**
- `key`: 环境变量名称
- `defaultValue`: 环境变量不存在或转换失败时的默认值

**返回值:**
- `float64`: 环境变量的浮点数值或默认值
- `error`: 类型转换失败时返回错误，否则为 `nil`

**示例:**

```go
// 设置环境变量: export RATE=0.95
rate, err := env.GetFloat64("RATE", 1.0)
if err != nil {
    log.Fatal(err)
}
// 返回: 0.95, nil

// 科学计数法
// export LARGE_NUM=1.23e10
num, err := env.GetFloat64("LARGE_NUM", 0.0)
// 返回: 12300000000.0, nil

// 整数自动转换
// export COUNT=42
count, err := env.GetFloat64("COUNT", 0.0)
// 返回: 42.0, nil

// 环境变量不存在
rate, err := env.GetFloat64("NON_EXISTENT", 1.0)
// 返回: 1.0, nil

// 格式错误
// export INVALID_RATE=not_a_number
rate, err := env.GetFloat64("INVALID_RATE", 1.0)
// 返回: 1.0, error("ParseFloat failed")
```

---

### GetBool

获取布尔类型的环境变量，如果不存在则返回默认值。

```go
func GetBool(key string, defaultValue bool) (bool, error)
```

**参数:**
- `key`: 环境变量名称
- `defaultValue`: 环境变量不存在或转换失败时的默认值

**返回值:**
- `bool`: 环境变量的布尔值或默认值
- `error`: 类型转换失败时返回错误，否则为 `nil`

**支持的布尔值格式:**
- `true`: `"1"`, `"t"`, `"T"`, `"true"`, `"TRUE"`, `"True"`
- `false`: `"0"`, `"f"`, `"F"`, `"false"`, `"FALSE"`, `"False"`

**示例:**

```go
// 设置环境变量: export DEBUG=true
debug, err := env.GetBool("DEBUG", false)
if err != nil {
    log.Fatal(err)
}
// 返回: true, nil

// 使用数字
// export ENABLED=1
enabled, err := env.GetBool("ENABLED", false)
// 返回: true, nil

// 使用简写
// export VERBOSE=t
verbose, err := env.GetBool("VERBOSE", false)
// 返回: true, nil

// 大小写不敏感
// export PRODUCTION=TRUE
prod, err := env.GetBool("PRODUCTION", false)
// 返回: true, nil

// 环境变量不存在
debug, err := env.GetBool("NON_EXISTENT", false)
// 返回: false, nil

// 格式错误
// export INVALID_BOOL=yes
flag, err := env.GetBool("INVALID_BOOL", false)
// 返回: false, error("ParseBool failed")
```

## 使用场景

### 1. 应用配置管理

```go
package config

import "github.com/Cospk/base-tools/env"

type Config struct {
    AppName     string
    Port        int
    Debug       bool
    Timeout     float64
}

func Load() (*Config, error) {
    cfg := &Config{
        AppName: env.GetString("APP_NAME", "MyApp"),
    }
    
    var err error
    cfg.Port, err = env.GetInt("PORT", 8080)
    if err != nil {
        return nil, err
    }
    
    cfg.Debug, err = env.GetBool("DEBUG", false)
    if err != nil {
        return nil, err
    }
    
    cfg.Timeout, err = env.GetFloat64("TIMEOUT", 30.0)
    if err != nil {
        return nil, err
    }
    
    return cfg, nil
}
```

### 2. 数据库连接配置

```go
package database

import (
    "fmt"
    "github.com/Cospk/base-tools/env"
)

func GetDSN() (string, error) {
    host := env.GetString("DB_HOST", "localhost")
    port, err := env.GetInt("DB_PORT", 3306)
    if err != nil {
        return "", err
    }
    
    user := env.GetString("DB_USER", "root")
    password := env.GetString("DB_PASSWORD", "")
    database := env.GetString("DB_NAME", "mydb")
    
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", 
        user, password, host, port, database)
    return dsn, nil
}
```

### 3. 功能开关

```go
package features

import "github.com/Cospk/base-tools/env"

func IsFeatureEnabled(featureName string) bool {
    enabled, err := env.GetBool(featureName, false)
    if err != nil {
        // 发生错误时默认关闭功能
        return false
    }
    return enabled
}

// 使用示例
func main() {
    if IsFeatureEnabled("FEATURE_NEW_UI") {
        // 启用新 UI
    }
    
    if IsFeatureEnabled("FEATURE_BETA_API") {
        // 启用 Beta API
    }
}
```

### 4. 性能调优参数

```go
package performance

import "github.com/Cospk/base-tools/env"

type PerformanceConfig struct {
    MaxConnections    int
    ConnectionTimeout float64
    EnableCache       bool
    CacheSize         int
}

func LoadPerformanceConfig() (*PerformanceConfig, error) {
    cfg := &PerformanceConfig{}
    
    var err error
    cfg.MaxConnections, err = env.GetInt("MAX_CONNECTIONS", 100)
    if err != nil {
        return nil, err
    }
    
    cfg.ConnectionTimeout, err = env.GetFloat64("CONNECTION_TIMEOUT", 5.0)
    if err != nil {
        return nil, err
    }
    
    cfg.EnableCache, err = env.GetBool("ENABLE_CACHE", true)
    if err != nil {
        return nil, err
    }
    
    cfg.CacheSize, err = env.GetInt("CACHE_SIZE", 1000)
    if err != nil {
        return nil, err
    }
    
    return cfg, nil
}
```

## 错误处理最佳实践

### 1. 优雅降级

```go
// 如果环境变量格式错误，使用默认值并记录日志
port, err := env.GetInt("PORT", 8080)
if err != nil {
    log.Printf("警告: PORT 环境变量格式错误，使用默认值 8080: %v", err)
    port = 8080
}
```

### 2. 严格验证

```go
// 关键配置必须正确，否则终止程序
port, err := env.GetInt("PORT", 0)
if err != nil {
    log.Fatalf("错误: PORT 环境变量配置无效: %v", err)
}
if port <= 0 || port > 65535 {
    log.Fatalf("错误: PORT 必须在 1-65535 范围内，当前值: %d", port)
}
```

### 3. 批量验证

```go
func ValidateConfig() error {
    var errs []error
    
    if _, err := env.GetInt("PORT", 8080); err != nil {
        errs = append(errs, fmt.Errorf("PORT: %w", err))
    }
    
    if _, err := env.GetFloat64("TIMEOUT", 30.0); err != nil {
        errs = append(errs, fmt.Errorf("TIMEOUT: %w", err))
    }
    
    if _, err := env.GetBool("DEBUG", false); err != nil {
        errs = append(errs, fmt.Errorf("DEBUG: %w", err))
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("配置验证失败: %v", errs)
    }
    
    return nil
}
```

## 性能基准

在 Apple M2 处理器上的基准测试结果：

```
BenchmarkGetString-8    	69443581	  17.30 ns/op	  0 B/op	  0 allocs/op
BenchmarkGetInt-8       	53619103	  20.74 ns/op	  0 B/op	  0 allocs/op
BenchmarkGetFloat64-8   	31810198	  36.25 ns/op	  0 B/op	  0 allocs/op
BenchmarkGetBool-8      	63886567	  18.81 ns/op	  0 B/op	  0 allocs/op
```

**性能特点:**
- 所有操作都在纳秒级别完成
- 零内存分配，无 GC 压力
- 适合高频调用场景

## 测试

运行测试：

```bash
# 运行所有测试
go test ./env/

# 查看详细输出
go test -v ./env/

# 查看覆盖率
go test -cover ./env/

# 生成覆盖率报告
go test -coverprofile=coverage.out ./env/
go tool cover -html=coverage.out

# 运行基准测试
go test -bench=. -benchmem ./env/
```

测试覆盖率：**100%**

## 注意事项

### 1. 环境变量的优先级

环境变量的值始终优先于默认值，即使环境变量的值为空字符串：

```go
// export EMPTY=""
value := env.GetString("EMPTY", "default")
// 返回: "" (空字符串，不是 "default")
```

### 2. 类型转换错误

当环境变量存在但无法转换为目标类型时，会返回默认值和错误：

```go
// export PORT=invalid
port, err := env.GetInt("PORT", 8080)
// 返回: 8080, error
// 建议检查 err 并记录日志
```

### 3. 布尔值格式

`GetBool` 使用 Go 标准库的 `strconv.ParseBool`，只接受特定格式：
- ✅ 支持: `"1"`, `"t"`, `"T"`, `"true"`, `"TRUE"`, `"True"`, `"0"`, `"f"`, `"F"`, `"false"`, `"FALSE"`, `"False"`
- ❌ 不支持: `"yes"`, `"no"`, `"on"`, `"off"`, `"enabled"`, `"disabled"`

### 4. 线程安全

所有函数都是线程安全的，可以在并发环境中安全使用。

## 依赖

- `os`: Go 标准库，用于环境变量读取
- `strconv`: Go 标准库，用于类型转换
- `github.com/Cospk/base-tools/errs`: 项目内部错误处理包

## 许可证

本项目采用与 base-tools 相同的许可证。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

### v1.0.0 (2025-10-20)
- ✨ 初始版本发布
- ✅ 支持 String、Int、Float64、Bool 四种类型
- ✅ 完整的单元测试和基准测试
- ✅ 100% 代码覆盖率
- 📚 完整的文档和示例
