<h1 align="center" style="border-bottom: none">
    <b>Tools - Go 工具包<br>
    </b>
</h1>
<h3 align="center" style="border-bottom: none">
      ⭐️   Go 基础设施工具库  ⭐️ <br>
<h3>


<p align=center>
<a href="https://github.com/Cospk/base-tools"><img src="https://img.shields.io/github/stars/Cospk/base-tools.svg?style=flat&logo=github&colorB=deeppink&label=stars"></a>
<a href="https://github.com/Cospk/base-tools/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-green"></a>
<a href="https://pkg.go.dev/github.com/Cospk/base-tools"><img src="https://img.shields.io/badge/Language-Go-blue.svg"></a>
</p>

</p>

----

## 📖 项目概述

本项目是一个企业级 Go 基础设施工具库，为构建高性能、可扩展的分布式系统提供完整的基础设施支持。本库采用分层架构设计，提供统一的接口抽象和多实现策略，支持多云部署和微服务架构。

### 基本信息

- **项目定位**: 企业级公共工具库，为分布式系统提供基础能力支撑
- **代码规模**: 约 21,000+ 行 Go 代码
- **Go 版本**: 1.22.0+
- **许可证**: Apache 2.0


### 核心特性

- ✅ **扁平架构**: 清晰的功能模块划分，最为工具库，不需要分层架构，但是说明文档按照分层讲解便于快速理解
- ✅ **接口抽象**: 统一接口设计，支持多种实现
- ✅ **多云支持**: 支持 AWS、阿里云、腾讯云、MinIO 等多个云厂商
- ✅ **微服务友好**: 内置服务发现、链路追踪、错误传递
- ✅ **生产就绪**: 包含熔断器、限流器等可靠性保障
- ✅ **高可观测性**: 结构化日志、错误堆栈追踪

---

## 🏗️ 架构设计

### 整体分层架构

本项目采用**扁平架构**，但是具体划分从下到上可分为五个层次：

```
┌─────────────────────────────────────────────────┐
│              应用层 (Application)                │
│         (使用此工具库的业务系统)                   │
└──────┬──────────────┬──────────────┬─────────────┘
       │              │              │
       ▼              │              │
┌─────────────────────┼──────────────┼─────────────┐
│            中间件层 (Middleware)   │              │
│    mw/ | mcontext/ | tokenverify/ | a2r/        │
└──────┬──────────────┴──────────────┼─────────────┘
       │                             │
       ▼                             │
┌─────────────────────────────────────┼─────────────┐
│         集成层 (Integration)        │             │
│  db/ | mq/ | s3/ | discovery/ | stability/      │
└──────┬──────────────────────────────┴─────────────┘
       │
       ▼
┌─────────────────────────────────────────────────┐
│          基础设施层 (Infrastructure)             │
│    errs/ | log/ | config/ | env/ | utils/      │
└─────────────────────────────────────────────────┘
```

### 依赖原则

- ✅ **上层依赖下层**: 高层模块可以依赖所有低层模块（包括跨层依赖）
  - 例如：应用层可以直接调用中间件层、集成层、基础设施层
  - 例如：中间件层可以直接调用集成层、基础设施层
  - 例如：集成层可以直接调用基础设施层
- ✅ **同层解耦**: 同层模块尽量独立，减少相互依赖
- ❌ **禁止反向依赖**: 低层模块不能依赖高层模块

---

## 📦 模块分层说明

### 🏗️ 基础设施层 (Infrastructure Layer)

> 无外部依赖，被其他所有层依赖

| 模块 | 功能说明 | 文档 |
|------|---------|------|
| [errs](./errs) | 统一错误处理体系 - 错误码、堆栈追踪、错误关系链 | [📖](./errs/README.md) |
| [log](./log) | 结构化日志系统 - 基于 Zap，支持日志轮转和多输出 | [📖](./log/README.md) |
| [config](./config) | 配置管理 - YAML 解析、动态加载、配置验证 | |
| [env](./env) | 环境变量管理 - 类型安全的环境变量读取 | |
| [version](./version) | 版本管理工具 - 版本信息定义和管理 | |

