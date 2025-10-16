# .github 目录说明文档

本文档详细说明了 `.github` 目录下各个文件的作用以及 GitHub Actions Workflows 的语法使用。

## 目录结构

```
.github/
├── codecov.yml                      # Codecov 代码覆盖率配置文件
└── workflows/                       # GitHub Actions 工作流目录
    ├── auto-tag.yml                 # 通过评论手动创建 Tag 工作流
    ├── coreCi.yml                   # 核心 CI/CD 流程
    ├── go-typeCheck.yml             # Go 类型检查工作流
    ├── golangci-lint.yml            # Go 代码质量检查工作流
    ├── gosec.yml                    # Go 安全扫描工作流
    ├── release-drafter.yaml         # 自动生成 Release Notes 草稿工作流
    └── release-on-tag.yml           # Tag 推送自动发布 Release 工作流
```

---

## 文件详解

### 1. codecov.yml

**文件作用**：配置 Codecov 代码覆盖率报告的行为和规则。

**详细说明**：

```yaml
coverage:
  status:
    project:                         # 项目整体覆盖率设置
      default:
        target: auto                 # 自动设置覆盖率目标
        informational: true          # 仅提供信息，不阻止 PR 合并
    patch:                           # PR 补丁覆盖率设置
      default:
        target: auto                 # 自动设置补丁覆盖率目标
        informational: true          # 仅提供信息，不影响 PR 状态

ignore:                              # 忽略的文件类型
  - "**/*.md"                        # 忽略所有 Markdown 文件
  - "**/*.sh"                        # 忽略所有 Shell 脚本
  - "**/*.yaml"                      # 忽略所有 YAML 文件
  - "**/*.yml"                       # 忽略所有 YML 文件
  - "./.github/**"                   # 忽略 .github 目录下所有文件
```

**关键配置说明**：
- `target: auto`：Codecov 会根据历史数据自动设置覆盖率目标
- `informational: true`：表示覆盖率检查结果仅供参考，不会阻止 PR 合并
- `ignore`：指定不计入覆盖率统计的文件模式

---

### 2. workflows/coreCi.yml

**文件作用**：核心 CI/CD 流程，负责代码构建、测试和覆盖率收集。

**详细说明**：

```yaml
name: base-tools ci Auto Build and Install
```
工作流的名称，会显示在 GitHub Actions 页面。

#### 触发条件 (on)

```yaml
on:
  push:
    branches: [ main ]               # 当推送到 main 分支时触发
  pull_request:
    branches: [ main ]               # 当创建针对 main 分支的 PR 时触发
```

#### 环境变量 (env)

```yaml
env:
  GO_VERSION: 1.23.4                 # 定义全局环境变量 GO 版本
```

#### 任务定义 (jobs)

```yaml
jobs:
  cospk-base-tools:                  # 任务名称
    name: Test with go ${{ matrix.go_version }} on ${{ matrix.os }}
    runs-on: ${{ matrix.os }}        # 运行环境，使用矩阵变量
    permissions:                      # 权限设置
      contents: write                 # 允许写入仓库内容
      pull-requests: write            # 允许写入 PR
```

**权限说明**：
- `contents: write`：允许工作流修改仓库内容
- `pull-requests: write`：允许工作流在 PR 上添加评论等操作

#### 矩阵策略 (strategy)

```yaml
strategy:
  matrix:
    go_version: [ "1.23.4" ]         # Go 版本矩阵
    os: [ ubuntu-latest ]            # 操作系统矩阵
```

**矩阵说明**：可以定义多个值，GitHub Actions 会为每个组合创建一个任务实例。

#### 步骤详解 (steps)

**步骤 1：检出代码**
```yaml
- name: Setup
  uses: actions/checkout@v4          # 使用 GitHub 官方 checkout action v4
```

**步骤 2：设置 Go 环境**
```yaml
- name: Set up Go ${{ matrix.go_version }}
  uses: actions/setup-go@v5          # 使用 GitHub 官方 setup-go action v5
  with:
    go-version: ${{ matrix.go_version }}
  id: go                             # 步骤 ID，可在后续步骤中引用
```

