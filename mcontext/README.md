# mcontext 包

`mcontext` 包提供了 `context.Context` 的扩展功能，用于在请求链路中传递操作ID、用户ID、平台信息等关键上下文信息，特别适用于分布式系统的链路追踪和用户行为跟踪。

## 特性

- ✅ **链路追踪**: 支持 OperationID 在整个请求链路中传递
- ✅ **用户信息**: 携带用户ID、平台等用户相关信息
- ✅ **连接管理**: 支持连接ID的传递和管理
- ✅ **触发器支持**: 支持触发器ID的传递
- ✅ **类型安全**: 使用常量键避免键冲突
- ✅ **简洁API**: 提供便捷的 Get/Set 方法
- ✅ **批量操作**: 支持批量获取和设置上下文信息

## 安装

```bash
go get github.com/Cospk/base-tools/mcontext
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "github.com/Cospk/base-tools/mcontext"
)

func main() {
    // 创建带有 OperationID 的新 context
    ctx := mcontext.NewCtx("op-123456")
    
    // 设置用户信息
    ctx = mcontext.SetOpUserID(ctx, "user-001")
    ctx = mcontext.WithOpUserPlatformContext(ctx, "iOS")
    
    // 设置连接ID
    ctx = mcontext.SetConnID(ctx, "conn-789")
    
    // 获取信息
    operationID := mcontext.GetOperationID(ctx)
    userID := mcontext.GetOpUserID(ctx)
    platform := mcontext.GetOpUserPlatform(ctx)
    connID := mcontext.GetConnID(ctx)
    
    fmt.Printf("OperationID: %s\n", operationID)
    fmt.Printf("UserID: %s\n", userID)
    fmt.Printf("Platform: %s\n", platform)
    fmt.Printf("ConnID: %s\n", connID)
}
```

## API 文档

### 创建 Context

#### NewCtx
创建新的 context 并设置 operationID，用于链路追踪。

```go
func NewCtx(operationID string) context.Context
```

**示例:**
```go
ctx := mcontext.NewCtx("op-123456")
```

#### WithMustInfoCtx
从字符串切片创建包含必需信息的 context。

```go
func WithMustInfoCtx(values []string) context.Context
```

**示例:**
```go
// 按顺序: OperationID, OpUserID, OpUserPlatform, ConnID
values := []string{"op-123", "user-001", "Android", "conn-456"}
ctx := mcontext.WithMustInfoCtx(values)
```

### 设置信息

#### SetOperationID
设置操作ID。

```go
func SetOperationID(ctx context.Context, operationID string) context.Context
```

#### SetOpUserID
设置操作用户ID。

```go
func SetOpUserID(ctx context.Context, opUserID string) context.Context
```

#### SetConnID
设置连接ID。

```go
func SetConnID(ctx context.Context, connID string) context.Context
```

#### WithOpUserIDContext
设置操作用户ID到 context（别名方法）。

```go
func WithOpUserIDContext(ctx context.Context, opUserID string) context.Context
```

#### WithOpUserPlatformContext
设置用户平台到 context。

```go
func WithOpUserPlatformContext(ctx context.Context, platform string) context.Context
```

#### WithTriggerIDContext
设置触发器ID到 context。

```go
func WithTriggerIDContext(ctx context.Context, triggerID string) context.Context
```

### 获取信息

#### GetOperationID
从 context 获取操作ID。

```go
func GetOperationID(ctx context.Context) string
```

#### GetOpUserID
从 context 获取用户ID。

```go
func GetOpUserID(ctx context.Context) string
```

#### GetConnID
从 context 获取连接ID。

```go
func GetConnID(ctx context.Context) string
```

#### GetTriggerID
从 context 获取触发器ID。

```go
func GetTriggerID(ctx context.Context) string
```

#### GetOpUserPlatform
从 context 获取用户平台。

```go
func GetOpUserPlatform(ctx context.Context) string
```

#### GetRemoteAddr
从 context 获取远程地址。

```go
func GetRemoteAddr(ctx context.Context) string
```

### 批量获取

#### GetMustCtxInfo
获取必需的上下文信息，如果缺少任何字段则返回错误。

