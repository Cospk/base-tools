# Go语言测试完全指南

> **一份从零到精通的Go测试实战手册**
>
> 涵盖：为什么测试？测什么？怎么测？如何写好测试？

---

## 📚 目录

### Part 1: 基础篇
- [1.1 什么是测试](#11-什么是测试)
- [1.2 为什么需要测试](#12-为什么需要测试)
- [1.3 测试类型详解](#13-测试类型详解)
- [1.4 Go测试基础](#14-go测试基础)

### Part 2: 实战篇
- [2.1 单元测试实战](#21-单元测试实战)
- [2.2 集成测试实战](#22-集成测试实战)
- [2.3 基准测试实战](#23-基准测试实战)
- [2.4 表驱动测试模式](#24-表驱动测试模式)

### Part 3: 进阶篇
- [3.1 如何写好测试](#31-如何写好测试)
- [3.2 测试的7个原则](#32-测试的7个原则)
- [3.3 实战技巧集](#33-实战技巧集)
- [3.4 测试覆盖率](#34-测试覆盖率)

### Part 4: 深度理解
- [4.1 为什么基础库必须有大量测试](#41-为什么基础库必须有大量测试)
- [4.2 为什么测试是重构的信心来源](#42-为什么测试是重构的信心来源)
- [4.3 测试ROI分析](#43-测试roi分析)

### Part 5: 实践指南
- [5.1 针对base-tools的测试示例](#51-针对base-tools的测试示例)
- [5.2 测试最佳实践](#52-测试最佳实践)
- [5.3 测试驱动开发(TDD)](#53-测试驱动开发tdd)
- [5.4 从今天开始行动](#54-从今天开始行动)

---

# Part 1: 基础篇

## 1.1 什么是测试

测试是验证代码是否按预期工作的过程。通过编写测试代码来自动化验证功能的正确性，而不是每次都手动测试。

**简单理解**：
```
手动测试：改代码 → 启动程序 → 点点点 → 看结果 → 重复
自动测试：改代码 → 运行测试 → 2秒看结果 ✅
```

## 1.2 为什么需要测试

### 1. 保证代码质量
- 确保代码按预期工作
- 捕获bug和边界情况
- 防止回归(regression)问题

### 2. 重构信心
- 修改代码时能快速发现破坏性变更
- 安全地优化和重构代码

### 3. 文档作用
- 测试代码展示了API的使用方式
- 说明了预期行为和边界条件

### 4. 开发效率
- 减少手动测试时间
- 快速定位问题
- 支持持续集成/持续部署(CI/CD)

### 5. 团队协作
- 其他开发者修改代码时有安全网
- 新成员快速了解代码行为

## 1.3 测试类型详解

### 单元测试 (Unit Test)

**目的**: 测试单个函数或方法的行为

**特点**:
- 测试范围小，只测试一个独立单元
- 运行速度快（微秒级）
- 不依赖外部系统(数据库、网络等)
- 使用mock/stub隔离依赖

**使用场景**:
- ✅ 测试工具函数
- ✅ 测试业务逻辑
- ✅ 测试数据转换
- ✅ 测试边界条件和错误处理

**示例**:
```go
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

### 集成测试 (Integration Test)

**目的**: 测试多个组件协同工作的情况

**特点**:
- 测试范围较大
- 可能涉及真实的外部依赖
- 运行速度较慢
- 验证组件间的交互

**使用场景**:
- ✅ 测试数据库操作
- ✅ 测试API端点
- ✅ 测试消息队列交互
- ✅ 测试文件系统操作

**示例**:
```go
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    repo := NewUserRepository(db)
    user := &User{Name: "测试用户"}

    err := repo.Create(user)
    if err != nil {
        t.Fatalf("创建用户失败: %v", err)
    }

    // 验证数据真的写入了
    found, _ := repo.FindByID(user.ID)
    if found.Name != user.Name {
        t.Error("用户数据不匹配")
    }
}
```

### 基准测试 (Benchmark Test)

**目的**: 测量代码的性能表现

**特点**:
- 测量执行时间
- 测量内存分配
- 对比不同实现的性能
- 发现性能瓶颈

**使用场景**:
- ✅ 优化关键路径代码
- ✅ 对比不同算法
- ✅ 验证性能改进
- ✅ 监控性能退化

**示例**:
```go
func BenchmarkStringConcatenation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = "hello" + "world"
    }
}
```

### 测试金字塔

```
        /\
       /  \    E2E测试(少量,慢,昂贵)
      /____\
     /      \  集成测试(适量,较慢)
    /________\
   /          \ 单元测试(大量,快速,便宜)
  /____________\
```

## 1.4 Go测试基础

### 文件命名规则
```
foo.go       ←→  foo_test.go
bar.go       ←→  bar_test.go
utils.go     ←→  utils_test.go
```

测试文件必须以 `_test.go` 结尾

### 函数命名规则
```go
// 单元测试
func TestXxx(t *testing.T) { ... }

// 基准测试
func BenchmarkXxx(b *testing.B) { ... }

// 示例(出现在文档中)
func ExampleXxx() { ... }

// 模糊测试(Go 1.18+)
func FuzzXxx(f *testing.F) { ... }
```

### 常用命令
```bash
# 运行当前目录测试
go test

# 运行所有子目录测试
go test ./...

# 运行特定测试
go test -run TestFunctionName

# 详细输出
go test -v

# 运行基准测试
go test -bench=.

# 测试覆盖率
go test -cover

# 生成覆盖率报告
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### testing.T 常用方法
```go
// 错误但继续
t.Error("错误信息")
t.Errorf("格式化: %v", err)

// 错误并停止
t.Fatal("致命错误")
t.Fatalf("致命: %v", err)

// 日志(-v时显示)
t.Log("日志")
t.Logf("格式化: %v", value)

// 跳过测试
t.Skip("跳过")
t.Skipf("跳过: %s", reason)

// 并行测试
t.Parallel()

// 子测试
t.Run("子测试", func(t *testing.T) {
    // ...
})

// 辅助函数标记
t.Helper()

// 清理
t.Cleanup(func() {
    // 清理代码
})
```

---

# Part 2: 实战篇

## 2.1 单元测试实战

### AAA模式 (Arrange-Act-Assert)

```go
func TestCalculatePrice(t *testing.T) {
    // Arrange: 准备测试数据
    items := []Item{
        {Name: "商品A", Price: 100},
        {Name: "商品B", Price: 200},
    }
    discount := 0.1

    // Act: 执行被测试的函数
    result := CalculatePrice(items, discount)

    // Assert: 验证结果
    expected := 270.0
    if result != expected {
        t.Errorf("CalculatePrice() = %.2f; want %.2f", result, expected)
    }
}
```

### 测试错误处理

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name      string
        a, b      float64
        want      float64
        expectErr bool
    }{
        {"正常除法", 10, 2, 5, false},
        {"除以零", 10, 0, 0, true},
        {"负数", -10, 2, -5, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)

            if tt.expectErr {
                if err == nil {
                    t.Error("期望错误但没有返回")
                }
                return
            }

            if err != nil {
                t.Fatalf("不期望错误: %v", err)
            }

            if got != tt.want {
                t.Errorf("got %.2f, want %.2f", got, tt.want)
            }
        })
    }
}
```

## 2.2 集成测试实战

### 数据库测试

```go
func TestUserRepository_Integration(t *testing.T) {
    // 跳过快速测试模式
    if testing.Short() {
        t.Skip("跳过集成测试")
    }

    // 设置
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    repo := NewUserRepository(db)

    t.Run("创建并查询", func(t *testing.T) {
        user := &User{Name: "张三", Email: "test@example.com"}

        // 创建
        err := repo.Create(user)
        if err != nil {
            t.Fatalf("创建失败: %v", err)
        }

        // 查询
        found, err := repo.FindByEmail("test@example.com")
        if err != nil {
            t.Fatalf("查询失败: %v", err)
        }

        if found.Name != user.Name {
            t.Errorf("名称不匹配")
        }
    })
}

// 辅助函数
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, _ := sql.Open("sqlite3", ":memory:")
    db.Exec(`CREATE TABLE users (...)`)
    return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
    t.Helper()
    db.Close()
}
```

### 运行集成测试

```bash
# 跳过集成测试(快速)
go test -short

# 运行所有测试
go test

# 只运行集成测试
go test -run Integration
```

## 2.3 基准测试实战

### 基本用法

```go
func BenchmarkStringBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var builder strings.Builder
        for j := 0; j < 100; j++ {
            builder.WriteString("hello")
        }
        _ = builder.String()
    }
}

func BenchmarkStringConcat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        s := ""
        for j := 0; j < 100; j++ {
            s += "hello"
        }
        _ = s
    }
}
```

### 运行基准测试

```bash
# 运行所有基准测试
go test -bench=.

# 显示内存分配
go test -bench=. -benchmem

# 运行更长时间获得准确结果
go test -bench=. -benchtime=10s

# 对比结果
go test -bench=. > old.txt
# 修改代码后
go test -bench=. > new.txt
benchstat old.txt new.txt
```

### 输出解读

```
BenchmarkStringBuilder-8   2000000   789 ns/op   512 B/op   1 allocs/op
```

- `StringBuilder-8`: 测试名-CPU核心数
- `2000000`: 运行次数(b.N)
- `789 ns/op`: 每次操作耗时
- `512 B/op`: 每次分配内存
- `1 allocs/op`: 每次内存分配次数

### 排除准备时间

```go
func BenchmarkWithSetup(b *testing.B) {
    // 准备(不计时)
    data := generateLargeData()

    // 重置计时器
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        process(data)
    }
}
```

## 2.4 表驱动测试模式

**Go社区推荐的最佳实践**

### 完整示例

```go
func TestFormatString(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        maxLen    int
        leftAlign bool
        want      string
    }{
        {
            name:      "短字符串左对齐",
            input:     "hello",
            maxLen:    10,
            leftAlign: true,
            want:      "hello     ",
        },
        {
            name:      "短字符串右对齐",
            input:     "hello",
            maxLen:    10,
            leftAlign: false,
            want:      "     hello",
        },
        {
            name:      "长字符串截断",
            input:     "hello world",
            maxLen:    5,
            leftAlign: true,
            want:      "hello",
        },
        {
            name:      "空字符串",
            input:     "",
            maxLen:    5,
            leftAlign: true,
            want:      "     ",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := FormatString(tt.input, tt.maxLen, tt.leftAlign)
            if got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

### 优势

✅ 代码简洁
✅ 易于添加新场景
✅ 每个场景独立运行
✅ 失败时能看到具体场景

---

# Part 3: 进阶篇

## 3.1 如何写好测试

### F.I.R.S.T原则

好测试的5个特征：

#### F - Fast (快速)
- 单个测试 < 100ms
- 所有测试 < 10秒
- 开发者愿意频繁运行

#### I - Independent (独立)
- 测试间不依赖
- 顺序无关
- 可以单独运行

#### R - Repeatable (可重复)
- 每次运行结果相同
- 不依赖外部环境
- 不依赖时间/随机数

#### S - Self-Validating (自验证)
- 测试自己检查结果
- 不需要人工判断
- 明确的通过/失败

#### T - Timely (及时)
- 和代码同时编写
- 最好先写测试(TDD)
- 不拖到最后补测试

### 好测试 vs 坏测试

#### ✅ 好测试

```go
func TestAdd(t *testing.T) {
    result := Add(2, 3)

    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
// ✅ 快速、独立、可重复、自验证
```

#### ❌ 坏测试

```go
func TestUserService(t *testing.T) {
    // ❌ 慢: 连接生产数据库
    db := connectToProductionDB()

    // ❌ 不独立: 依赖其他测试数据
    user := db.FindByEmail("test@example.com")

    // ❌ 不可重复: 依赖当前时间
    if user.CreatedAt.After(time.Now()) {
        t.Error("时间错误")
    }

    // ❌ 不自验证: 需要人工检查
    fmt.Println("请检查:", user.Name)
}
```

## 3.2 测试的7个原则

### 原则1: 一个测试只测一件事

#### ❌ 坏示例
```go
func TestUserService(t *testing.T) {
    // 一个测试测太多东西
    service.Create(user)
    service.FindByID(user.ID)
    service.Update(user)
    service.Delete(user.ID)
    // 中间某步失败，后面都不执行
}
```

#### ✅ 好示例
```go
func TestUserService_Create(t *testing.T) {
    err := service.Create(user)
    if err != nil {
        t.Fatal(err)
    }
    if user.ID == 0 {
        t.Error("ID应该被设置")
    }
}

func TestUserService_FindByID(t *testing.T) {
    // 自己准备数据
    service.Create(user)

    // 只测查询
    found, err := service.FindByID(user.ID)
    if err != nil {
        t.Fatal(err)
    }
    // ...
}
```

### 原则2: 测试应该是文档

#### ❌ 坏示例
```go
func TestFunc1(t *testing.T) {
    result := DoSomething(10, 20, true)
    if result != 200 {
        t.Error("错了")
    }
}
// 😕 这是在测什么？为什么是200？
```

#### ✅ 好示例
```go
func TestCalculatePrice_WithTaxAndDiscount(t *testing.T) {
    // 给定: 商品100元，税率10%，折扣20%
    basePrice := 100.0
    taxRate := 0.10
    discountRate := 0.20

    // 当: 计算最终价格
    result := CalculatePrice(basePrice, taxRate, discountRate)

    // 那么: (100 * 0.8) * 1.1 = 88元
    expected := 88.0
    if result != expected {
        t.Errorf("got %.2f, want %.2f", result, expected)
    }
}
// ✅ 清晰明了，一看就懂
```

### 原则3: 使用表驱动测试

见 [2.4 表驱动测试模式](#24-表驱动测试模式)

### 原则4: 使用辅助函数减少重复

```go
// 辅助函数
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, _ := sql.Open("sqlite3", ":memory:")
    db.Exec(`CREATE TABLE users (...)`)
    t.Cleanup(func() { db.Close() })
    return db
}

func setupTestService(t *testing.T) *UserService {
    t.Helper()
    return NewUserService(setupTestDB(t))
}

// 测试变简洁
func TestUserService_Create(t *testing.T) {
    service := setupTestService(t)
    // 测试逻辑...
}
```

### 原则5: 清晰的错误信息

#### ❌ 坏示例
```go
if result != expected {
    t.Error("错了") // 什么错了？
}
```

#### ✅ 好示例
```go
if result != expected {
    t.Errorf("Calculate(%d, %d) = %d; want %d",
        input1, input2, result, expected)
}
// 输出: Calculate(10, 20) = 30; want 200
```

### 原则6: 测试边界条件

常见边界：
- **空值**: nil, "", [], 0
- **最大/最小值**
- **边界值**: -1, 0, 1
- **特殊值**: NaN, Infinity
- **错误情况**

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name string
        a, b float64
        want float64
        expectErr bool
    }{
        {"正常", 10, 2, 5, false},
        {"除以1", 10, 1, 10, false},
        {"负数", 10, -2, -5, false},
        {"被除数为0", 0, 5, 0, false},
        {"除以0", 10, 0, 0, true}, // 边界
        {"都是0", 0, 0, 0, true},  // 边界
        {"大数", 1e100, 1e50, 1e50, false}, // 边界
    }
    // ...
}
```

### 原则7: 使用Mock隔离依赖

```go
// 定义接口
type EmailSender interface {
    Send(to, subject, body string) error
}

// Mock实现
type MockEmailSender struct {
    SendFunc func(to, subject, body string) error
    Calls    []EmailCall
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    m.Calls = append(m.Calls, EmailCall{to, subject, body})
    if m.SendFunc != nil {
        return m.SendFunc(to, subject, body)
    }
    return nil
}

// 测试
func TestSendWelcomeEmail(t *testing.T) {
    mockEmail := &MockEmailSender{}
    service := NewUserService(mockEmail)

    service.SendWelcomeEmail(&User{
        Name:  "张三",
        Email: "test@example.com",
    })

    // 验证调用
    if len(mockEmail.Calls) != 1 {
        t.Fatalf("期望调用1次，实际%d次", len(mockEmail.Calls))
    }

    call := mockEmail.Calls[0]
    if call.To != "test@example.com" {
        t.Error("收件人错误")
    }
}
```

## 3.3 实战技巧集

### 1. 使用testdata目录

```
project/
├── parser.go
├── parser_test.go
└── testdata/
    ├── valid_config.json
    └── invalid_config.json
```

```go
func TestParseConfig(t *testing.T) {
    data, _ := os.ReadFile("testdata/valid_config.json")
    config, err := ParseConfig(data)
    // ...
}
```

### 2. 并行测试加速

```go
func TestFastOperation(t *testing.T) {
    t.Parallel() // 并行运行

    tests := []struct{
        // ...
    }{...}

    for _, tt := range tests {
        tt := tt // 捕获变量
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // 子测试也并行
            // ...
        })
    }
}
```

### 3. Golden Files模式

```go
func TestRender(t *testing.T) {
    result := RenderTemplate(data)

    goldenFile := "testdata/render.golden"

    if *update {
        // go test -update 更新golden file
        os.WriteFile(goldenFile, []byte(result), 0644)
    }

    expected, _ := os.ReadFile(goldenFile)
    if result != string(expected) {
        t.Errorf("输出不匹配")
    }
}
```

### 4. 测试Panic

```go
func TestDivideByZero_Panics(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Error("期望panic但没有发生")
        }
    }()

    DivideOrPanic(10, 0)
}
```

### 5. 测试超时

```go
func TestSlowOperation(t *testing.T) {
    timeout := time.After(5 * time.Second)
    done := make(chan bool)

    go func() {
        SlowOperation()
        done <- true
    }()

    select {
    case <-timeout:
        t.Fatal("操作超时")
    case <-done:
        // 成功
    }
}
```

## 3.4 测试覆盖率

### 查看覆盖率

```bash
# 显示百分比
go test -cover

# 详细报告
go test -coverprofile=coverage.out

# HTML可视化
go tool cover -html=coverage.out

# 按函数显示
go tool cover -func=coverage.out
```

### 覆盖率目标

- **70-80%**: 基本覆盖，良好起点
- **80-90%**: 较好覆盖，推荐目标
- **90%+**: 优秀覆盖，关键系统应达到

**注意**: 不要盲目追求100%，应关注：
- 关键业务逻辑
- 复杂算法
- 错误处理
- 边界条件

---

# Part 4: 深度理解

## 4.1 为什么基础库必须有大量测试

### 真实案例：基础库的bug

假设你的`base-tools/log`包有个bug：

```go
func (l *ZapLogger) Error(ctx context.Context, msg string, err error) {
    l.zap.Errorw(msg) // 忘记记录err了！
}
```

**影响范围**：
```
项目A: 用户服务 (10个开发者)
├─ 找不到登录失败原因
├─ 找不到订单创建失败原因
└─ 找不到支付失败原因

项目B: 订单服务 (8个开发者)
├─ 找不到库存扣减失败原因
└─ 找不到退款失败原因

项目C: 支付服务 (5个开发者)
├─ 找不到支付回调失败原因
└─ 找不到微信支付失败原因

总影响: 23个开发者 × N小时浪费
```

**如果有测试**：

```go
func TestZapLogger_Error_LogsErrorMessage(t *testing.T) {
    logger, buf := setupTestLogger(t)

    testErr := errors.New("测试错误")
    logger.Error(ctx, "发生错误", testErr)

    output := buf.String()
    if !strings.Contains(output, "测试错误") {
        t.Error("错误信息没有被记录！") // ✅ 发布前发现
    }
}
```

### ROI对比

| 项目类型 | 测试投入 | 影响范围 | Bug修复成本 | ROI |
|---------|---------|---------|------------|-----|
| **基础库** | 3天 | 10+项目 | 极高 | ⭐⭐⭐⭐⭐ |
| **业务代码** | 3天 | 单个服务 | 中等 | ⭐⭐⭐ |
| **一次性脚本** | 3天 | 只跑一次 | 低 | ⭐ |

### 成本计算

```
基础库测试投入：
- 3天写测试

基础库无测试成本（一次bug）：
- 10个项目 × 2人 × 0.5天排查 = 10人天
- 修复基础库 = 0.5人天
- 所有项目回归测试 = 5人天
- 紧急发版风险 = 无法估量

总成本：15.5人天 + 线上事故

ROI = (15.5 - 3) / 3 = 417% 🚀
```

### 行业数据

```
测试阶段发现bug: 1x
开发阶段发现bug: 10x
测试阶段发现bug: 100x
生产环境发现bug: 1000x

基础库生产bug = 1000x × 使用项目数
```

### 测试策略对比

**基础库(base-tools) - 高标准**
```
✓ 所有公开API: 100%覆盖
✓ 工具函数: 90%+覆盖
✓ 边界条件: 必须覆盖
✓ 并发安全: 必须测试
✓ 性能基准: 关键路径必须有

目标覆盖率: 80%+
```

**业务代码 - 适度**
```
✓ 核心业务: 80%+覆盖
✓ 一般功能: 60%+覆盖
✓ 简单CRUD: 可以不测试
✓ 快速迭代: 先上线后补测试

目标覆盖率: 60-70%
```

### 公式

```
基础库测试价值 = 使用项目数 × 单项目节省时间

例如:
10个项目用base-tools
每个项目避免1天排查
= 10 × 1天 = 10天节省

写测试花3天
净收益 = 10 - 3 = 7天

每次可能的bug都节省10天！
```

**结论：基础库的测试不是成本，是投资！**

## 4.2 为什么测试是重构的信心来源

### 常见困惑

```
你的想法:
写测试 → 花时间
重构代码 → 花时间
修改测试 → 又花时间
= 工作量翻倍！😫
```

### 关键误解

**重构不需要改测试！**

```
❌ 错误理解：改功能、改需求、改API
✅ 正确理解：不改功能，只改内部实现

重构 = 改变代码结构 + 不改变外部行为

外部行为没变 → 测试就不用改！
```

### 案例对比

#### 场景：优化log包性能

**重构前（慢）**：
```go
func (l *ZapLogger) Info(ctx context.Context, msg string, kv ...any) {
    // 每次创建新切片
    fields := []any{}
    if opID := getOperationID(ctx); opID != "" {
        fields = append(fields, "opID", opID)
    }
    // ...
    l.zap.Infow(msg, fields...)
}
```

**重构后（快）**：
```go
func (l *ZapLogger) Info(ctx context.Context, msg string, kv ...any) {
    // 预分配容量
    fields := make([]any, 0, len(kv)+6)
    if opID := getOperationID(ctx); opID != "" {
        fields = append(fields, "opID", opID)
    }
    // ...
    l.zap.Infow(msg, fields...)
}
```

#### 没有测试的噩梦

```
1. 改完代码，看起来没问题 ✓
2. 提交代码
3. 10个项目开始使用

第二天:
- "日志没有打印userID了！"
- "ctx为nil时panic了！"
- "operationID不见了！"

你: 😱😱😱
- 赶紧回滚
- 排查问题3天
- 团队不敢用新版本
- 老板质疑能力
```

#### 有测试的愉快

```go
// 重构前就有的测试（不需要改！）
func TestZapLogger_Info_WithContext(t *testing.T) {
    tests := []struct {
        name string
        setupContext func() context.Context
        expectedFields []string
    }{
        {
            name: "包含operationID",
            setupContext: func() context.Context {
                return mcontext.WithOperationID(ctx, "op-123")
            },
            expectedFields: []string{"opID", "op-123"},
        },
        {
            name: "nil context不应panic",
            setupContext: func() context.Context {
                return nil
            },
            expectedFields: []string{},
        },
    }
    // ...
}
```

```
1. 决定重构 ✓
2. 改代码
3. 运行: go test

结果:
✅ TestZapLogger_Info_WithContext ........ PASS
✅ TestZapLogger_Info_NilContext ......... PASS
✅ TestZapLogger_Info_WithCustomFields ... PASS

4. 看到全绿，信心满满 ✓
5. 提交发布
6. 一切正常 ✓

总耗时: 2小时
信心: 100% 😊
```

### 对比

**没有测试**：
```
改代码(2h) → 手动测试(4h) → 发布 → 出问题 →
回滚 → 修复(2h) → 手动测试(4h) → 再发布...

总耗时: 12+小时
信心: 0%
```

**有测试**：
```
改代码(2h) → 运行测试(2秒) ✅ → 发布

总耗时: 2小时
信心: 100%
```

### 为什么测试不用改？

**理解重构本质**：

```go
// 重构前
func CalculateTotal(items []Item) float64 {
    total := 0.0
    for i := 0; i < len(items); i++ {
        total += items[i].Price
    }
    return total
}

// 重构后（更优雅）
func CalculateTotal(items []Item) float64 {
    total := 0.0
    for _, item := range items {
        total += item.Price
    }
    return total
}

// 测试（不需要改！）
func TestCalculateTotal(t *testing.T) {
    items := []Item{{Price: 10}, {Price: 20}}
    result := CalculateTotal(items)
    expected := 30.0

    if result != expected {
        t.Errorf("期望%v, 得到%v", expected, result)
    }
}
// 输入没变：items
// 输出没变：总价
// 测试不用改！✅
```

### 测试是安全网

```
想象走钢丝:

没有测试 = 没有安全网
- 每步小心翼翼
- 害怕摔下去(出bug)
- 不敢尝试(不敢重构)
- 压力巨大

有测试 = 有安全网
- 可以大胆尝试
- 失败能及时发现
- 敢于重构
- 信心满满
```

### 时间对比

**有测试（1年）**：
```
初次: 写代码2天 + 写测试1天 = 3天
重构5次 × 2小时 = 10小时
Bug修复: 0次（测试保护）

总计: 3天 + 10小时
```

**无测试（1年）**：
```
初次: 写代码2天
重构1次失败回滚 = 2天
不敢重构，代码越来越乱 = 每次慢1小时
Bug修复3次 × 1天 = 3天
手动测试20次 × 30分钟 = 10小时

总计: 7天 + 10小时
```

**结论：有测试省4天！**

### 公式

```
没测试的重构 = 赌博
- 改完能用？不知道
- 会出bug？不知道
- 敢重构？不敢

有测试的重构 = 科学
- 改完能用？跑测试2秒知道
- 会出bug？测试告诉你
- 敢重构？随便重构，有保护
```

**测试不是负担，是翅膀！**

## 4.3 测试ROI分析

### 短期 vs 长期

**短期看**（第1个月）：
```
无测试: 快 ⚡
- 只写代码: 2天

有测试: 慢 🐢
- 写代码: 2天
- 写测试: 1天
总计: 3天

短期: 无测试赢
```

**长期看**（1年）：
```
无测试:
- 开发: 2天
- Bug修复: 5天
- 手动测试: 2天
- 不敢重构的技术债: 无法估量
总计: 9+天

有测试:
- 开发+测试: 3天
- Bug修复: 0.5天
- 自动测试: 0.01天
- 重构自由: 代码越来越好
总计: 3.5天

长期: 有测试赢 🏆
```

### 何时写什么测试

**单元测试 - 必须写** ⭐⭐⭐⭐⭐
```
✅ 所有公开函数
✅ 复杂业务逻辑
✅ 工具函数
✅ 错误处理
✅ 边界条件
```

**集成测试 - 应该写** ⭐⭐⭐⭐
```
✅ 数据库操作
✅ 外部API调用
✅ 文件系统操作
✅ 多组件协作
```

**基准测试 - 选择性写** ⭐⭐⭐
```
✅ 性能关键代码
✅ 算法优化验证
✅ 对比不同方案
⚠️ 不是所有代码都需要
```

**不需要测试** ❌
```
简单getter/setter
纯数据结构
简单包装函数
自动生成代码
```

---

# Part 5: 实践指南

## 5.1 针对base-tools的测试示例

### log包测试示例

#### simplify_test.go（单元测试）

```go
package log

import "testing"

func TestSlice_Format(t *testing.T) {
    tests := []struct {
        name     string
        input    Slice[int]
        expected int
    }{
        {
            name:     "短切片不截断",
            input:    Slice[int]{1, 2, 3, 4, 5},
            expected: 5,
        },
        {
            name:     "长切片截断到30",
            input:    make(Slice[int], 50),
            expected: slicePrintLen,
        },
        {
            name:     "刚好30元素",
            input:    make(Slice[int], slicePrintLen),
            expected: slicePrintLen,
        },
        {
            name:     "空切片",
            input:    Slice[int]{},
            expected: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := tt.input.Format()
            resultSlice, ok := result.([]int)
            if !ok {
                t.Fatal("Format()应返回切片")
            }

            if len(resultSlice) != tt.expected {
                t.Errorf("期望长度%d, 得到%d",
                    tt.expected, len(resultSlice))
            }
        })
    }
}
```

#### zap_integration_test.go（集成测试）

```go
package log

import (
    "context"
    "os"
    "path/filepath"
    "testing"
    "time"
)

func TestZapLogger_FileRotation(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过集成测试")
    }

    // 创建临时目录
    tmpDir, err := os.MkdirTemp("", "log_test")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tmpDir)

    // 初始化logger
    logger, err := NewZapLogger(
        "test", "testmodule", "", "",
        LevelDebug, false, false,
        tmpDir, 2, 1*time.Hour,
        "v1.0.0", false,
    )
    if err != nil {
        t.Fatalf("创建logger失败: %v", err)
    }
    defer logger.Flush()

    // 写日志
    ctx := context.Background()
    logger.Info(ctx, "测试消息1", "key", "value")
    logger.Debug(ctx, "测试消息2")
    logger.Warn(ctx, "测试警告", nil)

    // 验证文件创建
    files, err := filepath.Glob(filepath.Join(tmpDir, "test*"))
    if err != nil {
        t.Fatal(err)
    }

    if len(files) == 0 {
        t.Error("期望创建日志文件")
    }

    // 验证内容
    content, err := os.ReadFile(files[0])
    if err != nil {
        t.Fatal(err)
    }

    if len(content) == 0 {
        t.Error("日志文件不应为空")
    }
}
```

#### zap_bench_test.go（基准测试）

```go
package log