**推荐学习顺序**: errs → log → config → env

### 🔌 中间件层 (Middleware Layer)

> 依赖基础设施层，提供请求处理和上下文管理

| 模块 | 功能说明 |
|------|---------|
| [mw](./mw) | gRPC/Gin 中间件 - 日志、鉴权、错误处理、责任链模式 |
| [mcontext](./mcontext) | 上下文管理 - 传递 OperationID、UserID 等请求级别信息 |
| [tokenverify](./tokenverify) | JWT 令牌验证 - Token 生成和验证 |
| [a2r](./a2r) | API 到 RPC 转换 - HTTP API 自动转换为 gRPC 调用 |

**核心功能**:
- 自动注入 OperationID 进行链路追踪
- 统一的参数校验和错误处理
- Panic 恢复机制
- 请求/响应日志记录

### 🔗 集成层 (Integration Layer)

> 第三方服务集成，提供统一抽象接口

#### 数据库 (Database)

| 模块 | 功能说明 |
|------|---------|
| [db/mongoutil](./db/mongoutil) | MongoDB 工具 - 连接池、查询构建器 |
| [db/redisutil](./db/redisutil) | Redis 工具 - 连接池、常用操作封装 |
| [db/cacheutil](./db/cacheutil) | 缓存工具 - 多级缓存、缓存穿透保护 |
| [db/pagination](./db/pagination) | 分页工具 - 统一分页查询 |

#### 消息队列 (Message Queue)

| 模块 | 功能说明 | 文档 |
|------|---------|------|
| [mq](./mq) | 消息队列 - Kafka、内存队列、简单队列的统一抽象 | [📖](./mq/README.md) |

**支持的实现**:
- Kafka - 生产环境分布式消息队列
- MemaMQ - 内存异步任务队列
- SimMQ - 简单内存消息队列

#### 对象存储 (Object Storage)

| 模块 | 功能说明 |
|------|---------|
| [s3](./s3) | 对象存储 - 多云厂商统一接口 |

**支持的云厂商**:
- AWS S3
- MinIO (S3 兼容)
- 阿里云 OSS
- 腾讯云 COS
- 七牛云 Kodo

**核心优势**:
- 云厂商无关性，避免厂商锁定
- 统一的 API 接口
- 支持分片上传、签名 URL、表单上传

#### 服务发现 (Service Discovery)

| 模块 | 功能说明 |
|------|---------|
| [discovery](./discovery) | 服务发现 - Etcd/Zookeeper/K8s/Standalone 统一接口 |

**支持的实现**:
- Etcd - 推荐生产环境
- Zookeeper - 传统服务发现
- Kubernetes - K8s 原生服务发现
- Standalone - 开发测试环境

### 🛡️ 可靠性层 (Reliability Layer)

> 提供稳定性保障和可靠性工具

| 模块 | 功能说明 |
|------|---------|
| [stability/circuitbreaker](./stability/circuitbreaker) | 熔断器 - 基于 Google SRE Breaker 算法 |
| [stability/ratelimit](./stability/ratelimit) | 限流器 - 基于 BBR 自适应限流算法 |
| [queue/task](./queue/task) | 任务队列 - 支持 Redis 和本地实现 |
| [timer](./timer) | 定时器工具 - 定时任务管理 |

**设计目标**:
- 快速失败，防止服务雪崩
- 自适应调整，无需手动配置
- 过载保护，保证服务质量

### 🔧 通用工具层 (Utility Layer)

> 提供各种通用工具函数

| 类别 | 模块 | 功能说明 |
|------|------|---------|
| **编码** | [utils/encoding](./utils/encoding) | Base64 编解码 |
| **加密** | [utils/encrypt](./utils/encrypt) | 加密工具 |
| **HTTP** | [utils/httputil](./utils/httputil) | HTTP 客户端封装 |
| **JSON** | [utils/jsonutil](./utils/jsonutil) | JSON 处理工具 |
| **网络** | [utils/network](./utils/network) | 网络工具 (IP 解析等) |
| **字符串** | [utils/stringutil](./utils/stringutil) | 字符串工具 |
| **时间** | [utils/timeutil](./utils/timeutil) | 时间工具 |
| **响应** | [apiresp](./apiresp) | API 响应封装 |
| **检查** | [checker](./checker) | 服务健康检查 |
| **文件** | [field](./field) | 文件操作工具 |
| **系统** | [system](./system) | 系统信息获取 |
| **TLS** | [xtls](./xtls) | TLS 证书管理 |

