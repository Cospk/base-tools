# utils 包

`utils` 包提供了一套全面的实用工具函数集合，涵盖数据处理、字符串操作、时间格式化、加密解密、网络工具等多个领域，是构建 Go 应用程序的基础工具库。

## 包结构

```
utils/
├── constants/      # 常量定义
├── datautil/       # 数据处理工具（切片、映射、集合操作）
├── encoding/       # 编码解码工具
├── encrypt/        # 加密解密工具
├── formatutil/     # 格式化工具
├── httputil/       # HTTP 工具
├── idutil/         # ID 生成工具
├── jsonutil/       # JSON 处理工具
├── mageutil/       # 图像处理工具
├── network/        # 网络工具
├── runtimeenv/     # 运行时环境工具
├── splitter/       # 分割器工具
├── stringutil/     # 字符串处理工具
└── timeutil/       # 时间处理工具
```

## 特性

- ✅ **泛型支持**: 充分利用 Go 1.18+ 泛型特性，提供类型安全的操作
- ✅ **高性能**: 优化的算法实现，最小化内存分配
- ✅ **零依赖**: 大部分工具仅依赖标准库
- ✅ **易用性**: 简洁直观的 API 设计
- ✅ **完整测试**: 每个包都有对应的单元测试
- ✅ **线程安全**: 所有工具函数都是线程安全的

## 安装

```bash
go get github.com/Cospk/base-tools/utils
```

## 子包详细说明

### 1. datautil - 数据处理工具

提供切片、映射等数据结构的高级操作。

#### 主要功能

- **切片操作**: 差集、交集、并集、去重、删除、查找
- **映射转换**: 切片与映射的相互转换
- **批量处理**: 批量操作和转换
- **泛型支持**: 支持任意类型的数据处理

#### 示例

```go
import "github.com/Cospk/base-tools/utils/datautil"

// 切片差集
a := []int{1, 2, 3, 4, 5}
b := []int{3, 4, 5, 6, 7}
diff := datautil.SliceSub(a, b) // [1, 2]

// 切片交集
intersection := datautil.SliceIntersect(a, b) // [3, 4, 5]

// 去重
data := []string{"a", "b", "a", "c", "b"}
unique := datautil.Distinct(data) // ["a", "b", "c"]

// 切片转映射
users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
userMap := datautil.SliceToMapAny(users, func(u User) (int, User) {
    return u.ID, u
})

// 批量操作
result := datautil.Slice(users, func(u User) string {
    return u.Name
}) // ["Alice", "Bob"]
```

### 2. stringutil - 字符串处理工具

提供字符串的各种操作和转换功能。

#### 主要功能

- **大小写转换**: CamelCase、snake_case、kebab-case 等格式转换
- **字符串操作**: 截取、填充、反转、分割
- **验证检查**: 判断字符串类型（数字、字母、中文等）
- **编码转换**: Base64、URL 编码等

#### 示例

```go
import "github.com/Cospk/base-tools/utils/stringutil"

// 大小写转换
camel := stringutil.ToCamelCase("hello_world")     // "helloWorld"
snake := stringutil.ToSnakeCase("HelloWorld")      // "hello_world"
kebab := stringutil.ToKebabCase("HelloWorld")      // "hello-world"

// 字符串截取
truncated := stringutil.Truncate("Hello World", 5) // "Hello..."

// 验证检查
isEmail := stringutil.IsEmail("user@example.com")  // true
isPhone := stringutil.IsPhone("13812345678")       // true

// 字符串生成
random := stringutil.RandomString(10)              // 随机字符串
```

### 3. timeutil - 时间处理工具

提供时间格式化、解析和计算功能。

#### 主要功能

- **格式化**: 支持多种时间格式的转换
- **解析**: 智能解析各种时间字符串
- **计算**: 时间差计算、时间段判断
- **时区处理**: 时区转换和处理

#### 示例

