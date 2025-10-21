# 脚本说明

## test-changed-packages.sh

### 功能
智能检测变更的包并只运行这些包的测试，适用于扁平化组织的 Go 组件库。

### 使用场景
- **CI/CD 环境**: 自动检测 PR 或 Push 中变更的包，只测试相关代码
- **本地开发**: 手动指定对比分支，快速测试变更的包

### 工作原理
1. 检测变更的文件（通过 git diff）
2. 从文件路径提取顶级包名（如 `config/config.go` → `config`）
3. 对提取的包运行测试（如 `go test ./config/... ./errs/...`）
4. 生成覆盖率报告

### 使用方法

#### 在 CI 中使用（自动）
```bash
make test.changed
```

CI 环境会自动设置以下环境变量：
- `GITHUB_BASE_REF`: PR 的目标分支
- `GITHUB_REF`: 当前分支引用
- `GITHUB_EVENT_NAME`: 事件类型（push/pull_request）

#### 本地使用
```bash
# 对比 main 分支（默认）
./scripts/test-changed-packages.sh

# 对比指定分支
./scripts/test-changed-packages.sh origin/develop

# 通过 make 命令
make test.changed
```

### 示例

#### 场景 1: 修改了 config 包
```bash
$ git diff --name-only origin/main...HEAD
config/config.go
config/config_test.go

$ make test.changed
=== 检测变更的包 ===
变更的文件:
  config/config.go
  config/config_test.go

需要测试的包:
  config

=== 运行测试 ===
测试命令: go test ./config/... -race -coverprofile=coverage.out -covermode=atomic
```

#### 场景 2: 修改了多个包
```bash
$ git diff --name-only origin/main...HEAD
config/config.go
errs/coderr.go
utils/encrypt/encryption.go

$ make test.changed
=== 检测变更的包 ===
变更的文件:
  config/config.go
  errs/coderr.go
  utils/encrypt/encryption.go

需要测试的包:
  config
  errs
  utils

=== 运行测试 ===
测试命令: go test ./config/... ./errs/... ./utils/... -race -coverprofile=coverage.out -covermode=atomic
```

#### 场景 3: 只修改了文档或配置文件
```bash
$ git diff --name-only origin/main...HEAD
README.md
.github/workflows/coreCi.yml

$ make test.changed
=== 检测变更的包 ===
变更的文件:
  README.md
  .github/workflows/coreCi.yml

变更的文件不包含 Go 包，跳过测试
```

### 优势
1. **提升 CI 效率**: 只测试变更的包，大幅减少测试时间
2. **节省资源**: 避免重复测试未变更的稳定包
3. **快速反馈**: 开发者能更快得到测试结果
4. **智能检测**: 自动识别 PR 和 Push 场景

### 注意事项
- 脚本依赖 git 历史，确保 `fetch-depth: 0` 或已 fetch 完整历史
- 对于跨包依赖的变更，建议定期运行全量测试（如发布前）
- 如果检测不到变更，会自动回退到全量测试

### 与现有命令对比
| 命令 | 测试范围 | 使用场景 |
|------|---------|---------|
| `make test` | 所有包 | 本地全量测试、发布前验证 |
| `make test.changed` | 仅变更的包 | CI 增量测试、快速验证 |
| `go test ./config/...` | 指定包 | 针对性测试单个包 |
