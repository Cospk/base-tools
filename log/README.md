# log 包文档

## 概述

`log` 包提供了基于 [Uber Zap](https://github.com/uber-go/zap) 的高性能结构化日志系统。它在 Zap 的基础上进行了封装和扩展，提供了更便捷的 API 和更丰富的功能。

## 核心特性

### 1. 高性能
- 基于 Uber Zap，提供极高的日志性能
- 零内存分配（在稳定状态下）
- 支持异步日志写入

### 2. 结构化日志
- 使用键值对记录日志
- 便于日志分析和查询
- 支持多种数据类型

### 3. 日志轮转
- 支持按时间自动轮转
- 支持按大小自动轮转
- 自动清理过期日志文件

### 4. 多输出
- 同时输出到文件和控制台
- 文件输出 JSON 格式
- 控制台输出彩色格式

### 5. 上下文感知
- 自动从 context 提取信息
- 支持 OperationID、UserID 等
- 方便追踪请求链路

### 6. 消息对齐
- 日志消息自动对齐
- 提高日志可读性
- 便于快速定位信息

## 快速开始

### 1. 初始化日志系统

```go
package main

import (
    "github.com/openimsdk/tools/log"
)

func main() {
    // 初始化日志配置
    logConfig := log.Config{
        Level:           "info",          // 日志级别
        LogFile:         "/var/log/app.log", // 日志文件路径
        MaxSize:         100,             // 单个文件最大大小（MB）
        MaxBackups:      10,              // 保留的旧文件数量
        MaxAge:          30,              // 保留的最大天数
        Compress:        true,            // 是否压缩旧文件
        LocalTime:       true,            // 使用本地时间
        Console:         true,            // 是否输出到控制台
        JSONFormat:      true,            // 文件是否使用JSON格式
        WithColor:       true,            // 控制台是否使用颜色
    }

    // 初始化日志
    if err := log.InitFromConfig(&logConfig); err != nil {
        panic(err)
    }

    // 确保程序退出前刷新所有日志
    defer log.Flush()
}
```

### 2. 基本使用

```go
import (
    "context"
    "github.com/openimsdk/tools/log"
)

func Example() {
    ctx := context.Background()

    // Debug 日志
    log.ZDebug(ctx, "这是调试信息", "userId", 123, "action", "login")

    // Info 日志
    log.ZInfo(ctx, "用户登录成功", "userId", 123, "ip", "192.168.1.1")

    // Warn 日志
    err := someOperation()
    if err != nil {
        log.ZWarn(ctx, "操作警告", err, "operation", "someOperation")
    }

    // Error 日志
    if err != nil {
        log.ZError(ctx, "操作失败", err,
            "userId", 123,
            "operation", "someOperation",
            "retries", 3)
    }
}
```

## 日志级别

系统支持以下日志级别（从低到高）：

| 级别 | 方法 | 说明 | 使用场景 |
|------|------|------|----------|
| Debug | `ZDebug` | 调试信息 | 开发调试，详细的执行流程 |
| Info | `ZInfo` | 一般信息 | 正常的业务操作，如用户登录 |
| Warn | `ZWarn` | 警告信息 | 潜在问题，如缓存未命中 |
| Error | `ZError` | 错误信息 | 需要关注的错误，如数据库查询失败 |
| Panic | `ZPanic` | 严重错误 | 致命错误，如空指针 |

### 日志级别配置

```go
// 配置文件中设置
logConfig.Level = "debug"  // 或 "info", "warn", "error"

// 运行时动态调整
log.SetLevel("warn")
```

## API 使用指南

### 1. 基础日志方法

#### ZDebug - 调试日志
```go
log.ZDebug(ctx, "查询参数",
    "userId", userId,
    "pageSize", pageSize,
    "pageNum", pageNum)
```

#### ZInfo - 信息日志
```go
log.ZInfo(ctx, "用户登录成功",
    "userId", userId,
    "username", username,
    "loginTime", time.Now())
```

#### ZWarn - 警告日志
```go
log.ZWarn(ctx, "缓存未命中", err,
    "key", cacheKey,
    "fallback", "database")
```

#### ZError - 错误日志
```go
log.ZError(ctx, "数据库查询失败", err,
    "table", "users",
    "query", query,
    "params", params)
```

### 2. 上下文感知

日志系统会自动从 context 中提取以下信息：
- OperationID：操作ID，用于追踪请求链路
- UserID：用户ID
- RequestID：请求ID

```go
// 在 context 中设置信息
ctx = context.WithValue(ctx, "operationID", "op-12345")
ctx = context.WithValue(ctx, "userID", "user-67890")

// 日志会自动包含这些信息
log.ZInfo(ctx, "处理请求", "action", "updateProfile")
// 输出类似：{"level":"info","msg":"处理请求","operationID":"op-12345","userID":"user-67890","action":"updateProfile"}
```

### 3. 自适应日志

`ZAdaptive` 方法会根据错误自动选择日志级别：

```go
func DoSomething(ctx context.Context) error {
    err := operation()

    // 如果 err 为 nil，使用 Info 级别
    // 如果 err 不为 nil，使用 Error 级别
    log.ZAdaptive(ctx, "操作完成", err, "operation", "DoSomething")

    return err
}
```

### 4. 创建子 Logger

#### WithValues - 添加固定字段

```go
// 创建带有固定字段的 logger
userLogger := log.WithValues(
    "userId", userId,
    "sessionId", sessionId,
)

// 后续日志都会自动包含这些字段
userLogger.Info(ctx, "开始处理请求")
userLogger.Info(ctx, "处理完成")
```

#### WithName - 添加名称前缀

```go
// 为不同模块创建不同的 logger
dbLogger := log.WithName("database")
cacheLogger := log.WithName("cache")
apiLogger := log.WithName("api")

dbLogger.Info(ctx, "查询用户")      // [database] 查询用户
cacheLogger.Info(ctx, "缓存命中")  // [cache] 缓存命中
apiLogger.Info(ctx, "处理请求")    // [api] 处理请求
```

#### WithCallDepth - 调整调用栈深度

```go
// 在封装日志函数时使用
func MyLogWrapper(ctx context.Context, msg string) {
    // 跳过1层栈帧，显示实际调用位置
    log.WithCallDepth(1).Info(ctx, msg)
}
```

### 5. 特殊场景使用

#### Panic 处理

```go
defer func() {
    if r := recover(); r != nil {
        log.ZPanic(ctx, "发生panic",
            fmt.Errorf("panic: %v", r),
            "stack", string(debug.Stack()))
    }
}()
```

#### 性能监控

```go
func MonitorPerformance(ctx context.Context) {
    start := time.Now()

    // 执行操作
    doSomething()

    duration := time.Since(start)

    // 记录性能指标
    log.ZInfo(ctx, "性能监控",
        "operation", "doSomething",
        "duration_ms", duration.Milliseconds(),
        "status", "success")
}
```

#### 批量操作日志

```go
func BatchProcess(ctx context.Context, items []Item) {
    successCount := 0
    errorCount := 0

    for _, item := range items {
        if err := processItem(item); err != nil {
            errorCount++
            log.ZError(ctx, "处理失败", err, "itemId", item.ID)
        } else {
            successCount++
        }
    }

    // 汇总日志
    log.ZInfo(ctx, "批量处理完成",
        "total", len(items),
        "success", successCount,
        "error", errorCount)
}
```

## 日志格式

### 控制台输出格式

控制台输出使用彩色格式，便于阅读：

```
2024-01-15 10:30:45.123  INFO  用户登录成功                                      userId=123, username=alice
2024-01-15 10:30:46.456  WARN  缓存未命中                                        key=user:123, fallback=database
2024-01-15 10:30:47.789  ERROR 数据库查询失败                                  table=users, error=connection timeout
```

特点：
- 时间戳精确到毫秒
- 日志级别彩色显示
- 消息对齐到50个字符宽度
- 键值对格式清晰

### 文件输出格式

文件输出使用 JSON 格式，便于日志分析工具处理：

```json
{
  "level": "info",
  "ts": "2024-01-15T10:30:45.123+0800",
  "caller": "service/user.go:45",
  "msg": "用户登录成功",
  "operationID": "op-12345",
  "userId": 123,
  "username": "alice",
  "ip": "192.168.1.1"
}
```

## 配置详解

### 完整配置示例

```go
logConfig := log.Config{
    // 基本配置
    Level:           "info",           // 日志级别: debug, info, warn, error

    // 文件配置
    LogFile:         "/var/log/app/app.log",  // 日志文件路径
    MaxSize:         100,              // 单个日志文件最大大小（MB）
    MaxBackups:      10,               // 保留的旧日志文件数量
    MaxAge:          30,               // 保留日志文件的最大天数
    Compress:        true,             // 是否压缩旧日志文件（使用gzip）
    LocalTime:       true,             // 使用本地时间而不是UTC

    // 输出配置
    Console:         true,             // 是否同时输出到控制台
    JSONFormat:      true,             // 文件是否使用JSON格式
    WithColor:       true,             // 控制台是否使用彩色输出

    // 性能配置
    BufferSize:      256 * 1024,       // 日志缓冲区大小（字节）
}
```

### 配置项说明

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| Level | string | "info" | 日志级别 |
| LogFile | string | "" | 日志文件路径 |
| MaxSize | int | 100 | 单个文件最大大小（MB） |
| MaxBackups | int | 10 | 保留的旧文件数量 |
| MaxAge | int | 30 | 保留的最大天数 |
| Compress | bool | false | 是否压缩旧文件 |
| LocalTime | bool | true | 是否使用本地时间 |
| Console | bool | true | 是否输出到控制台 |
| JSONFormat | bool | true | 文件是否使用JSON |
| WithColor | bool | true | 控制台是否使用颜色 |

### 环境特定配置

#### 开发环境

```go
logConfig := log.Config{
    Level:      "debug",
    Console:    true,
    WithColor:  true,
    JSONFormat: false,  // 开发环境使用文本格式更易读
}
```

#### 生产环境

```go
logConfig := log.Config{
    Level:      "info",
    LogFile:    "/var/log/app/app.log",
    MaxSize:    500,    // 生产环境可以设置更大
    MaxBackups: 30,     // 保留更多备份
    MaxAge:     90,     // 保留更长时间
    Compress:   true,   // 压缩节省空间
    JSONFormat: true,   // 生产环境使用JSON便于分析
    Console:    false,  // 生产环境不输出到控制台
}
```

#### 测试环境

```go
logConfig := log.Config{
    Level:      "warn",
    Console:    true,
    WithColor:  false,  // 测试环境可能不支持彩色
}
```

## 与第三方库集成

### 与 gRPC 集成

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/grpclog"
)

