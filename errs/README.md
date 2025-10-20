# errs 包文档

## 概述

`errs` 包提供了一个统一的错误处理系统，支持错误码、堆栈追踪和错误关系管理。它扩展了 Go 标准库的错误处理能力，提供了更丰富的错误上下文信息和更强大的错误判断机制。

## 核心特性

### 1. 错误码支持
- 每个错误都有唯一的数字错误码
- 支持预定义的常用错误码
- 便于错误分类和国际化

### 2. 堆栈追踪
- 自动捕获错误发生时的调用栈
- 格式化输出堆栈信息，便于调试
- 可配置跳过的栈帧数量

### 3. 错误码关系管理
- 支持错误码的父子关系
- 实现错误码的层级判断
- 方便统一处理某一类错误

### 4. 上下文信息
- 支持添加键值对形式的上下文信息
- 错误包装时保留原始错误
- 支持链式调用添加信息

## 核心接口

### CodeError 接口

`CodeError` 是最核心的错误接口，扩展了标准 `error` 接口：

```go
type CodeError interface {
    Code() int                       // 返回错误码
    Msg() string                     // 返回错误消息
    Detail() string                  // 返回详细信息
    WithDetail(detail string) CodeError // 添加详细信息
    Error                            // 嵌入 Error 接口
}
```

### Error 接口

基础错误接口，提供了错误判断和包装能力：

```go
type Error interface {
    Is(err error) bool               // 判断错误类型
    Wrap() error                     // 添加堆栈追踪
    WrapMsg(msg string, kv ...any) error // 添加消息和堆栈追踪
    error
}
```

### CodeRelation 接口

错误码关系管理接口：

```go
type CodeRelation interface {
    Add(codes ...int) error        // 添加错误码关系链
    Is(parent, child int) bool     // 判断父子关系
}
```

## 使用指南

### 1. 定义错误码

```go
package main

import "github.com/Cospk/base-tools/errs"

// 定义业务错误码
const (
    ErrUserNotFoundCode = 20001
    ErrUserExistsCode   = 20002
)

// 创建错误实例
var (
    ErrUserNotFound = errs.NewCodeError(ErrUserNotFoundCode, "用户不存在")
    ErrUserExists   = errs.NewCodeError(ErrUserExistsCode, "用户已存在")
)
```

### 2. 使用预定义错误码

包中预定义了常用的错误码：

```go
// 参数错误
return errs.ErrArgs.WrapMsg("用户ID不能为空", "userId", userId)

// 权限错误
return errs.ErrNoPermission.Wrap()

// 记录不存在
return errs.ErrRecordNotFound.WrapMsg("查询失败", "table", "users", "id", id)

// 服务器内部错误
return errs.ErrInternalServer.WrapMsg("数据库连接失败", "error", err)

// Token 相关错误
return errs.ErrTokenExpired.Wrap()
return errs.ErrTokenInvalid.Wrap()
```

### 3. 包装错误并添加上下文

#### 基本包装

```go
func GetUser(userId string) (*User, error) {
    user, err := db.FindUser(userId)
    if err != nil {
        // 添加堆栈追踪
        return nil, errs.Wrap(err)
    }
    return user, nil
}
```

#### 添加消息和键值对

```go
func GetUser(userId string) (*User, error) {
    user, err := db.FindUser(userId)
    if err != nil {
        // 添加消息、键值对和堆栈追踪
        return nil, errs.WrapMsg(err, "查询用户失败",
            "userId", userId,
            "table", "users",
            "error", err)
    }
    return user, nil
}
```

#### 使用 CodeError 包装

```go
func CreateUser(user *User) error {
    err := db.Insert(user)
    if err != nil {
        // 判断是否是唯一键冲突
        if isDuplicateKeyError(err) {
            return ErrUserExists.WrapMsg("用户已存在",
                "username", user.Username)
        }
        // 其他数据库错误
        return errs.ErrInternalServer.WrapMsg("创建用户失败",
            "username", user.Username,
            "error", err)
    }
    return nil
}
```

### 4. 判断错误类型

#### 基本判断

```go
err := GetUser("123")
if err != nil {
    // 判断是否是特定错误
    if ErrUserNotFound.Is(err) {
        // 处理用户不存在的情况
        return nil, nil
    }
    // 其他错误
    return nil, err
}
```

#### 使用 errors.As 获取 CodeError

