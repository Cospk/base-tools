SHELL := /usr/bin/env bash

# 检测 GOBIN（为空时回退到 $HOME/go/bin）
GOBIN := $(shell go env GOBIN)
ifeq ($(strip $(GOBIN)),)
  GOBIN := $(HOME)/go/bin
endif

.PHONY: all tidy fmt vet lint test cover tools.verify.go-gitlint tools.install.golangci tools.install.go-gitlint

all: tidy fmt vet lint test cover

## 整理并校验 Go 模块依赖
tidy:
	go mod tidy
	go mod verify

## 代码格式化
fmt:
	go fmt ./...

## Go 静态检查（vet）
vet:
	go vet ./...

## 若未安装 golangci-lint 则自动安装
tools.install.golangci:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		mkdir -p $(GOBIN); \
		GO111MODULE=on GOBIN=$(GOBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	else \
		echo "golangci-lint already installed"; \
	fi

## 运行 golangci-lint
lint: tools.install.golangci
	$(GOBIN)/golangci-lint run --timeout=5m || golangci-lint run --timeout=5m

## 若未安装 go-gitlint 则自动安装（go install 生成的二进制名为 go-gitlint）
tools.install.go-gitlint:
	@if ! command -v go-gitlint >/dev/null 2>&1; then \
		echo "Installing go-gitlint..."; \
		mkdir -p $(GOBIN); \
		GO111MODULE=on GOBIN=$(GOBIN) go install github.com/llorllale/go-gitlint@latest; \
	else \
		echo "go-gitlint already installed"; \
	fi

## 使用 go-gitlint 校验提交信息
tools.verify.go-gitlint: tools.install.go-gitlint
	@if [ -x "$(GOBIN)/go-gitlint" ]; then \
		"$(GOBIN)/go-gitlint"; \
	elif command -v go-gitlint >/dev/null 2>&1; then \
		go-gitlint; \
	elif [ -x "$(GOBIN)/gitlint" ]; then \
		"$(GOBIN)/gitlint"; \
	elif command -v gitlint >/dev/null 2>&1; then \
		gitlint; \
	else \
		echo "未找到 go-gitlint/gitlint，请检查安装与 GOBIN（当前：$(GOBIN)）"; \
		exit 127; \
	fi

## 启用竞态检测并生成覆盖率文件的测试
test:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

## 显示覆盖率摘要（若不存在则先生成）
cover:
	@if [ ! -f coverage.out ]; then \
		echo "coverage.out not found; running tests to generate it"; \
		$(MAKE) test; \
	fi
	go tool cover -func=coverage.out | tail -n 1
