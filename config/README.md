# Viper Config - 基于 Viper 的配置管理

`viper_config.go` 提供了一个基于 [spf13/viper](https://github.com/spf13/viper) 的配置管理封装，相比原生实现，它提供了更多开箱即用的功能。

## 为什么选择 Viper？

Viper 是 Go 生态中最流行的配置管理库，提供了：

- ✅ **多格式支持**: JSON, YAML, TOML, INI, HCL, envfile, Java properties
- ✅ **多数据源**: 文件、环境变量、命令行参数、远程配置系统
- ✅ **热加载**: 配置文件变更自动重载
- ✅ **默认值**: 支持设置默认配置
- ✅ **环境变量**: 自动绑定和覆盖
- ✅ **类型安全**: 自动类型转换
- ✅ **子配置**: 支持获取配置子集

## 安装

```bash
go get github.com/spf13/viper
go get github.com/fsnotify/fsnotify  # 用于配置热加载
```

## 快速开始

### 1. 最简单的使用方式

```go
import "github.com/Cospk/base-tools/config"

// 定义配置结构体
type AppConfig struct {
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"server"`
    Database struct {
        DSN string `mapstructure:"dsn"`
    } `mapstructure:"database"`
}

// 快速加载
var cfg AppConfig
err := config.QuickLoad("./config/app.yaml", &cfg, "APP")
if err != nil {
    log.Fatal(err)
}
```

### 2. 使用配置管理器

```go
// 创建配置管理器
vc := config.NewViperConfig(
    config.WithConfigName("app"),        // 配置文件名
    config.WithConfigType("yaml"),       // 配置格式
    config.WithConfigPath("./config"),   // 搜索路径
    config.WithEnvPrefix("APP"),         // 环境变量前缀
)

// 加载配置
if err := vc.Load(); err != nil {
    log.Fatal(err)
}

// 解析到结构体
var cfg AppConfig
if err := vc.Unmarshal(&cfg); err != nil {
    log.Fatal(err)
}
```

## 核心功能

### 1. 环境变量覆盖

环境变量会自动覆盖配置文件中的值：

```go
vc := config.NewViperConfig(
    config.WithEnvPrefix("APP"),  // 设置前缀 APP_
)

// 配置文件: server.port = 8080
// 环境变量: APP_SERVER_PORT=9090
// 最终值:    9090
port := vc.GetInt("server.port")  // 返回 9090
```

### 2. 配置热加载

监听配置文件变化并自动重新加载：

```go
vc := config.NewViperConfig(
    config.WithConfigName("app"),
    config.WithConfigPath("./config"),
)

// 加载初始配置
vc.Load()

// 注册变更回调
vc.OnConfigChange(func() {
    fmt.Println("配置已更新")
    // 重新加载配置到结构体
    var cfg AppConfig
    vc.Unmarshal(&cfg)
})

// 开始监听
vc.WatchConfig()
```

### 3. 默认值设置

```go
// 设置默认值
vc.SetDefault("server.port", 8080)
vc.SetDefault("server.timeout", 30)
vc.SetDefault("database.max_connections", 100)

// 如果配置文件和环境变量都没有设置，将使用默认值
port := vc.GetInt("server.port")
```

### 4. 动态配置

```go
// 运行时修改配置
vc.Set("feature.new_feature", true)
vc.Set("api.rate_limit", 1000)

// 保存配置到文件
err := vc.WriteConfig()  // 写入原配置文件
// 或
err := vc.WriteConfigAs("./config/new_config.yaml")  // 另存为
```

### 5. 全局配置

```go
// 初始化全局配置（通常在 main 函数中）
err := config.InitGlobalConfig(
    config.WithConfigName("app"),
    config.WithConfigPath("./config"),
    config.WithEnvPrefix("APP"),
)

// 在任何地方使用
gc := config.GetGlobalConfig()
dbHost := gc.GetString("database.host")
```

### 6. 多配置文件合并

```go
// 主配置
mainConfig := config.NewViperConfig(
    config.WithConfigName("app"),
    config.WithConfigPath("./config"),
)
mainConfig.Load()

// 环境特定配置
envConfig := config.NewViperConfig(
    config.WithConfigName("app.production"),
    config.WithConfigPath("./config"),
)
envConfig.Load()

// 合并配置
config.MergeConfig(mainConfig, envConfig)
```

## 配置文件示例

### YAML 格式 (app.yaml)

```yaml
server:
  host: localhost
  port: 8080
  timeout: 30

database:
  driver: mysql
  dsn: "user:pass@tcp(localhost:3306)/dbname?charset=utf8mb4"
  max_connections: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600

redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 10

log:
  level: info
  format: json
  output: stdout

feature:
  cache_enabled: true
  rate_limit: 1000
```

### JSON 格式 (app.json)

```json
{
  "server": {
    "host": "localhost",
    "port": 8080,
    "timeout": 30
  },
  "database": {
    "driver": "mysql",
    "dsn": "user:pass@tcp(localhost:3306)/dbname",
    "max_connections": 100
  }
}
```

### TOML 格式 (app.toml)

```toml
[server]
host = "localhost"
port = 8080
timeout = 30