```go
import "github.com/Cospk/base-tools/utils/timeutil"

// 格式化时间
now := time.Now()
formatted := timeutil.Format(now, "2006-01-02 15:04:05")

// 解析时间字符串
parsed, _ := timeutil.Parse("2024-01-01 12:00:00")

// 获取时间戳
timestamp := timeutil.GetTimestamp()        // 秒级时间戳
millis := timeutil.GetMilliTimestamp()      // 毫秒级时间戳

// 时间计算
startOfDay := timeutil.GetStartOfDay(now)   // 当天开始时间
endOfDay := timeutil.GetEndOfDay(now)       // 当天结束时间
age := timeutil.GetAge(birthDate)           // 计算年龄

// 时间判断
isToday := timeutil.IsToday(someTime)
isWeekend := timeutil.IsWeekend(someTime)
```

### 4. jsonutil - JSON 处理工具

提供 JSON 的序列化、反序列化和操作功能。

#### 主要功能

- **序列化**: 对象转 JSON，支持美化输出
- **反序列化**: JSON 转对象，支持泛型
- **路径操作**: 通过路径获取/设置 JSON 值
- **验证**: JSON 格式验证

#### 示例

```go
import "github.com/Cospk/base-tools/utils/jsonutil"

// 序列化
user := User{Name: "Alice", Age: 30}
jsonStr := jsonutil.Marshal(user)
prettyJson := jsonutil.MarshalIndent(user)

// 反序列化
var newUser User
jsonutil.Unmarshal(jsonStr, &newUser)

// 路径操作
value := jsonutil.GetPath(jsonData, "user.address.city")
jsonutil.SetPath(jsonData, "user.age", 31)

// 合并 JSON
merged := jsonutil.Merge(json1, json2)
```

### 5. encrypt - 加密解密工具

提供常用的加密解密算法实现。

#### 主要功能

- **哈希算法**: MD5、SHA1、SHA256、SHA512
- **对称加密**: AES、DES、3DES
- **非对称加密**: RSA
- **编码**: Base64、Hex

#### 示例

```go
import "github.com/Cospk/base-tools/utils/encrypt"

// 哈希
md5Hash := encrypt.MD5("hello")
sha256Hash := encrypt.SHA256("hello")

// AES 加密解密
key := []byte("1234567890123456")
encrypted := encrypt.AESEncrypt([]byte("secret"), key)
decrypted := encrypt.AESDecrypt(encrypted, key)

// RSA 加密解密
publicKey, privateKey := encrypt.GenerateRSAKeyPair(2048)
ciphertext := encrypt.RSAEncrypt([]byte("message"), publicKey)
plaintext := encrypt.RSADecrypt(ciphertext, privateKey)
```

### 6. idutil - ID 生成工具

提供各种 ID 生成策略。

#### 主要功能

- **UUID**: 生成标准 UUID
- **雪花算法**: 分布式 ID 生成
- **短 ID**: 生成短唯一标识符
- **自定义 ID**: 按规则生成 ID

#### 示例

```go
import "github.com/Cospk/base-tools/utils/idutil"

// UUID
uuid := idutil.NewUUID()                    // "550e8400-e29b-41d4-a716-446655440000"

// 雪花 ID
snowflake := idutil.NewSnowflakeID()        // 1234567890123456789

// 短 ID
shortId := idutil.NewShortID()              // "Kb8J2xm"

// 自定义前缀 ID
orderId := idutil.NewPrefixID("ORD")        // "ORD20240101120000001"
```

### 7. httputil - HTTP 工具

提供 HTTP 请求和响应的辅助功能。

#### 主要功能

- **请求构建**: 简化 HTTP 请求的创建
- **响应处理**: 统一的响应处理
- **中间件**: 常用 HTTP 中间件
- **工具函数**: IP 获取、User-Agent 解析等

#### 示例

