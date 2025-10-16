# GitHub Actions Workflows 语法参考手册

本文档提供 GitHub Actions Workflows 的常用语法和最佳实践，方便编写和维护工作流时快速查阅。

---

## 目录

- [1. 工作流触发器 (on)](#1-工作流触发器-on)
- [2. 环境变量 (env)](#2-环境变量-env)
- [3. 矩阵构建 (strategy)](#3-矩阵构建-strategy)
- [4. 条件执行 (if)](#4-条件执行-if)
- [5. 输出和依赖](#5-输出和依赖)
- [6. 常用上下文变量](#6-常用上下文变量)
- [7. Action 使用](#7-action-使用)
- [8. 权限控制 (permissions)](#8-权限控制-permissions)
- [9. 最佳实践建议](#9-最佳实践建议)
- [10. 相关资源](#10-相关资源)

---

## 1. 工作流触发器 (on)

### 单一事件触发

```yaml
on: push
```

### 多个事件触发

```yaml
on: [push, pull_request]
```

### 事件过滤 - 分支过滤

```yaml
on:
  push:
    branches:
      - main                         # 仅 main 分支
      - 'releases/**'                # 通配符匹配
      - '!dev'                       # 排除 dev 分支
```

### 事件过滤 - 标签过滤

```yaml
on:
  push:
    tags:
      - v1.*                         # 匹配 v1.x 标签
      - 'v*.*.*'                     # 匹配语义化版本
```

### 事件过滤 - 路径过滤

```yaml
on:
  push:
    paths:
      - '**.go'                      # 仅 Go 文件变更时触发
      - 'src/**'                     # src 目录下的任何文件
    paths-ignore:
      - '**.md'                      # 忽略 Markdown 文件
      - 'docs/**'                    # 忽略 docs 目录
```

### Pull Request 事件类型

```yaml
on:
  pull_request:
    branches: [main, dev]
    types:
      - opened                       # PR 创建时
      - synchronize                  # PR 更新时
      - reopened                     # PR 重新打开时
      - closed                       # PR 关闭时
```

### 定时触发

```yaml
on:
  schedule:
    - cron: '0 0 * * *'              # 每天午夜触发
    - cron: '0 */6 * * *'            # 每 6 小时触发一次
```

**Cron 表达式格式**：`分 时 日 月 周`
- `*`：任意值
- `*/n`：每 n 个单位
- `a,b,c`：指定多个值
- 示例：`'30 5,17 * * *'` 表示每天 5:30 和 17:30

### 手动触发

```yaml
on:
  workflow_dispatch:                 # 允许手动触发
    inputs:                          # 手动触发时的输入参数
      environment:
        description: 'Environment to deploy'
        required: true
        default: 'staging'
        type: choice
        options:
          - staging
          - production
```

### 其他触发事件

```yaml
on:
  issues:
    types: [opened, edited, deleted]

  release:
    types: [published, edited]

  workflow_run:                      # 当其他工作流完成时触发
    workflows: ["CI"]
    types: [completed]
```

---

## 2. 环境变量 (env)

### 全局环境变量

```yaml
env:
  GO_VERSION: 1.23.4
  NODE_ENV: production
```

### 任务级环境变量

```yaml
jobs:
  build:
    env:
      BUILD_ENV: staging
      DEBUG: true
```

### 步骤级环境变量

```yaml
steps:
  - name: Build
    env:
      CGO_ENABLED: 0
      GOOS: linux
    run: go build
```

### 使用环境变量

```yaml
steps:
  - run: echo "Go version is $GO_VERSION"
  - run: echo "Using ${{ env.GO_VERSION }}"
```

### 动态设置环境变量

```yaml
steps:
  - name: Set environment variable
    run: |
      echo "BUILD_TIME=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" >> $GITHUB_ENV
      echo "COMMIT_SHA=${GITHUB_SHA::7}" >> $GITHUB_ENV

  - name: Use environment variable
    run: |
      echo "Build time: $BUILD_TIME"
      echo "Commit: $COMMIT_SHA"
```

---

## 3. 矩阵构建 (strategy)

### 基本矩阵

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
    go: [1.21, 1.22, 1.23]
```

这会创建 3 × 3 = 9 个任务实例。

### 矩阵配置选项

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest]
    go: [1.22, 1.23]
  fail-fast: false                   # 一个失败不影响其他任务
  max-parallel: 2                    # 最大并行任务数
```

### 引用矩阵变量

```yaml
runs-on: ${{ matrix.os }}

steps:
  - uses: actions/setup-go@v5
    with:
      go-version: ${{ matrix.go }}
```

### 包含额外组合

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest]
    go: [1.22, 1.23]
    include:
      - os: windows-latest           # 添加额外的组合
        go: 1.23
        experimental: true
```

### 排除特定组合

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
    go: [1.21, 1.22, 1.23]
    exclude:
      - os: macos-latest             # 排除 macOS + Go 1.21
        go: 1.21
```

---

## 4. 条件执行 (if)

### 基于分支条件

```yaml
steps:
  - name: Deploy to production
    if: github.ref == 'refs/heads/main'
    run: make deploy
```

### 基于事件类型

```yaml
- name: Comment on PR
  if: github.event_name == 'pull_request'
  run: gh pr comment ${{ github.event.number }} --body "CI passed!"
```

### 基于步骤状态

```yaml
- name: Debug on failure
  if: failure()                      # 前序步骤失败时执行
  run: cat debug.log

- name: Cleanup
  if: always()                       # 总是执行
  run: make clean

- name: Success notification
  if: success()                      # 仅成功时执行
  run: echo "All tests passed!"
```

### 组合条件

```yaml
- name: Deploy to staging
  if: |
    github.event_name == 'push' &&
    github.ref == 'refs/heads/develop' &&
    success()
  run: make deploy-staging
```

### 基于输出值

```yaml
- id: check_files
  run: |
    if [[ -f "go.mod" ]]; then
      echo "has_gomod=true" >> $GITHUB_OUTPUT
    fi

- name: Run Go tests
  if: steps.check_files.outputs.has_gomod == 'true'
  run: go test ./...
```

### 常用条件表达式

```yaml
# 检查环境
if: runner.os == 'Linux'

# 检查 actor
if: github.actor == 'dependabot[bot]'

# 检查标签
if: startsWith(github.ref, 'refs/tags/')

# 检查提交消息
if: contains(github.event.head_commit.message, '[skip ci]') == false
```

---

## 5. 输出和依赖

### 步骤输出

```yaml
steps:
  - id: get_version
    run: |
      VERSION=$(cat VERSION)
      echo "version=$VERSION" >> $GITHUB_OUTPUT
      echo "short_sha=${GITHUB_SHA::7}" >> $GITHUB_OUTPUT

  - name: Use output
    run: |
      echo "Version: ${{ steps.get_version.outputs.version }}"
      echo "SHA: ${{ steps.get_version.outputs.short_sha }}"
```

### 任务输出

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      version: ${{ steps.get_version.outputs.version }}
      artifact_name: ${{ steps.build.outputs.name }}
    steps:
      - id: get_version
        run: echo "version=1.0.0" >> $GITHUB_OUTPUT

      - id: build
        run: echo "name=myapp-${{ steps.get_version.outputs.version }}" >> $GITHUB_OUTPUT
```

### 任务依赖

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: make build

  test:
    needs: build                     # 依赖 build 任务
    runs-on: ubuntu-latest
    steps:
      - run: make test

  deploy:
    needs: [build, test]             # 依赖多个任务
    runs-on: ubuntu-latest
    steps:
      - run: make deploy
```

### 使用前置任务的输出

```yaml
jobs:
  build:
    outputs:
      version: ${{ steps.version.outputs.value }}
    steps:
      - id: version
        run: echo "value=1.2.3" >> $GITHUB_OUTPUT

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - run: echo "Deploying version ${{ needs.build.outputs.version }}"
```

---

## 6. 常用上下文变量

### GitHub 上下文

```yaml
${{ github.event_name }}             # 事件名称 (push, pull_request 等)
${{ github.ref }}                    # 分支或标签引用 (refs/heads/main)
${{ github.ref_name }}               # 分支或标签名称 (main, v1.0.0)
${{ github.sha }}                    # 完整的 commit SHA
${{ github.actor }}                  # 触发工作流的用户
${{ github.repository }}             # 仓库名称 (owner/repo)
${{ github.repository_owner }}       # 仓库所有者
${{ github.workspace }}              # 工作目录路径
${{ github.run_id }}                 # 工作流运行 ID
${{ github.run_number }}             # 工作流运行编号
```

### 事件上下文

```yaml
# Pull Request 事件
${{ github.event.pull_request.number }}
${{ github.event.pull_request.title }}
${{ github.event.pull_request.head.ref }}

# Issue 事件
${{ github.event.issue.number }}
${{ github.event.issue.title }}

# Push 事件
${{ github.event.head_commit.message }}
${{ github.event.head_commit.author.name }}
```

### Runner 上下文

```yaml
${{ runner.os }}                     # 运行器操作系统 (Linux, macOS, Windows)
${{ runner.arch }}                   # 架构 (X64, ARM, ARM64)
${{ runner.name }}                   # Runner 名称
${{ runner.temp }}                   # 临时目录路径
```

### Job 上下文

```yaml
${{ job.status }}                    # 任务状态
${{ job.container.id }}              # 容器 ID
```

### Steps 上下文

```yaml
${{ steps.step_id.outputs.name }}    # 步骤输出
${{ steps.step_id.conclusion }}      # 步骤结论 (success, failure, etc.)
```

### Secrets 上下文

```yaml
${{ secrets.GITHUB_TOKEN }}          # GitHub 自动提供的 token
${{ secrets.MY_SECRET }}             # 自定义密钥
```

### 环境变量

```yaml
${{ env.GO_VERSION }}                # 环境变量
```

---

## 7. Action 使用

### 使用特定版本

```yaml
- uses: actions/checkout@v4          # 推荐：使用主版本号
- uses: actions/checkout@v4.1.0      # 使用具体版本
```

### 使用分支

```yaml
- uses: owner/repo@main              # 使用 main 分支
- uses: owner/repo@develop           # 使用 develop 分支
```

### 使用提交 SHA（最安全）

```yaml
- uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada57f0ab
```

### 传递参数

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: 1.23.4
    cache: true
    cache-dependency-path: go.sum
```

### 引用本地 Action

```yaml
- uses: ./.github/actions/my-action  # 本地 action
  with:
    param1: value1
```

### 常用官方 Actions

```yaml
# 检出代码
- uses: actions/checkout@v4
  with:
    fetch-depth: 0                   # 获取完整历史
    submodules: true                 # 包含子模块

# 设置 Go 环境
- uses: actions/setup-go@v5
  with:
    go-version: '1.23.4'

# 设置 Node.js 环境
- uses: actions/setup-node@v4
  with:
    node-version: '20'

# 缓存依赖
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-

# 上传构建产物
- uses: actions/upload-artifact@v3
  with:
    name: my-artifact
    path: dist/

# 下载构建产物
- uses: actions/download-artifact@v3
  with:
    name: my-artifact
    path: dist/
```

---

## 8. 权限控制 (permissions)

### 详细权限设置

```yaml
permissions:
  contents: read                     # 读取仓库内容
  contents: write                    # 写入仓库内容（提交、创建 Release）
  pull-requests: read                # 读取 PR
  pull-requests: write               # 写入 PR（评论、标签）
  issues: read                       # 读取 Issue
  issues: write                      # 写入 Issue（评论、标签、关闭）
  packages: write                    # 发布包
  deployments: write                 # 创建部署
  id-token: write                    # OIDC token（用于云服务认证）
  checks: write                      # 创建检查
  statuses: write                    # 创建状态
```

### 工作流级别权限

```yaml
# 所有任务的默认权限
permissions:
  contents: read
  pull-requests: write
```

### 任务级别权限

```yaml
jobs:
  build:
    permissions:
      contents: read                 # 覆盖工作流级别权限

  deploy:
    permissions:
      contents: write
      deployments: write
```

### 批量设置

```yaml
permissions: read-all                # 所有权限只读
```

```yaml
permissions: write-all               # 所有权限可写（不推荐）
```

### 禁用所有权限

```yaml
permissions: {}
```

### 最小权限原则示例

```yaml
name: CI

permissions:
  contents: read                     # 默认只读

jobs:
  test:
    # 继承默认只读权限
    steps:
      - uses: actions/checkout@v4
      - run: make test

  release:
    permissions:
      contents: write                # 仅发布任务需要写权限
    steps:
      - uses: actions/checkout@v4
      - run: make release
```

---

## 9. 最佳实践建议

### 1. 版本固定

**推荐做法**：
```yaml
- uses: actions/checkout@v4          # ✅ 使用主版本号
- uses: actions/setup-go@v5.0.0      # ✅ 使用具体版本
```

**不推荐做法**：
```yaml
- uses: actions/checkout@latest      # ❌ 不稳定
- uses: actions/checkout@main        # ❌ 可能突然变更
```

**原因**：提高稳定性和可重复性，避免上游变更导致的意外失败。

---

### 2. 缓存利用

#### Go 项目缓存

```yaml
- uses: actions/cache@v3
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

#### Node.js 项目缓存

```yaml
- uses: actions/cache@v3
  with:
    path: ~/.npm
    key: ${{ runner.os }}-node-${{ hashFiles('**/package-lock.json') }}
    restore-keys: |
      ${{ runner.os }}-node-
```

#### Docker 层缓存

```yaml
- uses: docker/build-push-action@v5
  with:
    context: .
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

---

### 3. 失败处理

#### 继续执行

```yaml
- name: Lint code
  run: make lint
  continue-on-error: true            # 失败不中断后续步骤
```

#### 失败后清理

```yaml
- name: Cleanup on failure
  if: failure()
  run: |
    docker-compose down
    rm -rf temp/
```

#### 总是执行清理

```yaml
- name: Cleanup
  if: always()
  run: make clean
```

#### 超时控制

```yaml
- name: Long running test
  run: make integration-test
  timeout-minutes: 30                # 30 分钟超时
```

---

### 4. 安全性

#### 使用 Secrets

```yaml
# ❌ 不要硬编码敏感信息
- run: echo "password123" | docker login

# ✅ 使用 secrets
- run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login
```

#### 最小化权限

```yaml
# ✅ 默认只读
permissions:
  contents: read

jobs:
  deploy:
    # 仅需要时授予写权限
    permissions:
      contents: write
```

#### 防止脚本注入

```yaml
# ❌ 危险：用户输入直接用于命令
- run: echo "${{ github.event.issue.title }}"

# ✅ 安全：使用环境变量
- env:
    ISSUE_TITLE: ${{ github.event.issue.title }}
  run: echo "$ISSUE_TITLE"
```

#### 使用 paths-ignore

```yaml
on:
  push:
    paths-ignore:
      - '**.md'
      - 'docs/**'
```

---

### 5. 性能优化

#### 并发控制

```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true           # 取消进行中的旧工作流
```

#### 合理使用矩阵

```yaml
strategy:
  matrix:
    os: [ubuntu-latest]              # 仅在必要时测试多平台
    go: [1.23]                       # 仅测试关键版本
```

#### 减少 checkout 时间

```yaml
- uses: actions/checkout@v4
  with:
    fetch-depth: 1                   # 浅克隆，仅最新提交
```

#### 条件执行

```yaml
- name: Deploy docs
  if: github.ref == 'refs/heads/main' && contains(github.event.head_commit.message, '[docs]')
  run: make deploy-docs
```

---

### 6. 可维护性

#### 使用可复用工作流

```yaml
# .github/workflows/reusable-test.yml
on:
  workflow_call:
    inputs:
      go-version:
        required: true
        type: string

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ inputs.go-version }}
      - run: go test ./...
```

调用：
```yaml
jobs:
  test:
    uses: ./.github/workflows/reusable-test.yml
    with:
      go-version: '1.23.4'
```

#### 环境隔离

```yaml
jobs:
  deploy-staging:
    environment: staging
    steps:
      - run: deploy to staging

  deploy-prod:
    environment: production          # 可设置审批要求
    needs: deploy-staging
    steps:
      - run: deploy to production
```

---

## 10. 安全配置和注意事项

### 1. 敏感信息管理

#### 绝对不要硬编码敏感信息

```yaml
# ❌ 错误示例 - 永远不要这样做
env:
  DATABASE_PASSWORD: "mypassword123"
  API_KEY: "sk_live_abc123xyz"
  AWS_SECRET_KEY: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

# ❌ 错误示例 - 不要在脚本中硬编码
- run: |
    export TOKEN="ghp_abc123xyz"
    curl -H "Authorization: token $TOKEN" api.github.com
```

#### ✅ 正确做法：使用 GitHub Secrets

**在仓库中配置 Secrets**：
1. 进入仓库 Settings → Secrets and variables → Actions
2. 点击 "New repository secret"
3. 添加密钥名称和值

```yaml
# ✅ 正确示例 - 使用 secrets
env:
  DATABASE_PASSWORD: ${{ secrets.DATABASE_PASSWORD }}
  API_KEY: ${{ secrets.API_KEY }}
  AWS_SECRET_KEY: ${{ secrets.AWS_SECRET_KEY }}

# ✅ 使用 secrets 进行身份验证
- name: Login to Docker Hub
  run: |
    echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin

# ✅ 使用 secrets 调用 API
- name: Call API
  run: |
    curl -H "Authorization: token ${{ secrets.GITHUB_TOKEN }}" \
         https://api.github.com/repos/${{ github.repository }}
```

---

### 2. 防止代码注入攻击

#### 用户输入注入风险

```yaml
# ❌ 危险：直接使用用户输入
- name: Echo issue title
  run: echo "Issue title: ${{ github.event.issue.title }}"
  # 如果 issue 标题是: "; rm -rf / #" 将会执行删除命令！

# ❌ 危险：PR 标题或评论可被利用
- run: |
    TITLE="${{ github.event.pull_request.title }}"
    echo $TITLE > title.txt

# ✅ 安全：使用环境变量传递
- name: Echo issue title safely
  env:
    ISSUE_TITLE: ${{ github.event.issue.title }}
  run: echo "Issue title: $ISSUE_TITLE"

# ✅ 安全：使用输入文件
- name: Process PR title
  run: |
    cat << 'EOF' > title.txt
    ${{ github.event.pull_request.title }}
    EOF
    cat title.txt
```

#### 防止脚本注入的最佳实践

```yaml
# ✅ 使用环境变量而不是直接插值
- name: Handle user input
  env:
    USER_INPUT: ${{ github.event.comment.body }}
    PR_TITLE: ${{ github.event.pull_request.title }}
    COMMIT_MSG: ${{ github.event.head_commit.message }}
  run: |
    echo "Comment: $USER_INPUT"
    echo "Title: $PR_TITLE"
    echo "Message: $COMMIT_MSG"
```

---

### 3. 权限最小化原则

#### 默认使用最小权限

```yaml
# ✅ 推荐：工作流级别设置最小权限
name: CI

permissions:
  contents: read                     # 默认只读

jobs:
  test:
    runs-on: ubuntu-latest
    # 继承只读权限，足够运行测试
    steps:
      - uses: actions/checkout@v4
      - run: make test

  release:
    runs-on: ubuntu-latest
    # 仅在需要时提升权限
    permissions:
      contents: write                # 创建 release 需要写权限
      packages: write                # 发布包需要写权限
    steps:
      - uses: actions/checkout@v4
      - run: make release
```

#### 避免过度授权

```yaml
# ❌ 不推荐：授予所有权限
permissions: write-all

# ❌ 不推荐：授予不必要的权限
permissions:
  contents: write
  pull-requests: write
  issues: write
  packages: write
  deployments: write
  # ... 实际只需要 contents: read

# ✅ 推荐：按需授权
permissions:
  contents: read                     # 只给需要的权限
```

---

### 4. 第三方 Actions 安全

#### 固定 Action 版本

```yaml
# ❌ 不安全：使用分支
- uses: some-org/some-action@main    # main 分支可能被恶意修改

# ⚠️  一般：使用标签
- uses: some-org/some-action@v1      # 标签可以被移动

# ✅ 最安全：使用完整 SHA
- uses: some-org/some-action@8e5e7e5ab8b370d6c329ec480221332ada57f0ab
  # SHA 是不可变的，最安全
```

#### 审查第三方 Actions

使用第三方 Actions 前应该：
1. 检查 Action 的源代码仓库
2. 查看 Action 的权限要求
3. 检查 Action 的维护状态和社区评价
4. 优先使用官方或知名组织的 Actions
5. 定期审查和更新使用的 Actions

```yaml
# ✅ 可信的官方 Actions
- uses: actions/checkout@v4
- uses: actions/setup-go@v5
- uses: docker/build-push-action@v5

# ⚠️  使用不熟悉的 Actions 时要谨慎
- uses: unknown-user/suspicious-action@v1
  # 先审查源代码！
```

---

### 5. Secrets 管理最佳实践

#### Secrets 的生命周期管理

```yaml
# ✅ 为不同环境使用不同的 secrets
jobs:
  deploy-staging:
    environment: staging
    steps:
      - run: |
          echo "Deploying to staging"
          # 使用 staging 环境的 secrets
          echo "${{ secrets.STAGING_API_KEY }}"

  deploy-production:
    environment: production
    needs: deploy-staging
    steps:
      - run: |
          echo "Deploying to production"
          # 使用 production 环境的 secrets
          echo "${{ secrets.PROD_API_KEY }}"
```

#### Secrets 使用注意事项

```yaml
# ❌ 不要在日志中打印 secrets
- run: echo "My secret is ${{ secrets.MY_SECRET }}"
  # GitHub 会自动屏蔽，但仍不安全

# ❌ 不要将 secrets 写入文件后提交
- run: |
    echo "${{ secrets.API_KEY }}" > config.txt
    git add config.txt
    git commit -m "Add config"

# ✅ 仅在内存中使用 secrets
- run: |
    curl -H "Authorization: Bearer ${{ secrets.API_TOKEN }}" \
         https://api.example.com/data

# ✅ 使用临时文件并清理
- run: |
    echo "${{ secrets.CERTIFICATE }}" > /tmp/cert.pem
    use-cert /tmp/cert.pem
    rm -f /tmp/cert.pem
```

#### 定期轮换 Secrets

建议：
- 定期更新 API 密钥和令牌
- 当团队成员离职时轮换相关 secrets
- 使用有时效性的令牌（如临时 AWS credentials）
- 监控 secrets 的使用情况

---

### 6. 防止敏感文件泄露

#### 使用 .gitignore

确保敏感文件不会被提交：

```gitignore
# .gitignore
.env
.env.local
.env.*.local
*.pem
*.key
*.p12
*.pfx
credentials.json
secrets.yml
config/secrets.yml
terraform.tfstate
.aws/credentials
```

#### Workflow 中避免提交敏感信息

```yaml
# ✅ 生成配置文件但不提交
- name: Create config
  run: |
    cat << EOF > config.yml
    api_key: ${{ secrets.API_KEY }}
    database:
      password: ${{ secrets.DB_PASSWORD }}
    EOF
    # 使用配置但不提交
    ./app --config config.yml

# ✅ 如果必须提交，使用加密
- name: Encrypt and commit
  run: |
    echo "${{ secrets.GPG_KEY }}" | gpg --import
    echo "${{ secrets.API_KEY }}" | gpg --encrypt --recipient user@example.com > api_key.gpg
    git add api_key.gpg
    git commit -m "Add encrypted API key"
```

---

### 7. 网络安全

#### 限制出站连接

```yaml
# ✅ 仅连接到已知的可信主机
- name: Download dependencies
  run: |
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    # 验证签名

# ⚠️  避免下载未知来源的脚本并直接执行
# ❌ 危险
- run: curl https://unknown-site.com/install.sh | bash

# ✅ 更安全的做法
- run: |
    curl -o install.sh https://known-site.com/install.sh
    # 验证签名或校验和
    sha256sum -c install.sh.sha256
    bash install.sh
```

#### 使用 HTTPS

```yaml
# ✅ 始终使用 HTTPS
- run: curl https://api.example.com/data

# ❌ 避免使用 HTTP
- run: curl http://api.example.com/data
```

---

### 8. 依赖安全

#### 依赖锁定和审计

```yaml
# Go 项目
- name: Verify dependencies
  run: |
    go mod verify                    # 验证依赖完整性
    go mod tidy                      # 清理未使用的依赖

- name: Security scan
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...                # 检查已知漏洞

# Node.js 项目
- name: Audit dependencies
  run: |
    npm audit                        # 检查安全漏洞
    npm audit fix                    # 自动修复

# 使用 Dependabot 自动更新依赖
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
```

---

### 9. 容器安全

#### 使用可信的基础镜像

```yaml
# ✅ 使用官方镜像和固定标签
- name: Build Docker image
  run: |
    docker build \
      --build-arg BASE_IMAGE=golang:1.23.4-alpine \
      -t myapp:${{ github.sha }} .

# ❌ 避免使用 latest 标签
- run: docker build --build-arg BASE_IMAGE=golang:latest .
```

#### 扫描镜像漏洞

```yaml
- name: Build image
  run: docker build -t myapp:${{ github.sha }} .

- name: Scan for vulnerabilities
  uses: aquasecurity/trivy-action@master
  with:
    image-ref: myapp:${{ github.sha }}
    format: 'sarif'
    output: 'trivy-results.sarif'

- name: Upload scan results
  uses: github/codeql-action/upload-sarif@v2
  with:
    sarif_file: 'trivy-results.sarif'
```

---

### 10. 日志和审计

#### 避免在日志中泄露敏感信息

```yaml
# ❌ 不要记录敏感信息
- run: |
    echo "API Key: ${{ secrets.API_KEY }}"
    echo "Database password: $DB_PASSWORD"

# ✅ 记录非敏感信息
- run: |
    echo "Deploying version: ${{ github.sha }}"
    echo "Environment: production"
    echo "Triggered by: ${{ github.actor }}"

# ✅ 使用占位符
- run: |
    echo "API Key: ***"
    echo "Connecting to database..."
```

#### 启用审计日志

利用 GitHub 的审计功能：
- 在组织设置中启用审计日志
- 定期审查工作流执行历史
- 监控异常的工作流运行

---

### 11. Pull Request 安全

#### 限制 PR 的工作流权限

```yaml
name: PR Check

on:
  pull_request_target:               # 注意：慎用 pull_request_target

permissions:
  contents: read                     # PR 检查通常只需要读权限
  pull-requests: write               # 如需评论 PR

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          ref: ${{ github.event.pull_request.head.sha }}  # 检查 PR 代码

      # ❌ 危险：不要在 pull_request_target 中执行不受信任的代码
      # - run: npm install && npm test

      # ✅ 安全：仅进行静态检查
      - run: |
          go vet ./...
          golangci-lint run
```

#### Fork 的 PR 处理

```yaml
# ✅ 对于外部贡献，使用 pull_request 而非 pull_request_target
on:
  pull_request:                      # 对 fork 的 PR 权限受限

# 如果需要使用 secrets，考虑手动审批
on:
  pull_request_target:
    types: [labeled]                 # 需要添加特定标签才触发

jobs:
  test:
    if: contains(github.event.pull_request.labels.*.name, 'safe-to-test')
    steps:
      - uses: actions/checkout@v4
```

---

### 12. 环境保护规则

#### 生产环境的保护

在仓库设置中配置环境保护规则：

```yaml
jobs:
  deploy-production:
    environment:
      name: production
      url: https://example.com
    steps:
      - run: ./deploy.sh

# 在 GitHub Settings → Environments → production 中配置：
# - Required reviewers: 需要指定人员审批
# - Wait timer: 延迟部署时间
# - Deployment branches: 限制可部署的分支
```

---

### 13. 安全检查清单

在编写 Workflow 前，检查以下项目：

- [ ] 所有敏感信息都使用 GitHub Secrets 存储
- [ ] 没有在代码中硬编码密码、API 密钥或令牌
- [ ] 使用环境变量传递用户输入，避免代码注入
- [ ] 遵循最小权限原则设置 permissions
- [ ] 第三方 Actions 使用 SHA 固定版本
- [ ] 定期审查和更新依赖
- [ ] 启用依赖项安全扫描（Dependabot、govulncheck 等）
- [ ] 容器镜像使用固定标签并进行漏洞扫描
- [ ] 日志中不包含敏感信息
- [ ] 生产环境配置了保护规则和审批流程
- [ ] 敏感文件已添加到 .gitignore
- [ ] PR 工作流不会执行不受信任的代码
- [ ] 使用 HTTPS 进行所有网络通信
- [ ] 定期轮换 secrets 和访问令牌

---

## 11. 相关资源

### 官方文档

- [GitHub Actions 官方文档](https://docs.github.com/en/actions)
- [Workflow 语法参考](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
- [上下文和表达式](https://docs.github.com/en/actions/learn-github-actions/contexts)
- [事件触发器](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows)

### 常用工具和 Actions

- [GitHub Marketplace](https://github.com/marketplace?type=actions)
- [actions/checkout](https://github.com/actions/checkout)
- [actions/setup-go](https://github.com/actions/setup-go)
- [actions/cache](https://github.com/actions/cache)
- [golangci-lint-action](https://github.com/golangci/golangci-lint-action)
- [softprops/action-gh-release](https://github.com/softprops/action-gh-release)

### 社区资源

- [Awesome Actions](https://github.com/sdras/awesome-actions)
- [GitHub Actions 示例](https://github.com/actions/starter-workflows)

---

**文档版本**：1.0
**最后更新**：2025-10-16
**维护者**：项目团队