[database]
driver = "mysql"
dsn = "user:pass@tcp(localhost:3306)/dbname"
max_connections = 100
```

## API 参考

### 配置选项

| 选项 | 说明 | 示例 |
|-----|------|-----|
| `WithConfigName(name)` | 配置文件名（不含扩展名） | `WithConfigName("app")` |
| `WithConfigType(type)` | 配置文件格式 | `WithConfigType("yaml")` |
| `WithConfigPath(paths...)` | 配置文件搜索路径 | `WithConfigPath(".", "./config")` |
| `WithEnvPrefix(prefix)` | 环境变量前缀 | `WithEnvPrefix("APP")` |

### 主要方法

| 方法 | 说明 | 返回值 |
|-----|------|--------|
| `Load()` | 加载配置文件 | `error` |
| `LoadWithFile(file)` | 从指定文件加载 | `error` |
| `Unmarshal(v)` | 解析到结构体 | `error` |
| `UnmarshalKey(key, v)` | 解析指定键到结构体 | `error` |
| `Get(key)` | 获取配置值 | `interface{}` |
| `GetString(key)` | 获取字符串值 | `string` |
| `GetInt(key)` | 获取整数值 | `int` |
| `GetBool(key)` | 获取布尔值 | `bool` |
| `GetFloat64(key)` | 获取浮点数值 | `float64` |
| `GetStringSlice(key)` | 获取字符串切片 | `[]string` |
| `GetStringMap(key)` | 获取字符串映射 | `map[string]interface{}` |
| `Set(key, value)` | 设置配置值 | - |
| `SetDefault(key, value)` | 设置默认值 | - |
| `IsSet(key)` | 检查键是否存在 | `bool` |
| `WatchConfig()` | 监听配置变化 | - |
| `OnConfigChange(fn)` | 注册变更回调 | - |
| `AllSettings()` | 获取所有配置 | `map[string]interface{}` |

## 配置验证系统

### 验证器架构

配置验证系统提供了多层次的验证机制：

1. **规则验证器** (`RuleValidator`) - 基于预定义规则验证
2. **结构体验证器** (`StructValidator`) - 基于结构体标签验证
3. **自定义验证器** - 实现 `ConfigValidator` 接口

### 内置验证规则

| 规则 | 说明 | 示例 |
|-----|------|------|
| `RequiredRule` | 必需字段验证 | `RequiredRule{Field: "server.host"}` |
| `RangeRule` | 数值范围验证 | `RangeRule{Field: "port", Min: 1, Max: 65535}` |
| `PatternRule` | 正则表达式验证 | `PatternRule{Field: "version", Pattern: `^\d+\.\d+\.\d+$`}` |
| `EnumRule` | 枚举值验证 | `EnumRule{Field: "level", Options: []string{"info", "warn"}}` |
| `EmailRule` | 邮箱格式验证 | `EmailRule{Field: "admin.email"}` |
| `URLRule` | URL格式验证 | `URLRule{Field: "webhook", RequireHTTPS: true}` |
| `IPRule` | IP地址验证 | `IPRule{Field: "bind_ip", Version: 4}` |
| `PortRule` | 端口号验证 | `PortRule{Field: "server.port"}` |
| `PathExistsRule` | 路径存在性验证 | `PathExistsRule{Field: "log_dir", MustBeDir: true}` |
| `CustomRule` | 自定义验证 | `CustomRule{Field: "password", ValidateFn: checkPassword}` |
| `DependencyRule` | 依赖关系验证 | `DependencyRule{Field: "ssl.cert", DependsOn: "ssl.enabled"}` |

### 使用验证系统

#### 1. 基本验证

```go
vc := NewViperConfig(
    WithConfigName("app"),
    WithConfigPath("./config"),
)

// 创建验证器
validator := NewRuleValidator(false) // false = 收集所有错误

// 添加验证规则
validator.
    AddRule(RequiredRule{Field: "server.host"}).
    AddRule(PortRule{Field: "server.port"}).
    AddRule(EmailRule{Field: "admin.email"})

// 添加到配置管理器
vc.AddValidator(validator)

// 启用自动验证（加载时自动执行）
vc.SetAutoValidate(true)

// 加载配置（会自动验证）
if err := vc.Load(); err != nil {
    log.Fatal(err)
}
```

#### 2. 结构体标签验证

```go
type Config struct {
    Host  string `validate:"required"`
    Port  int    `validate:"required,min=1,max=65535"`
    Email string `validate:"email"`
    URL   string `validate:"url"`
}

var cfg Config
vc.Unmarshal(&cfg)