import (
    "context"
    "testing"
)

func BenchmarkZapLogger_Info(b *testing.B) {
    logger, err := NewZapLogger(
        "bench", "benchmodule", "", "",
        LevelInfo, false, false,
        "/tmp", 1, 24,
        "v1.0.0", false,
    )
    if err != nil {
        b.Fatal(err)
    }
    defer logger.Flush()

    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        logger.Info(ctx, "基准测试", "iteration", i)
    }
}

func BenchmarkLogLevels(b *testing.B) {
    levels := []struct {
        name  string
        level int
    }{
        {"Debug", LevelDebug},
        {"Info", LevelInfo},
        {"Warn", LevelWarn},
    }

    for _, l := range levels {
        b.Run(l.name, func(b *testing.B) {
            logger, _ := NewZapLogger(
                "bench", "benchmodule", "", "",
                l.level, false, false,
                "/tmp", 1, 24,
                "v1.0.0", false,
            )
            defer logger.Flush()

            ctx := context.Background()
            b.ResetTimer()

            for i := 0; i < b.N; i++ {
                logger.Info(ctx, "消息")
            }
        })
    }
}
```

### errs包测试示例

```go
// errs/errors_test.go

func TestWrap(t *testing.T) {
    originalErr := errors.New("原始错误")

    wrappedErr := Wrap(originalErr, "包装信息")

    if wrappedErr == nil {
        t.Fatal("Wrap不应返回nil")
    }

    errMsg := wrappedErr.Error()
    if !strings.Contains(errMsg, "包装信息") {
        t.Error("错误应包含包装信息")
    }
    if !strings.Contains(errMsg, "原始错误") {
        t.Error("错误应包含原始错误")
    }

    // 验证unwrap
    if errors.Unwrap(wrappedErr) != originalErr {
        t.Error("Unwrap应返回原始错误")
    }
}