---

## 🚀 快速开始

### 安装

```bash
go get -u github.com/Cospk/base-tools
```

### 基础使用示例

#### 错误处理

```go
import "github.com/Cospk/base-tools/errs"

// 定义错误码
var ErrUserNotFound = errs.NewCodeError(20001, "用户不存在")

func GetUser(ctx context.Context, userID string) (*User, error) {
    user, err := db.FindUser(userID)
    if err != nil {
        // 包装错误，添加上下文信息和堆栈追踪
        return nil, ErrUserNotFound.WrapMsg("查询失败",
            "userId", userID,
            "error", err)
    }
    return user, nil
}
```

#### 日志记录

```go
import "github.com/Cospk/base-tools/log"

// 结构化日志
log.ZInfo(ctx, "用户登录成功", "userId", userId, "ip", clientIP)
log.ZError(ctx, "数据库查询失败", err, "table", "users", "query", query)
log.ZWarn(ctx, "缓存未命中", "key", cacheKey)

// 自适应日志级别 - 根据错误类型自动选择级别
log.ZAdaptive(ctx, "操作完成", err)
```

#### 服务发现

```go
import "github.com/Cospk/base-tools/discovery/etcd"

// 初始化服务注册中心
registry, err := etcd.NewDiscovery([]string{"localhost:2379"})

// 注册服务
err = registry.Register(ctx, "user-service", "localhost", 50001)

// 获取服务连接
conn, err := registry.GetConn(ctx, "user-service")
```

#### 对象存储 (多云支持)

```go
import "github.com/Cospk/base-tools/s3/minio"

// 初始化 MinIO (可轻松切换为 AWS/阿里云/腾讯云)
s3Client, err := minio.NewMinIO(config)

// 上传文件
result, err := s3Client.PresignedPutObject(ctx, "avatar/user123.jpg", 1*time.Hour, nil)

// 获取访问 URL
url, err := s3Client.AccessURL(ctx, "avatar/user123.jpg", 24*time.Hour, nil)
```

#### 熔断器

```go
import "github.com/Cospk/base-tools/stability/circuitbreaker"

breaker := circuitbreaker.NewSREBreaker()

func CallRemoteService() error {
    // 检查熔断器状态
    if err := breaker.Allow(); err != nil {
        return err  // 熔断器打开，快速失败
    }

    // 执行远程调用
    err := doRemoteCall()

    // 标记结果
    if err != nil {
        breaker.MarkFailed()
        return err
    }
    breaker.MarkSuccess()
    return nil
}
```

#### 消息队列

```go
import "github.com/Cospk/base-tools/mq/kafka"

// 创建生产者
producer, err := kafka.NewKafkaProducerV2(config, []string{"localhost:9092"}, "my-topic")
defer producer.Close()

// 发送消息
err = producer.SendMessage(ctx, "key1", []byte("消息内容"))

// 创建消费者
consumer, err := kafka.NewMConsumerGroupV2(ctx, config, "my-group", []string{"my-topic"}, false)
defer consumer.Close()

// 消费消息
handler := func(msg mq.Message) error {
    fmt.Printf("收到: %s\n", string(msg.Value()))
    msg.Mark()
    msg.Commit()
    return nil
}

for {
    consumer.Subscribe(ctx, handler)
}
```

---

## 📚 详细文档

### 核心模块文档

- **[errs 包文档](./errs/README.md)** - 错误处理系统完整文档
- **[log 包文档](./log/README.md)** - 日志系统完整文档
- **[mq 包文档](./mq/README.md)** - 消息队列完整文档

### 架构设计文档