```go
func GetMustCtxInfo(ctx context.Context) (operationID, opUserID, platform, connID string, err error)
```

**示例:**
```go
operationID, userID, platform, connID, err := mcontext.GetMustCtxInfo(ctx)
if err != nil {
    log.Printf("缺少必需的上下文信息: %v", err)
    return
}
```

#### GetCtxInfos
获取上下文信息，只有 operationID 是必需的。

```go
func GetCtxInfos(ctx context.Context) (operationID, opUserID, platform, connID string, err error)
```

**示例:**
```go
operationID, userID, platform, connID, err := mcontext.GetCtxInfos(ctx)
if err != nil {
    // 只有当 operationID 缺失时才会返回错误
    log.Printf("缺少 operationID: %v", err)
    return
}
// userID、platform、connID 可能为空字符串
```

## 使用场景

### 1. 分布式链路追踪

```go
package service

import (
    "context"
    "github.com/Cospk/base-tools/mcontext"
    "github.com/Cospk/base-tools/log"
)

func HandleRequest(ctx context.Context, req *Request) (*Response, error) {
    // 创建带有追踪ID的 context
    ctx = mcontext.NewCtx(generateOperationID())
    
    // 设置用户信息
    ctx = mcontext.SetOpUserID(ctx, req.UserID)
    ctx = mcontext.WithOpUserPlatformContext(ctx, req.Platform)
    
    // 记录请求开始
    log.Info(ctx, "开始处理请求", 
        "operationID", mcontext.GetOperationID(ctx),
        "userID", mcontext.GetOpUserID(ctx))
    
    // 调用下游服务
    result, err := callDownstreamService(ctx, req)
    if err != nil {
        log.Error(ctx, "调用下游服务失败", err)
        return nil, err
    }
    
    return result, nil
}

func callDownstreamService(ctx context.Context, req *Request) (*Response, error) {
    // 从 context 中获取追踪信息
    operationID := mcontext.GetOperationID(ctx)
    
    // 将追踪信息传递给下游服务
    headers := map[string]string{
        "X-Operation-ID": operationID,
        "X-User-ID": mcontext.GetOpUserID(ctx),
    }
    
    // 发起请求...
    return doRequest(headers, req)
}
```

### 2. WebSocket 连接管理

```go
package websocket

import (
    "context"
    "github.com/Cospk/base-tools/mcontext"
)

type Connection struct {
    ID     string
    UserID string
    ctx    context.Context
}

func NewConnection(userID, connID string) *Connection {
    ctx := context.Background()
    ctx = mcontext.SetOpUserID(ctx, userID)
    ctx = mcontext.SetConnID(ctx, connID)
    
    return &Connection{
        ID:     connID,
        UserID: userID,
        ctx:    ctx,
    }
}

func (c *Connection) HandleMessage(msg []byte) error {
    // 为每个消息创建操作ID
    ctx := mcontext.SetOperationID(c.ctx, generateMessageID())
    
    log.Info(ctx, "收到消息",
        "connID", mcontext.GetConnID(ctx),
        "userID", mcontext.GetOpUserID(ctx))
    
    // 处理消息...
    return processMessage(ctx, msg)
}
```

### 3. 中间件集成

```go
package middleware

import (
    "net/http"
    "github.com/Cospk/base-tools/mcontext"
)

func TraceMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 从请求头获取或生成追踪ID
        operationID := r.Header.Get("X-Operation-ID")
        if operationID == "" {
            operationID = generateOperationID()
        }
        
        // 创建带追踪信息的 context
        ctx := mcontext.NewCtx(operationID)
        
        // 从请求中提取用户信息
        userID := r.Header.Get("X-User-ID")
        if userID != "" {
            ctx = mcontext.SetOpUserID(ctx, userID)
        }
        
        platform := r.Header.Get("X-Platform")
        if platform != "" {
            ctx = mcontext.WithOpUserPlatformContext(ctx, platform)
        }
        
        // 将 context 传递给下一个处理器
        r = r.WithContext(ctx)
        next(w, r)
    }
}
```

### 4. 批量操作场景