**步骤 3：模块操作**
```yaml
- name: Module Operations
  run: |                             # 执行多行 shell 命令
    sudo make tidy                   # 整理 Go 模块依赖
    sudo make tools.verify.go-gitlint  # 验证 git commit 信息
```

**步骤 4：代码格式化**
```yaml
- name: Format Code
  run: sudo make lint
  continue-on-error: true            # 即使此步骤失败，也继续执行后续步骤
```

**步骤 5：运行测试**
```yaml
- name: test
  run: sudo make test                # 执行测试
```

**步骤 6：收集覆盖率**
```yaml
- name: Collect and Display Test Coverage
  id: collect_coverage
  run: |
    sudo make cover                  # 生成覆盖率报告
```

---

### 3. workflows/go-typeCheck.yml

**文件作用**：执行 Go 代码的类型检查，确保代码类型安全。

**详细说明**：

```yaml
name: Go Typecheck Workflow Test

on: [push, pull_request]             # 简化语法：任何推送和 PR 都触发
```

**触发条件语法说明**：
- `[push, pull_request]`：数组形式，表示多个事件都会触发
- 等同于：
  ```yaml
  on:
    push:
    pull_request:
  ```

```yaml
jobs:
  go-language-typecheck:
    runs-on: ubuntu-latest           # 在最新版 Ubuntu 上运行
    steps:
      - name: Checkout Repository
        uses: actions/checkout@v4    # 检出代码

      - name: Go Code Language typecheck
        uses: kubecub/typecheck@main # 使用第三方 Action 进行类型检查
```

**第三方 Action 说明**：
- `kubecub/typecheck@main`：使用 kubecub/typecheck 仓库的 main 分支
- `@main`：指定使用的版本/分支

---

### 4. workflows/golangci-lint.yml

**文件作用**：使用 golangci-lint 进行代码质量检查和静态分析。

**详细说明**：

```yaml
name: Cospk base-tools golangCi-lint

on:
  push:
    branches: [main]                 # 仅在推送到 main 分支时触发
  pull_request:                      # 所有 PR 都触发（不限制分支）
```

```yaml
jobs:
  golangci:
    name: lint                       # 任务显示名称
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4    # 简化语法：省略 name

      - uses: actions/setup-go@v5
        with:
          go-version: 1.23.4
          cache: false               # 禁用 Go 缓存
```

**缓存说明**：
- `cache: false`：禁用 setup-go 的默认缓存行为
- 可能原因：与 golangci-lint 自己的缓存机制冲突

```yaml
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4.0.0
        with:
          version: v1.54             # golangci-lint 版本

          # Optional: working directory, useful for monorepos
          # working-directory: server

          # Optional: golangci-lint command line arguments.
          # Note: by default the `.golangci.yml` file should be at the root of the repository.
          # args: --timeout=30m --config=/scripts/.golangci.yml --issues-exit-code=0

          only-new-issues: true      # 仅显示新问题（PR 模式）

          # Optional:The mode to install golangci-lint.
          # install-mode: "goinstall"
```

**配置参数说明**：
- `version`：指定 golangci-lint 的版本
- `only-new-issues: true`：仅报告 PR 中新增的问题，不报告已存在的问题
- 注释中的可选参数：
  - `working-directory`：适用于 monorepo（单仓多项目）场景
  - `args`：传递额外的命令行参数
  - `install-mode`：安装模式选择

---

### 5. workflows/gosec.yml

**文件作用**：使用 gosec 工具进行 Go 代码安全扫描，发现潜在的安全漏洞。

**头部注释说明**：

```yaml
# gosec 是一个用于 Go 语言的源代码安全审计工具。它通过静态分析 Go 代码，查找潜在的安全问题。
# gosec 的主要功能包括：
#
# 1、发现常见的安全漏洞，例如 SQL 注入、命令注入以及跨站脚本攻击（XSS）。
# 2、根据常见的安全标准对代码进行审计，查找不符合规范的代码。
# 3、帮助 Go 语言工程师编写安全可靠的代码。
#
# https://github.com/securego/gosec/
```

**详细说明**：