// 使用结构体验证器
validator := NewStructValidator()
if err := validator.ValidateStruct(&cfg); err != nil {
    log.Fatal(err)
}
```

#### 3. 自定义验证

```go
validator.AddRule(CustomRule{
    Field: "password",
    ValidateFn: func(value interface{}) error {
        pwd := value.(string)
        if len(pwd) < 8 {
            return fmt.Errorf("密码长度至少8位")
        }
        // 添加更多检查...
        return nil
    },
})
```

#### 4. 依赖关系验证

```go
validator.AddRule(DependencyRule{
    Field:     "ssl.cert_file",
    DependsOn: "ssl.enabled",
    Condition: func(fieldValue, dependsValue interface{}) bool {
        if enabled := dependsValue.(bool); enabled {
            certFile := fieldValue.(string)
            _, err := os.Stat(certFile)
            return err == nil
        }
        return true
    },
    Message: "SSL启用时证书文件必须存在",
})
```

### 错误处理

验证系统提供了丰富的错误信息：

```go
if err := vc.Validate(); err != nil {
    switch e := err.(type) {
    case ValidationErrors:
        // 多个验证错误
        for _, verr := range e {
            if ve, ok := verr.(ValidationError); ok {
                fmt.Printf("字段: %s, 值: %v, 规则: %s, 消息: %s\n",
                    ve.Field, ve.Value, ve.Rule, ve.Message)
            }
        }
    case ValidationError:
        // 单个验证错误
        fmt.Printf("验证失败: %s\n", e.Message)
    default:
        // 其他错误
        fmt.Printf("错误: %v\n", err)
    }
}
```

### 快捷函数

```go
// 快速验证必需字段
err := ValidateRequiredFields(vc, "server.host", "server.port", "database.dsn")

// 使用规则快速验证
err := ValidateWithRules(vc,
    PortRule{Field: "server.port"},
    EmailRule{Field: "admin.email"},
)

// 验证失败时 panic
MustValidate(vc, validator)
```

## 最佳实践

### 1. 配置结构体设计

```go
type Config struct {
    // 使用 mapstructure 标签映射配置键
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
    Host    string        `mapstructure:"host"`
    Port    int           `mapstructure:"port"`
    Timeout time.Duration `mapstructure:"timeout"`
}
```

### 2. 环境特定配置

```go
env := os.Getenv("APP_ENV")
if env == "" {
    env = "development"
}

vc := config.NewViperConfig(
    config.WithConfigName(fmt.Sprintf("app.%s", env)),
    config.WithConfigPath("./config"),
)
```

### 3. 配置验证

```go
func ValidateConfig(cfg *AppConfig) error {
    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return fmt.Errorf("无效的端口: %d", cfg.Server.Port)
    }
    
    if cfg.Database.DSN == "" {
        return fmt.Errorf("数据库连接字符串不能为空")
    }
    
    return nil
}
```

### 4. 敏感信息处理

```go
// 使用环境变量存储敏感信息
vc := config.NewViperConfig(
    config.WithEnvPrefix("APP"),
)

// 绑定特定环境变量
vc.BindEnv("database.password", "DB_PASSWORD")
vc.BindEnv("api.secret", "API_SECRET")
```

### 5. 配置分层

```go
// 1. 默认值（最低优先级）
vc.SetDefault("server.port", 8080)

// 2. 配置文件
vc.Load()

// 3. 环境变量（最高优先级）
// APP_SERVER_PORT=9090

// 最终值由优先级决定
port := vc.GetInt("server.port")  // 9090
```

## 对比原生实现

| 特性 | Viper 封装 | 原生实现 |
|-----|-----------|----------|
| 实现复杂度 | 低（封装现成库） | 高（从零开始） |
| 配置格式 | 多种格式 | 需自行实现 |
| 环境变量 | 自动支持 | 需自行实现 |
| 热加载 | 内置支持 | 需自行实现 |
| 默认值 | 内置支持 | 需自行实现 |
| 类型转换 | 自动 | 需自行实现 |
| 远程配置 | 支持 | 需自行实现 |
| 性能 | 略低（功能多） | 较高（功能少） |
| 依赖 | viper + fsnotify | 仅 yaml |
| 社区支持 | 活跃 | 无 |

## 常见问题

### 1. 环境变量键名规则

环境变量名会自动转换：
- `.` 转换为 `_`
- 自动转大写
- 添加前缀

例如：`server.port` -> `APP_SERVER_PORT`

### 2. 配置文件搜索顺序

按添加的路径顺序搜索，找到第一个匹配的文件即停止。

### 3. 配置合并规则

后加载的配置会覆盖先加载的配置（相同键的情况下）。

### 4. 热加载的限制

- 仅支持文件系统的配置文件
- 不支持环境变量的热加载
- 需要文件系统支持 fsnotify

## 迁移指南

从原生实现迁移到 Viper：

```go
// 原生实现
loader := config.NewLoader(pathResolver)
err := loader.InitConfig(&cfg, "config.yaml", "./config")

// Viper 实现
err := config.QuickLoad("./config/config.yaml", &cfg, "APP")
```

## 总结

基于 Viper 的封装提供了：

1. **更少的代码**: 无需自己实现解析、验证、路径解析等
2. **更多的功能**: 热加载、环境变量、多格式等开箱即用
3. **更好的生态**: 利用成熟的社区方案，减少维护成本
4. **更高的可靠性**: 经过大量项目验证的稳定方案

适合大多数项目的配置管理需求，特别是需要复杂配置管理功能的项目。