```go
package batch

import (
    "context"
    "github.com/Cospk/base-tools/mcontext"
)

func ProcessBatch(ctx context.Context, items []Item) error {
    // 验证必需的上下文信息
    operationID, userID, platform, connID, err := mcontext.GetMustCtxInfo(ctx)
    if err != nil {
        return fmt.Errorf("缺少必需的上下文信息: %w", err)
    }
    
    log.Info(ctx, "开始批量处理",
        "operationID", operationID,
        "userID", userID,
        "platform", platform,
        "connID", connID,
        "itemCount", len(items))
    
    for i, item := range items {
        // 为每个项目创建子操作ID
        itemCtx := mcontext.SetOperationID(ctx, 
            fmt.Sprintf("%s-item-%d", operationID, i))
        
        if err := processItem(itemCtx, item); err != nil {
            log.Error(itemCtx, "处理项目失败", err)
            // 根据业务需求决定是否继续
        }
    }
    
    return nil
}
```

### 5. 触发器场景

```go
package trigger

import (
    "context"
    "github.com/Cospk/base-tools/mcontext"
)

type Trigger struct {
    ID   string
    Name string
}

func (t *Trigger) Execute(ctx context.Context) error {
    // 设置触发器ID
    ctx = mcontext.WithTriggerIDContext(ctx, t.ID)
    
    log.Info(ctx, "触发器开始执行",
        "triggerID", mcontext.GetTriggerID(ctx),
        "triggerName", t.Name)
    
    // 执行触发器逻辑
    if err := t.runActions(ctx); err != nil {
        log.Error(ctx, "触发器执行失败", err)
        return err
    }
    
    return nil
}
```

## 最佳实践

### 1. 始终传递 Context

```go
// ✅ 好的做法
func ProcessOrder(ctx context.Context, order *Order) error {
    ctx = mcontext.SetOperationID(ctx, generateOrderID())
    // 使用 ctx...
}

// ❌ 不好的做法
func ProcessOrder(order *Order) error {
    ctx := context.Background() // 丢失了上游的追踪信息
    // ...
}
```

### 2. 在服务入口设置追踪信息

```go
func main() {
    http.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
        // 在入口处设置追踪信息
        ctx := mcontext.NewCtx(generateRequestID())
        ctx = mcontext.SetOpUserID(ctx, extractUserID(r))
        ctx = mcontext.WithOpUserPlatformContext(ctx, r.Header.Get("X-Platform"))
        
        // 传递给业务逻辑
        handleUserRequest(ctx, w, r)
    })
}
```

### 3. 验证关键信息

```go
func CriticalOperation(ctx context.Context) error {
    // 对于关键操作，验证必需的上下文信息
    _, _, _, _, err := mcontext.GetMustCtxInfo(ctx)
    if err != nil {
        return fmt.Errorf("关键操作需要完整的上下文信息: %w", err)
    }
    
    // 执行关键操作...
    return nil
}
```

### 4. 日志集成

```go
func LogWithContext(ctx context.Context, level, message string, fields ...interface{}) {
    // 自动添加上下文信息到日志
    baseFields := []interface{}{
        "operationID", mcontext.GetOperationID(ctx),
        "userID", mcontext.GetOpUserID(ctx),
        "platform", mcontext.GetOpUserPlatform(ctx),
    }
    
    allFields := append(baseFields, fields...)
    log.Log(level, message, allFields...)
}
```

## 注意事项

1. **Context 不可变性**: Context 是不可变的，每次设置值都会返回新的 context
2. **值类型安全**: 所有值都以字符串形式存储，避免类型断言失败
3. **空值处理**: Get 方法对于不存在的值会返回空字符串，不会 panic
4. **必需字段验证**: 使用 `GetMustCtxInfo` 时会验证所有必需字段，缺少任何字段都会返回错误

## 依赖

- `context`: Go 标准库
- `github.com/Cospk/base-tools/errs`: 错误处理包
- `github.com/Cospk/base-tools/utils/constants`: 常量定义包

## 许可证

本项目采用与 base-tools 相同的许可证。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

### v1.0.0 (2025-10-20)
- ✨ 初始版本发布
- ✅ 支持 OperationID、UserID、Platform、ConnID、TriggerID 等字段
- ✅ 提供批量获取和验证方法
- 📚 完整的文档和示例