```yaml
name: OpenIM Run Gosec

on:
  push:
    branches: "*"                    # 推送到任何分支都触发
  pull_request:
    branches: "*"                    # 针对任何分支的 PR 都触发
    paths-ignore:                    # 忽略特定文件变更
      - '*.md'                       # 忽略 Markdown 文件
      - '*.yml'                      # 忽略 YML 文件
      - '.github'                    # 忽略 .github 目录
```

**paths-ignore 说明**：
- 当 PR 仅修改了列表中的文件时，不触发此工作流
- 可以减少不必要的工作流运行，节省资源

```yaml
jobs:
  golang-security-action:
    runs-on: ubuntu-latest
    env:
      GO111MODULE: on                # 启用 Go Modules
    steps:
      - name: Check out code
        uses: actions/checkout@v4

      - name: Run Gosec Security Scanner
        uses: securego/gosec@master  # 使用 gosec 的 master 分支
        with:
          args: ./...                # 扫描所有包（递归）
        continue-on-error: true      # 发现安全问题不中断流程
```

**关键配置**：
- `args: ./...`：Go 语言的标准语法，表示递归扫描当前目录及所有子目录
- `continue-on-error: true`：即使发现安全问题也不失败，仅作为警告

---

### 6. workflows/release-drafter.yaml

**文件作用**：自动为 PR 添加标签，并根据合并的 PR 自动生成 Release Notes 草稿。

**详细说明**：

```yaml
name: Release Drafter

on:
  push:
    branches:
      - main                         # 推送到 main 分支时触发
  pull_request:                      # PR 事件用于自动标签功能
```

**触发事件说明**：
- `push` 事件：当代码合并到 main 分支时，更新 Release 草稿
- `pull_request` 事件：为 PR 自动添加标签（可选功能）
- 注释中提到的 `pull_request_target`：支持来自 fork 的 PR，但需要额外的安全考虑

#### 权限配置

```yaml
permissions:
  contents: read                     # 顶层默认权限为只读

jobs:
  update_release_draft:
    permissions:
      contents: write                # 创建 Release 需要写权限
      pull-requests: write           # 自动标签功能需要写权限
    runs-on: ubuntu-latest
```

**权限层级说明**：
- 顶层 `permissions`：设置工作流默认权限为只读（最小权限原则）
- 任务级 `permissions`：为特定任务授予必要的写权限
- 这是一种安全最佳实践：默认限制，按需授权

#### 步骤详解

**使用 Release Drafter Action**
```yaml
steps:
  - uses: release-drafter/release-drafter@v6
    env:
      GITHUB_TOKEN: ${{ secrets.REDBOT_GITHUB_TOKEN }}
```

**配置说明**：
- `release-drafter/release-drafter@v6`：第三方 Action，用于自动生成 Release Notes
- `GITHUB_TOKEN`：使用自定义的 GitHub Token（REDBOT_GITHUB_TOKEN）
  - 如果使用默认 `${{ secrets.GITHUB_TOKEN }}`，则无需在 Secrets 中配置
  - 自定义 Token 可以有更长的有效期或特定的权限

**可选配置（注释部分）**：

```yaml
# 1. GitHub Enterprise 配置
# - name: 设置 GHE_HOST
#   run: |
#     echo "GHE_HOST=${GITHUB_SERVER_URL##https:\/\/}" >> $GITHUB_ENV
```
仅在使用 GitHub Enterprise Server 时需要。

```yaml
# 2. 自定义配置文件
# with:
#   config-name: my-config.yml      # 指定配置文件名（默认为 release-drafter.yml）
#   disable-autolabeler: true       # 禁用自动标签功能
```

#### Release Drafter 工作原理

1. **监听 PR 合并**：当 PR 合并到 main 分支时触发
2. **分析 PR 信息**：
   - PR 标题
   - PR 标签（如 `feature`、`bug`、`enhancement`）
   - PR 描述
3. **生成/更新 Release 草稿**：
   - 按类别分组（Features、Bug Fixes、Dependencies 等）
   - 包含贡献者列表
   - 自动生成版本号（基于语义化版本）