- **[架构分析报告](./ARCHITECTURE_ANALYSIS.md)** - 完整的架构设计分析和实现细节
- **[包总览](./README_PACKAGES.md)** - 所有包的快速索引
- **[工具包使用指南](./TOOLS_GUIDE.md)** - env、mcontext、tokenverify 等包的使用说明
- **[包设计说明](./PACKAGE_OVERVIEW.md)** - 核心包的设计理念

---

## 🎯 核心设计理念

### 1. 接口抽象 + 多实现策略

所有外部依赖都通过接口抽象，支持多种实现：

```go
// 统一接口
type Interface interface {
    Method1()
    Method2()
}

// 多种实现
- AWS 实现
- 阿里云实现
- 腾讯云实现
- MinIO 实现
```

**优势**:
- 云厂商无关性
- 易于切换实现
- 可测试性强
- 避免厂商锁定

### 2. 统一错误处理

```go
type CodeError interface {
    Code() int                          // 错误码
    Msg() string                        // 错误消息
    Detail() string                     // 详细信息
    WithDetail(detail string) CodeError // 链式调用
    Wrap() error                        // 堆栈追踪
}
```

**特性**:
- 错误码标准化
- 自动堆栈追踪
- 跨服务错误传递
- 错误码继承关系

### 3. 中间件责任链

```go
server := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        mw.RpcServerInterceptor,  // 基础拦截
        mw.MetricsInterceptor,     // 指标收集
        mw.TracingInterceptor,     // 链路追踪
    ),
)
```

**功能**:
- 自动注入上下文
- 参数校验
- Panic 恢复
- 请求/响应日志
- 错误转换

---

## 🔄 依赖关系

### 技术栈

**核心依赖**:
- gRPC: `google.golang.org/grpc v1.71.0`
- Zap Logger: `go.uber.org/zap v1.24.0`
- Redis: `github.com/redis/go-redis/v9 v9.2.1`
- MongoDB: `go.mongodb.org/mongo-driver v1.12.0`
- Kafka: `github.com/IBM/sarama v1.43.0`
- Etcd: `go.etcd.io/etcd/client/v3 v3.5.13`

**云服务 SDK**:
- AWS S3: `github.com/aws/aws-sdk-go-v2`
- 阿里云 OSS: `github.com/aliyun/aliyun-oss-go-sdk`
- 腾讯云 COS: `github.com/tencentyun/cos-go-sdk-v5`
- 七牛云 Kodo: `github.com/qiniu/go-sdk/v7`
- MinIO: `github.com/minio/minio-go/v7`

---

## 🧪 开发指南

### 环境要求

- Go 1.22.0+
- Make
- golangci-lint

### 常用命令

```bash
# 代码格式化
make fmt

# 静态检查
make vet

# Lint 检查
make lint

# 运行测试
make test

# 测试覆盖率 (要求 ≥75%)
make cover

# 代码质量检查
make style

# 添加版权信息
make copyright-add

# 完整检查
make all
```

### 代码质量要求

- **测试覆盖率**: ≥75%
- **代码规范**: 通过 golangci-lint 检查
- **版权检查**: 所有文件包含 Apache 2.0 license header
- **Git 提交**: 符合 go-gitlint 规范

---

## 🤝 贡献指南

我们欢迎各种形式的贡献！

### 如何贡献

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 贡献方向

- 📝 完善文档和示例
- 🐛 修复 Bug
- ✨ 新增功能
- ⚡ 性能优化
- 🧪 补充测试用例

---

## 📊 项目状态

- **活跃维护**: ✅
- **生产就绪**: ✅
- **测试覆盖率**: 75%+
- **文档完善度**: 持续改进中

---

## 📄 许可证

本项目采用 Apache 2.0 许可证 - 详见 [LICENSE](./LICENSE) 文件

---

## 🔗 相关链接

- **代码仓库**: https://github.com/Cospk/base-tools
- **GitHub Issues**: https://github.com/Cospk/base-tools/issues
- **Go Package**: https://pkg.go.dev/github.com/Cospk/base-tools

---

## 💬 联系我们

- **GitHub Issues**: https://github.com/Cospk/base-tools/issues

---

<p align="center">
    如果这个项目对你有帮助，请给我们一个 ⭐️
</p>