func TestCodeError(t *testing.T) {
    err := NewCodeError(1001, "用户不存在")

    codeErr, ok := err.(CodeError)
    if !ok {
        t.Fatal("应实现CodeError接口")
    }

    if codeErr.Code() != 1001 {
        t.Errorf("错误码错误: got %d, want 1001", codeErr.Code())
    }

    if !strings.Contains(err.Error(), "用户不存在") {
        t.Error("错误信息不正确")
    }
}
```

## 5.2 测试最佳实践

### 文件组织

```
base-tools/
├── log/
│   ├── zap.go
│   ├── zap_test.go          # 单元测试
│   ├── zap_integration_test.go  # 集成测试
│   ├── zap_bench_test.go    # 基准测试
│   └── simplify_test.go
├── errs/
│   ├── errors.go
│   └── errors_test.go
└── testdata/               # 共享测试数据
    ├── test_config.json
    └── sample_input.txt
```

### 测试命名

```go
// ✅ 好的命名
func TestZapLogger_Error_LogsErrorMessage(t *testing.T)
func TestCalculateTotal_WithDiscount(t *testing.T)
func TestUserRepository_Create_ReturnsErrorOnDuplicate(t *testing.T)

// ❌ 不好的命名
func TestFunc1(t *testing.T)
func Test2(t *testing.T)
func TestIt(t *testing.T)
```

### 辅助函数

```go
// 标记为Helper
func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("期望无错误: %v", err)
    }
}