4. **自动标签**（可选）：
   - 根据 PR 内容自动添加标签
   - 例如：检测到 `fix:` 前缀则添加 `bug` 标签

#### 配置文件示例

Release Drafter 需要配置文件（通常位于 `.github/release-drafter.yml`）：

```yaml
# 示例配置（需要单独创建）
name-template: 'v$RESOLVED_VERSION'
tag-template: 'v$RESOLVED_VERSION'

categories:
  - title: '🚀 Features'
    labels:
      - 'feature'
      - 'enhancement'
  - title: '🐛 Bug Fixes'
    labels:
      - 'bug'
      - 'fix'
  - title: '📚 Documentation'
    labels:
      - 'documentation'
  - title: '🔧 Maintenance'
    labels:
      - 'chore'
      - 'dependencies'

change-template: '- $TITLE @$AUTHOR (#$NUMBER)'

autolabeler:
  - label: 'bug'
    branch:
      - '/fix\/.+/'
    title:
      - '/^fix/i'
  - label: 'feature'
    branch:
      - '/feature\/.+/'
    title:
      - '/^feat/i'
```

#### 使用场景

**适用于**：
- 需要维护清晰的 Release Notes 的项目
- 使用语义化版本管理的项目
- 多人协作的开源项目
- 需要自动化发布流程的项目

**优势**：
- 自动化 Release Notes 编写，减少手动工作
- 统一的发布日志格式
- 自动识别 PR 类型并分类
- 支持贡献者名单自动生成

#### 注释中的高级特性

**1. PR 事件类型过滤**
```yaml
# types: [opened, reopened, synchronize]
```
- `opened`：PR 首次创建
- `reopened`：PR 重新打开
- `synchronize`：PR 有新的提交

**2. 支持 Fork 仓库的 PR**
```yaml
# pull_request_target:
#   types: [opened, reopened, synchronize]
```
- `pull_request_target` 事件在基础仓库的上下文中运行
- 可以访问仓库的 secrets
- 需要注意安全风险（恶意 PR 可能利用此权限）

---

### 7. workflows/release-on-tag.yml

**文件作用**：监听 tag 推送事件，自动创建正式的 GitHub Release 版本。

**详细说明**：

```yaml
name: Release on Tag Push

on:
  push:
    tags:
      - 'v*.*.*'                     # 匹配语义化版本号标签
```

**触发条件说明**：
- 仅在推送符合 `v*.*.*` 格式的 tag 时触发
- 例如：`v1.0.0`、`v2.3.4`、`v10.20.30` 等
- 不会匹配：`v1.0`、`1.0.0`（缺少 v 前缀）、`beta-1.0.0` 等

#### 权限配置

```yaml
permissions:
  contents: write                    # 创建 Release 需要写权限
```

**权限说明**：
- `contents: write`：允许工作流创建 Release 和上传资源

#### 任务步骤详解

**步骤 1：检出代码**
```yaml
- name: Checkout code
  uses: actions/checkout@v4
  with:
    fetch-depth: 0                   # 获取完整的 Git 历史
```
- `fetch-depth: 0`：获取所有提交历史，用于生成完整的 changelog

**步骤 2：提取 Tag 信息**
```yaml
- name: Get tag information
  id: tag_info
  run: |
    TAG_NAME=${GITHUB_REF#refs/tags/}
    echo "tag_name=$TAG_NAME" >> $GITHUB_OUTPUT
    echo "📦 准备发布版本: $TAG_NAME"
```
- 从 `GITHUB_REF` 中提取 tag 名称
- 将 tag 名称保存到输出变量，供后续步骤使用

