# 变量定义
APP_NAME := bamboo-document
MAIN_FILE := .
GO := go
GOFLAGS := -v

# 版本信息
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

# 帮助信息
.PHONY: help
help:
	@echo "Bamboo Document MCP - Makefile 使用说明:"
	@echo ""
	@echo "  构建命令:"
	@echo "    make build          - 本地构建（当前平台）"
	@echo "    make run            - 构建并运行"
	@echo ""
	@echo "  GoReleaser 命令:"
	@echo "    make release-local - 本地快照构建（所有平台）"
	@echo "    make release-build - 单一目标构建"
	@echo ""
	@echo "  开发工具:"
	@echo "    make deps          - 安装依赖"
	@echo "    make fmt           - 格式化代码"
	@echo "    make lint          - 代码检查"
	@echo "    make test          - 运行测试"
	@echo "    make clean         - 清理构建产物"

# 本地构建
.PHONY: build
build:
	@echo "🔨 构建本地版本..."
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(APP_NAME) $(MAIN_FILE)
	@echo "✅ 构建完成: $(APP_NAME)"

# 安装依赖
.PHONY: deps
deps:
	@echo "📦 下载依赖..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✅ 依赖安装完成"

# 代码格式化
.PHONY: fmt
fmt:
	@echo "📝 格式化代码..."
	$(GO) fmt ./...
	@echo "✅ 格式化完成"

# 代码检查
.PHONY: lint
lint:
	@echo "🔍 代码检查..."
	$(GO) vet ./...
	@echo "✅ 检查完成"

# 运行测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	$(GO) test -v ./...

# 本地运行（开发模式）
.PHONY: run
run: build
	@echo "🚀 启动 MCP 服务器..."
	./$(APP_NAME)

# 清理构建产物
.PHONY: clean
clean:
	@echo "🧹 清理构建产物..."
	@rm -f $(APP_NAME) $(APP_NAME).exe
	@echo "✅ 清理完成"

# 使用 GoReleaser 本地构建（需要安装 goreleaser）
.PHONY: release-local
release-local:
	@echo "📦 使用 GoReleaser 本地构建..."
	goreleaser release --snapshot --clean
	@echo "✅ 本地发布构建完成"

# 使用 GoReleaser 构建单一二进制
.PHONY: release-build
release-build:
	@echo "📦 使用 GoReleaser 构建..."
	goreleaser build --single-target --clean
	@echo "✅ 构建完成"