func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

// 使用
func TestSomething(t *testing.T) {
    result, err := DoSomething()
    assertNoError(t, err)
    assertEqual(t, result, "expected")
}
```

### Mock模式

```go
// 接口
type UserRepository interface {
    FindByID(id int) (*User, error)
    Create(user *User) error
}

// Mock
type MockUserRepository struct {
    FindByIDFunc func(id int) (*User, error)
    CreateFunc   func(user *User) error
    Calls        []string  // 记录调用
}

func (m *MockUserRepository) FindByID(id int) (*User, error) {
    m.Calls = append(m.Calls, fmt.Sprintf("FindByID(%d)", id))
    if m.FindByIDFunc != nil {
        return m.FindByIDFunc(id)
    }
    return nil, nil
}

// 测试
func TestUserService(t *testing.T) {
    mockRepo := &MockUserRepository{
        FindByIDFunc: func(id int) (*User, error) {
            return &User{ID: id, Name: "测试"}, nil
        },
    }

    service := NewUserService(mockRepo)
    user, err := service.GetUser(1)

    assertNoError(t, err)
    assertEqual(t, user.Name, "测试")

    // 验证调用
    if len(mockRepo.Calls) != 1 {
        t.Error("应调用1次")
    }
}
```

### 测试清理

```go
func TestWithCleanup(t *testing.T) {
    tmpFile, err := os.CreateTemp("", "test")
    if err != nil {
        t.Fatal(err)
    }

    // 自动清理
    t.Cleanup(func() {
        os.Remove(tmpFile.Name())
    })

    // 测试逻辑...
}
```

## 5.3 测试驱动开发(TDD)

### TDD循环

```
1. 红 🔴: 写测试 → 运行 → 失败
2. 绿 🟢: 写代码 → 运行 → 通过
3. 重构 ♻️: 优化代码 → 测试仍通过
```

### 示例

```go
// 1. 红：先写测试
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        price    float64
        discount float64
        want     float64
    }{
        {100, 0.1, 90},
        {200, 0.2, 160},
    }

    for _, tt := range tests {
        got := CalculateDiscount(tt.price, tt.discount)
        if got != tt.want {
            t.Errorf("got %v, want %v", got, tt.want)
        }
    }
}
// 运行: FAIL（函数还不存在）