**步骤 3：创建 Release**
```yaml
- name: Create Release
  uses: softprops/action-gh-release@v1
  with:
    tag_name: ${{ steps.tag_info.outputs.tag_name }}
    name: Release ${{ steps.tag_info.outputs.tag_name }}
    generate_release_notes: true     # 自动生成 Release Notes
    draft: false                     # 正式发布，非草稿
    prerelease: false                # 正式版本，非预发布
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**配置参数说明**：
- `tag_name`：关联的 tag 名称
- `name`：Release 标题
- `generate_release_notes: true`：GitHub 自动根据 commits 和 PR 生成 Release Notes
- `draft: false`：直接发布，不创建草稿
- `prerelease: false`：标记为正式版本（如果是 `v1.0.0-beta.1` 这样的预发布版本，可设为 `true`）

**步骤 4 & 5：通知步骤**
```yaml
- name: Release Success Notification
  if: success()
  run: |
    echo "✅ Release ${{ steps.tag_info.outputs.tag_name }} 创建成功！"
    echo "🔗 访问: https://github.com/${{ github.repository }}/releases/tag/${{ steps.tag_info.outputs.tag_name }}"

- name: Release Failed Notification
  if: failure()
  run: |
    echo "❌ Release 创建失败，请检查日志"
```
- `if: success()`：仅在前序步骤成功时执行
- `if: failure()`：仅在前序步骤失败时执行

#### 使用方法

**本地创建并推送 tag**：
```bash
# 创建 tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# 推送 tag 到远程仓库
git push origin v1.0.0

# 或者推送所有 tag
git push origin --tags
```

**GitHub 网页创建 Release**：
1. 进入仓库页面 → Releases → Create a new release
2. 输入 tag 名称（如 `v1.0.1`）
3. 点击 "Publish release"
4. 工作流会自动触发

#### 工作流程

```
开发者推送 tag (v1.0.0)
       ↓
GitHub 检测到 tag 推送
       ↓
触发 release-on-tag.yml 工作流
       ↓
检出代码并获取完整历史
       ↓
提取 tag 信息
       ↓
调用 GitHub API 创建 Release
       ↓
自动生成 Release Notes
       ↓
发布正式版本 ✅
```

#### 生成的 Release 包含

1. **Release 标题**：`Release v1.0.0`
2. **自动生成的 Release Notes**：
   - 自上一个 tag 以来的所有 commits
   - 相关的 Pull Requests
   - 贡献者列表
3. **下载资源**：
   - 源代码压缩包（.zip 和 .tar.gz）
   - 可以手动添加编译后的二进制文件

#### 与 release-drafter.yaml 的区别

| 特性 | release-drafter.yaml | release-on-tag.yml |
|------|---------------------|-------------------|
| 触发时机 | PR 合并到 main | 推送 tag |
| Release 状态 | 草稿（Draft） | 正式发布 |
| Release Notes | 基于 PR 标签分类 | 基于 commits 自动生成 |
| 使用场景 | 持续更新草稿 | 正式发布版本 |
| 版本控制 | 自动递增 | 手动指定 tag |

#### 最佳实践

1. **语义化版本管理**：
   - 主版本号：不兼容的 API 修改（v2.0.0）
   - 次版本号：向下兼容的功能新增（v1.1.0）
   - 修订号：向下兼容的 bug 修复（v1.0.1）

2. **Tag 命名规范**：
   ```bash
   v1.0.0        # 正式版本
   v1.0.0-beta.1 # Beta 测试版
   v1.0.0-rc.1   # Release Candidate
   v1.0.0-alpha  # Alpha 测试版
   ```

3. **组合使用建议**：
   - 使用 `release-drafter.yaml` 维护 Release 草稿
   - 准备发布时，手动编辑草稿并创建 tag
   - 或直接推送 tag，触发 `release-on-tag.yml` 自动发布

4. **保护 tag**：
   - 在 GitHub 仓库设置中启用 tag 保护规则
   - 限制谁可以创建和删除 tag

---

## 相关文档

- **[GitHub Actions Workflows 语法参考手册](./WORKFLOWS_SYNTAX.md)** - 详细的 Workflows 语法说明和最佳实践
- [GitHub Actions 官方文档](https://docs.github.com/en/actions)
- [golangci-lint 文档](https://golangci-lint.run/)
- [gosec 项目](https://github.com/securego/gosec)
- [Codecov 文档](https://docs.codecov.com/)

---

**文档版本**：2.0
**最后更新**：2025-10-16
**变更说明**：将 Workflows 语法部分拆分到独立文档 `WORKFLOWS_SYNTAX.md`