func init() {
    // 将 gRPC 日志重定向到我们的日志系统
    grpclog.SetLoggerV2(log.NewGRPCLogger())
}
```

### 与 SQL 集成

```go
import (
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewDB() *gorm.DB {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: log.NewGormLogger(),  // 使用自定义的 GORM logger
    })
    return db
}
```

### 与 Gin 集成

```go
import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.New()

    // 使用自定义日志中间件
    r.Use(func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)

        log.ZInfo(c.Request.Context(), "HTTP请求",
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "status", c.Writer.Status(),
            "duration_ms", duration.Milliseconds(),
            "client_ip", c.ClientIP())
    })

    r.Run()
}
```

## 日志分析

### 使用 jq 分析 JSON 日志

```bash
# 查看所有错误日志
cat app.log | grep '"level":"error"' | jq .

# 统计各级别日志数量
cat app.log | jq -r .level | sort | uniq -c

# 查找特定用户的日志
cat app.log | jq 'select(.userId == 123)'

# 查找慢查询（超过1秒）
cat app.log | jq 'select(.duration_ms > 1000)'

# 按操作ID追踪请求链路
cat app.log | jq 'select(.operationID == "op-12345")'
```

### 使用 ELK Stack

```yaml
# Filebeat 配置示例
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/app/*.log
  json.keys_under_root: true
  json.add_error_key: true

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "app-logs-%{+yyyy.MM.dd}"
```

### 使用 Grafana Loki

```yaml
# Promtail 配置示例
clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: app
    static_configs:
    - targets:
        - localhost
      labels:
        job: app
        __path__: /var/log/app/*.log
    pipeline_stages:
    - json:
        expressions:
          level: level
          msg: msg
```

## 性能优化

### 1. 选择合适的日志级别

生产环境避免使用 Debug 级别：

```go
// 开发环境
log.SetLevel("debug")

// 生产环境
log.SetLevel("info")
```

### 2. 避免昂贵的操作

```go
// 不好：每次都执行昂贵的操作
log.ZDebug(ctx, "用户详情", "user", user.ToDetailString())  // ToDetailString() 很慢

// 好：先判断日志级别
if log.GetLevel() <= log.DebugLevel {
    log.ZDebug(ctx, "用户详情", "user", user.ToDetailString())
}
```

### 3. 批量操作时控制日志量

```go
// 不好：大量日志
for _, item := range items {
    log.ZInfo(ctx, "处理项目", "itemId", item.ID)  // 可能产生数百万条日志
}

// 好：定期记录或汇总
successCount := 0
for i, item := range items {
    if processItem(item) {
        successCount++
    }
    // 每1000条记录一次进度
    if (i+1)%1000 == 0 {
        log.ZInfo(ctx, "处理进度", "processed", i+1, "total", len(items))
    }
}
log.ZInfo(ctx, "处理完成", "total", len(items), "success", successCount)
```

### 4. 使用适当的缓冲区大小

```go
logConfig.BufferSize = 256 * 1024  // 256KB，根据实际情况调整
```

## 最佳实践

### 1. 日志内容规范

#### 消息格式

- 使用动词开头，描述发生的事情
- 保持简洁，不超过50个字符
- 使用中文或英文，保持一致

```go
// 好的示例
log.ZInfo(ctx, "用户登录成功")
log.ZError(ctx, "数据库查询失败", err)

// 不好的示例
log.ZInfo(ctx, "用户123在2024-01-15登录系统成功")  // 太冗长
log.ZError(ctx, "错误", err)  // 太简单，信息不足
```

#### 键值对规范

- 使用驼峰命名：`userId`, `operationId`
- 保持一致性：统一使用 `userId` 而不是混用 `user_id`
- 避免嵌套：直接展开对象的字段

```go
// 好的示例
log.ZInfo(ctx, "创建订单",
    "orderId", order.ID,
    "userId", order.UserID,
    "amount", order.Amount)

// 不好的示例
log.ZInfo(ctx, "创建订单",
    "order", order)  // 对象会被转换为字符串，不便于分析
```

### 2. 错误日志

错误日志应该包含足够的上下文信息：

```go
if err != nil {
    log.ZError(ctx, "更新用户失败", err,
        "userId", userId,
        "updateFields", fields,
        "previousValue", oldValue,
        "newValue", newValue)
    return err
}
```

### 3. 敏感信息处理

避免记录敏感信息：

```go
// 不要记录密码、token等敏感信息
log.ZInfo(ctx, "用户登录",
    "username", username)
    // "password", password)  // 绝不记录密码

// 可以记录部分信息
log.ZInfo(ctx, "Token验证",
    "tokenPrefix", token[:10])  // 只记录token前缀
```

### 4. 日志采样

对于高频操作，考虑采样：

```go
var logCounter int64

func HighFrequencyOperation(ctx context.Context) {
    // 每100次记录一次日志
    if atomic.AddInt64(&logCounter, 1)%100 == 0 {
        log.ZInfo(ctx, "高频操作", "count", logCounter)
    }
}
```

### 5. 结构化数据

对于复杂数据，考虑使用结构化方式：

```go
log.ZInfo(ctx, "订单统计",
    "date", date,
    "total_orders", totalOrders,
    "total_amount", totalAmount,
    "avg_amount", avgAmount,
    "top_product", topProduct)
```

## 故障排查

### 常见问题

#### 1. 日志文件没有生成

检查：
- 文件路径是否正确
- 目录是否存在且有写权限
- 配置的 `LogFile` 是否为空

#### 2. 日志不输出

检查：
- 日志级别是否设置正确
- 是否调用了 `log.Flush()`
- 程序是否正常退出

#### 3. 日志文件过大

解决方案：
- 调整 `MaxSize` 配置
- 减少 `MaxBackups` 和 `MaxAge`
- 启用压缩：`Compress: true`
- 降低日志级别

#### 4. 性能问题

优化建议：
- 提高日志级别（info -> warn -> error）
- 减少日志记录频率
- 增大缓冲区大小
- 避免在日志中进行昂贵操作

## 迁移指南

### 从标准库 log 迁移

```go
// 标准库
log.Printf("用户 %s 登录成功", username)

// 迁移后
log.ZInfo(ctx, "用户登录成功", "username", username)
```

### 从其他日志库迁移

```go
// logrus
logrus.WithFields(logrus.Fields{
    "userId": userId,
}).Info("用户登录成功")

// 迁移后
log.ZInfo(ctx, "用户登录成功", "userId", userId)

// zap
zap.L().Info("用户登录成功",
    zap.String("userId", userId))

// 迁移后
log.ZInfo(ctx, "用户登录成功", "userId", userId)
```

## 总结

`log` 包提供了一个完整的日志解决方案，主要优势包括：

1. ✅ 基于 Zap 的高性能
2. ✅ 结构化日志，便于分析
3. ✅ 自动日志轮转和清理
4. ✅ 多输出支持
5. ✅ 上下文感知
6. ✅ 易于集成和使用

通过合理使用这个包，可以显著提升日志系统的性能和可维护性。