```go
err := SomeFunction()
if err != nil {
    var codeErr errs.CodeError
    if errors.As(err, &codeErr) {
        // 获取错误码
        code := codeErr.Code()
        // 获取错误消息
        msg := codeErr.Msg()
        // 获取详细信息
        detail := codeErr.Detail()
    }
}
```

### 5. 使用错误码关系

建立错误码的父子关系，便于统一处理某一类错误：

```go
const (
    ErrDatabase      = 10000  // 数据库错误（父类）
    ErrDBConnection  = 10001  // 连接错误（子类）
    ErrDBTimeout     = 10002  // 超时错误（子类）
    ErrDBQuery       = 10003  // 查询错误（子类）
)

func init() {
    // 建立错误码关系
    errs.DefaultCodeRelation.Add(ErrDatabase, ErrDBConnection, ErrDBTimeout, ErrDBQuery)
}

// 使用时可以统一判断
err := QueryDatabase()
if err != nil {
    var codeErr errs.CodeError
    if errors.As(err, &codeErr) {
        // 判断是否是任何数据库相关错误
        if errs.DefaultCodeRelation.Is(ErrDatabase, codeErr.Code()) {
            // 统一处理所有数据库错误
            log.Error("数据库操作失败", err)
            return handleDatabaseError(err)
        }
    }
}
```

### 6. 处理 Panic

```go
// 在 defer 中捕获 panic
defer func() {
    if r := recover(); r != nil {
        // 将 panic 转换为错误
        err := errs.ErrPanic(r)
        log.Error("发生panic", err)
        // 进行清理或上报
    }
}()

// 或者使用自定义错误码
defer func() {
    if r := recover(); r != nil {
        err := errs.ErrPanicMsg(r, 10500, "服务内部panic", 9)
        handleError(err)
    }
}()
```

### 7. 包装标准库错误

```go
func OpenFile(path string) (*os.File, error) {
    file, err := os.Open(path)
    if err != nil {
        // 包装标准库错误
        if os.IsNotExist(err) {
            return nil, errs.ErrRecordNotFound.WrapMsg("文件不存在",
                "path", path,
                "error", err)
        }
        return nil, errs.WrapMsg(err, "打开文件失败", "path", path)
    }
    return file, nil
}
```

## 错误输出格式

### 基本错误格式

```
错误码 错误消息 详细信息
```

示例：
```
10001 用户不存在 数据库查询失败, userId=123, table=users, error=连接超时
```

### 带堆栈追踪的格式

```
Error: 错误码 错误消息 详细信息 | -> 函数名() 文件路径:行号 -> ...
```

示例：
```
Error: 10001 用户不存在 数据库查询失败, userId=123 | -> GetUser() /app/service/user.go:45 -> HandleRequest() /app/handler/user.go:23
```

## 预定义错误码

| 错误码 | 常量名 | 说明 | HTTP 状态码 |
|-------|--------|------|-------------|
| 500 | ServerInternalError | 服务器内部错误 | 500 |
| 1001 | ArgsError | 参数错误 | 400 |
| 1002 | NoPermissionError | 权限不足 | 403 |
| 1003 | DuplicateKeyError | 重复键错误 | 409 |
| 1004 | RecordNotFoundError | 记录不存在 | 404 |
| 1501 | TokenExpiredError | Token已过期 | 401 |
| 1502 | TokenInvalidError | Token无效 | 401 |
| 1503 | TokenMalformedError | Token格式错误 | 401 |
| 1504 | TokenNotValidYetError | Token尚未生效 | 401 |
| 1505 | TokenUnknownError | Token未知错误 | 401 |
| 1506 | TokenKickedError | Token被踢出 | 401 |
| 1507 | TokenNotExistError | Token不存在 | 401 |
| 1520 | OrgUserNoPermissionError | 组织用户无权限 | 403 |

## 最佳实践

### 1. 错误码规划

建议按功能模块划分错误码范围：

```go
const (
    // 10000-10999: 用户相关错误
    ErrUserNotFound = 10001
    ErrUserExists   = 10002

    // 11000-11999: 订单相关错误
    ErrOrderNotFound = 11001
    ErrOrderCanceled = 11002

    // 12000-12999: 支付相关错误
    ErrPaymentFailed = 12001
    ErrBalanceInsufficient = 12002
)
```

### 2. 错误消息设计

- 错误消息应该简洁明了
- 详细信息使用键值对，便于日志分析
- 不要在错误消息中包含敏感信息