// 2. 绿：写代码让测试通过
func CalculateDiscount(price, discount float64) float64 {
    return price * (1 - discount)
}
// 运行: PASS

// 3. 重构：优化代码
func CalculateDiscount(price, discount float64) float64 {
    if discount < 0 || discount > 1 {
        return price
    }
    return price * (1 - discount)
}
// 运行: PASS（测试保护重构）
```

## 5.4 从今天开始行动

### 3步走

#### 第1步：写第一个测试（今天）

```go
// log/simplify_test.go
func TestSlice_Format_ShortSlice(t *testing.T) {
    input := Slice[int]{1, 2, 3}
    result := input.Format()

    if len(result.([]int)) != 3 {
        t.Error("短切片不应截断")
    }
}
```

运行：
```bash
cd log
go test -run TestSlice_Format_ShortSlice -v
```

看到 `PASS` ✅，你就开始了！

#### 第2步：逐步覆盖（本周）

```
Day 1: log/simplify_test.go ✅
Day 2: log/zap_test.go（基本功能）
Day 3: errs/errors_test.go（核心功能）
Day 4: 添加边界测试
Day 5: 添加集成测试
```

#### 第3步：持续改进（持续）

```
Week 2: 提高覆盖率到60%
Week 3: 添加基准测试
Week 4: 覆盖率到80%
长期: 新功能必须带测试
```

### 检查清单

**写测试前** ✓
- [ ] 理解要测试的功能
- [ ] 确定测试类型
- [ ] 识别边界条件
- [ ] 设计测试用例

**写测试中** ✓
- [ ] 清晰的测试名称
- [ ] 一个测试一个点
- [ ] 使用表驱动
- [ ] 提取辅助函数
- [ ] 清晰的错误信息

**写测试后** ✓
- [ ] 测试通过
- [ ] 速度快(<100ms)
- [ ] 独立可重复
- [ ] 覆盖率合理(>70%)

**Code Review时** ✓
- [ ] 测试正确性
- [ ] 测试不脆弱
- [ ] 测试清晰
- [ ] 测试有价值

### 资源

**学习**：
- Go标准库测试
- Kubernetes测试
- Docker测试

**工具**：
```bash
# 测试覆盖率
go test -cover