```go
import "github.com/Cospk/base-tools/utils/httputil"

// 发送 GET 请求
resp, err := httputil.Get("https://api.example.com/users")

// 发送 POST 请求
data := map[string]interface{}{"name": "Alice"}
resp, err := httputil.PostJSON("https://api.example.com/users", data)

// 获取客户端 IP
ip := httputil.GetClientIP(request)

// 解析 User-Agent
ua := httputil.ParseUserAgent(request.Header.Get("User-Agent"))
```

### 8. network - 网络工具

提供网络相关的工具函数。

#### 主要功能

- **IP 操作**: IP 验证、转换、计算
- **端口检测**: 端口可用性检查
- **网络信息**: 获取本机网络信息

#### 示例

```go
import "github.com/Cospk/base-tools/utils/network"

// IP 验证
isValid := network.IsValidIP("192.168.1.1")      // true
isIPv4 := network.IsIPv4("192.168.1.1")          // true
isIPv6 := network.IsIPv6("::1")                  // true

// 获取本机 IP
localIP := network.GetLocalIP()                   // "192.168.1.100"
publicIP := network.GetPublicIP()                 // "203.0.113.1"

// 端口检测
available := network.IsPortAvailable(8080)        // true
```

### 9. formatutil - 格式化工具

提供各种数据格式化功能。

#### 主要功能

- **数字格式化**: 千分位、百分比、货币等
- **文件大小**: 字节数转人类可读格式
- **时间格式化**: 相对时间、持续时间等

#### 示例

```go
import "github.com/Cospk/base-tools/utils/formatutil"

// 数字格式化
thousands := formatutil.ThousandsSeparator(1234567)    // "1,234,567"
percent := formatutil.ToPercent(0.1234, 2)             // "12.34%"
currency := formatutil.ToCurrency(1234.56, "￥")       // "￥1,234.56"

// 文件大小
size := formatutil.FormatFileSize(1024*1024*5.5)       // "5.5 MB"

// 相对时间
relative := formatutil.RelativeTime(time.Now().Add(-time.Hour)) // "1小时前"
```

### 10. mageutil - 图像处理工具

提供图像的基本处理功能。

#### 主要功能

- **缩放**: 按比例或指定尺寸缩放
- **裁剪**: 裁剪图像指定区域
- **水印**: 添加文字或图片水印
- **格式转换**: 支持常见图片格式转换

#### 示例

```go
import "github.com/Cospk/base-tools/utils/mageutil"

// 缩放图像
resized := mageutil.Resize(img, 800, 600)

// 裁剪图像
cropped := mageutil.Crop(img, 0, 0, 200, 200)

// 添加水印
watermarked := mageutil.AddWatermark(img, "© 2024", position)

// 格式转换
mageutil.ConvertFormat("input.png", "output.jpg", 90)
```

## 性能优化

所有工具函数都经过性能优化：

- 使用对象池减少内存分配
- 避免不必要的复制
- 使用并发处理大数据集
- 缓存计算结果

## 测试

每个包都有完整的单元测试：

```bash
# 运行所有测试
go test ./utils/...

# 运行特定包的测试
go test ./utils/datautil/

# 查看测试覆盖率
go test -cover ./utils/...

# 运行基准测试
go test -bench=. ./utils/...
```

## 最佳实践

1. **选择合适的工具**: 根据具体需求选择最合适的工具函数
2. **错误处理**: 始终检查返回的错误
3. **性能考虑**: 对于大数据集，考虑使用流式处理或分批处理
4. **并发安全**: 所有工具函数都是并发安全的，可以在 goroutine 中使用

## 贡献指南

欢迎贡献新的工具函数或改进现有实现：

1. Fork 项目
2. 创建功能分支
3. 编写代码和测试
4. 提交 Pull Request

## 许可证

本项目采用与 base-tools 相同的许可证。

## 更新日志

### v1.0.0 (2025-10-20)
- ✨ 初始版本发布
- ✅ 包含 14 个工具包
- ✅ 完整的泛型支持
- ✅ 全面的单元测试
- 📚 完整的文档和示例