```go
// 好的示例
return ErrUserNotFound.WrapMsg("查询失败",
    "userId", userId,
    "operation", "GetUser")

// 不好的示例 - 消息太冗长
return ErrUserNotFound.WrapMsg("在数据库中查询用户信息时发现用户不存在，请检查用户ID是否正确",
    "userId", userId)
```

### 3. 堆栈追踪

- 只在关键位置添加堆栈追踪
- 避免重复包装错误
- 在最外层统一添加堆栈

```go
// 推荐：在最外层添加堆栈
func ServiceLayer() error {
    err := DataLayer()
    if err != nil {
        return errs.WrapMsg(err, "服务层处理失败")  // 只包装一次
    }
    return nil
}

// 不推荐：每层都包装
func ServiceLayer() error {
    err := DataLayer()
    if err != nil {
        err = errs.Wrap(err)  // 包装第一次
    }
    if err != nil {
        err = errs.Wrap(err)  // 包装第二次，多余
    }
    return err
}
```

### 4. 错误判断

- 优先使用 `CodeError.Is()` 方法
- 需要获取详细信息时使用 `errors.As`
- 利用错误码关系实现分类处理

```go
// 方式1：直接判断
if ErrUserNotFound.Is(err) {
    return handleNotFound()
}

// 方式2：获取详细信息
var codeErr errs.CodeError
if errors.As(err, &codeErr) {
    log.Info("错误码", codeErr.Code(), "详情", codeErr.Detail())
}

// 方式3：分类判断
if errs.DefaultCodeRelation.Is(ErrDatabase, codeErr.Code()) {
    return handleDatabaseError(err)
}
```

### 5. 与 HTTP API 集成

```go
func HandleError(c *gin.Context, err error) {
    var codeErr errs.CodeError
    if errors.As(err, &codeErr) {
        // 根据错误码映射 HTTP 状态码
        httpCode := mapErrorCodeToHTTP(codeErr.Code())
        c.JSON(httpCode, gin.H{
            "code":    codeErr.Code(),
            "message": codeErr.Msg(),
            "detail":  codeErr.Detail(),
        })
        return
    }

    // 未知错误
    c.JSON(500, gin.H{
        "code":    500,
        "message": "Internal Server Error",
        "detail":  err.Error(),
    })
}

func mapErrorCodeToHTTP(code int) int {
    switch code {
    case errs.ArgsError:
        return 400
    case errs.NoPermissionError:
        return 403
    case errs.RecordNotFoundError:
        return 404
    case errs.DuplicateKeyError:
        return 409
    default:
        return 500
    }
}
```

## 包结构

```
errs/
├── coderr.go       # CodeError 接口和实现
├── error.go        # Error 接口和实现
├── panic.go        # Panic 处理
├── predefine.go    # 预定义错误码
├── wrap_err.go     # 错误包装器
└── stack/
    └── stack.go    # 堆栈追踪实现
```

## 性能考虑

1. **堆栈追踪的开销**：堆栈追踪会增加一定的性能开销，但对于错误场景（通常是异常路径）来说，这个开销是可以接受的。

2. **错误码关系查询**：错误码关系使用 map 存储，查询效率为 O(1)。

3. **字符串拼接**：使用 `strings.Builder` 和预分配容量来优化字符串拼接性能。

## 常见问题

### Q1: 何时使用 Wrap 还是 WrapMsg？

- 如果只需要添加堆栈追踪，使用 `Wrap()`
- 如果需要添加额外的上下文信息，使用 `WrapMsg()`

### Q2: 如何避免堆栈信息过长？

- 只在关键边界添加堆栈（如服务层、API 层）
- 不要在每个函数中都包装错误
- 可以使用日志系统记录完整堆栈，HTTP 响应只返回简化信息

### Q3: 错误码冲突怎么办？

- 按模块划分错误码范围
- 使用常量定义错误码，避免硬编码
- 建立错误码文档，团队共享

### Q4: 如何处理第三方库的错误？

```go
err := thirdPartyLib.DoSomething()
if err != nil {
    // 包装第三方库错误
    return errs.WrapMsg(err, "第三方库调用失败",
        "library", "xxx",
        "operation", "DoSomething")
}
```

## 总结

`errs` 包提供了一个完整的错误处理解决方案，主要优势包括：

1. ✅ 统一的错误码管理
2. ✅ 完善的堆栈追踪
3. ✅ 灵活的错误关系
4. ✅ 丰富的上下文信息
5. ✅ 与标准库兼容

通过合理使用这个包，可以显著提升错误处理的规范性和调试效率。