# 基准测试
go test -bench=.

# 竞态检测
go test -race

# 测试某个包
go test ./log

# 详细输出
go test -v
```

---

## 总结

### 核心要点

1. **为什么测试**
   - 保证质量
   - 重构信心
   - 团队协作
   - 开发效率

2. **测什么**
   - 单元测试：必须写（公开API、业务逻辑、边界）
   - 集成测试：应该写（数据库、API、协作）
   - 基准测试：选择写（性能关键）

3. **怎么测**
   - 表驱动测试
   - AAA模式
   - Mock隔离
   - 辅助函数

4. **写好测试**
   - F.I.R.S.T原则
   - 7个原则
   - 实战技巧
   - 避免坏味道

5. **基础库特殊性**
   - 影响范围大（10×效应）
   - ROI极高（500%+）
   - 必须高覆盖率（80%+）

6. **测试与重构**
   - 重构不改测试
   - 测试是安全网
   - 2秒获得信心
   - 长期省时间

### 最后的话

```
测试不是负担，是投资！
测试不是成本，是收益！
测试不是约束，是自由！

好的测试 = 代码质量的保障
好的测试 = 重构的信心来源
好的测试 = 团队协作的基础

基础库必须有测试
业务代码应该有测试
关键路径一定有测试
```

**从今天开始，为你的base-tools写第一个测试！** 🚀

---

**记住**：

> 花1天写测试，省1个月debug
>
> 写3天测试，节省15天成本
>
> 测试是翅膀，让你的代码自由翱翔！

💡 **现在就开始**: `cd log && touch simplify_test.go`
