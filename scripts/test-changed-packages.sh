#!/usr/bin/env bash

# 脚本用于检测变更的包并只运行这些包的测试
# 适用于扁平化组织的 Go 组件库

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取变更的文件列表
get_changed_files() {
    local base_ref="${1:-origin/main}"
    
    # 检查是否在 CI 环境中
    if [ -n "$GITHUB_BASE_REF" ]; then
        # PR 场景：对比 PR 的 base 分支
        base_ref="origin/$GITHUB_BASE_REF"
        echo -e "${YELLOW}检测到 PR 环境，对比分支: $base_ref${NC}" >&2
    elif [ -n "$GITHUB_REF" ] && [ "$GITHUB_EVENT_NAME" = "push" ]; then
        # Push 场景：对比上一次提交
        echo -e "${YELLOW}检测到 Push 环境，对比上一次提交${NC}" >&2
        git diff --name-only HEAD~1 HEAD
        return
    fi
    
    # 确保 fetch 了远程分支
    git fetch origin "$GITHUB_BASE_REF" 2>/dev/null || true
    
    # 获取变更的文件
    git diff --name-only "$base_ref"...HEAD
}

# 从文件路径提取包路径
extract_packages() {
    local changed_files="$1"
    local packages=()
    
    while IFS= read -r file; do
        # 跳过非 Go 文件和根目录文件
        if [[ ! "$file" =~ \.go$ ]] || [[ ! "$file" =~ / ]]; then
            continue
        fi
        
        # 提取包路径（第一级目录）
        # 例如: config/config.go -> config
        #      errs/stack/stack.go -> errs
        #      utils/encrypt/encryption.go -> utils
        local pkg_path=$(echo "$file" | cut -d'/' -f1)
        
        # 去重添加
        if [[ ! " ${packages[@]} " =~ " ${pkg_path} " ]]; then
            packages+=("$pkg_path")
        fi
    done <<< "$changed_files"
    
    printf '%s\n' "${packages[@]}"
}

# 主函数
main() {
    local base_ref="${1:-origin/main}"
    local test_all=false
    
    echo -e "${GREEN}=== 检测变更的包 ===${NC}"
    
    # 获取变更的文件
    changed_files=$(get_changed_files "$base_ref")
    
    if [ -z "$changed_files" ]; then
        echo -e "${YELLOW}未检测到变更的文件，将测试所有包${NC}"
        test_all=true
    else
        echo -e "${GREEN}变更的文件:${NC}"
        echo "$changed_files" | sed 's/^/  /'
        echo ""
    fi
    
    # 提取变更的包
    if [ "$test_all" = false ]; then
        packages=$(extract_packages "$changed_files")
        
        if [ -z "$packages" ]; then
            echo -e "${YELLOW}变更的文件不包含 Go 包，跳过测试${NC}"
            exit 0
        fi
        
        echo -e "${GREEN}需要测试的包:${NC}"
        echo "$packages" | sed 's/^/  /'
        echo ""
        
        # 构建测试路径
        test_paths=""
        while IFS= read -r pkg; do
            if [ -d "$pkg" ]; then
                test_paths="$test_paths ./$pkg/..."
            fi
        done <<< "$packages"
        
        if [ -z "$test_paths" ]; then
            echo -e "${YELLOW}没有找到有效的包目录，跳过测试${NC}"
            exit 0
        fi
        
        echo -e "${GREEN}=== 运行测试 ===${NC}"
        echo -e "${YELLOW}测试命令: go test $test_paths -race -coverprofile=coverage.out -covermode=atomic${NC}"
        echo ""
        
        # 运行测试
        go test $test_paths -race -coverprofile=coverage.out -covermode=atomic
        
        echo ""
        echo -e "${GREEN}=== 测试完成 ===${NC}"
        
        # 显示覆盖率
        if [ -f coverage.out ]; then
            echo -e "${GREEN}覆盖率统计:${NC}"
            go tool cover -func=coverage.out | tail -n 1
        fi
    else
        # 测试所有包
        echo -e "${GREEN}=== 运行所有测试 ===${NC}"
        go test ./... -race -coverprofile=coverage.out -covermode=atomic
        
        echo ""
        echo -e "${GREEN}=== 测试完成 ===${NC}"
        
        # 显示覆盖率
        if [ -f coverage.out ]; then
            echo -e "${GREEN}覆盖率统计:${NC}"
            go tool cover -func=coverage.out | tail -n 1
        fi
    fi
}

# 执行主函数
main "$@"
