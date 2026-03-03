# 变量定义
APP_NAME := document
MAIN_FILE := .
GO := go
GOFLAGS := -v

# 版本信息
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# 插件目录定义
PLUGINS_DIR := plugins
PLUGIN_DESC := Bamboo 文档 - 竹简库是一套跨语言的后端服务组件库，旨在提供统一风格的 API 设计和开发体验

# 平台配置
# macOS
DARWIN_ARM64_DIR := document-macos-aarch
DARWIN_AMD64_DIR := document-macos-amd
# Linux
LINUX_ARM64_DIR := document-linux-aarch
LINUX_AMD64_DIR := document-linux-amd
# Windows
WINDOWS_AMD64_DIR := document-windows-amd

# 默认目标
.PHONY: all
all: build

# 本地构建 (自动检测当前平台)
.PHONY: build
build:
	@echo "🔨 构建本地版本..."
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(APP_NAME) $(MAIN_FILE)
	@echo "✅ 构建完成: $(APP_NAME)"

# macOS ARM64 (Apple Silicon)
.PHONY: package-darwin-arm64
package-darwin-arm64:
	@echo "📦 构建 macOS ARM64 插件包..."
	@mkdir -p $(PLUGINS_DIR)/$(DARWIN_ARM64_DIR)/.claude-plugin
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(PLUGINS_DIR)/$(DARWIN_ARM64_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo '{"mcpServers":{"shared-server":{"command":"./$(APP_NAME)","args":[],"env":{}}}}' > $(PLUGINS_DIR)/$(DARWIN_ARM64_DIR)/.mcp.json
	@perl -pi -e 's/"version":\s*"[^"]*"/"version": "$(VERSION)"/' $(PLUGINS_DIR)/$(DARWIN_ARM64_DIR)/.claude-plugin/plugin.json
	@echo "✅ 完成: $(PLUGINS_DIR)/$(DARWIN_ARM64_DIR)/"

# macOS AMD64 (Intel)
.PHONY: package-darwin-amd64
package-darwin-amd64:
	@echo "📦 构建 macOS AMD64 插件包..."
	@mkdir -p $(PLUGINS_DIR)/$(DARWIN_AMD64_DIR)/.claude-plugin
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(PLUGINS_DIR)/$(DARWIN_AMD64_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo '{"mcpServers":{"shared-server":{"command":"./$(APP_NAME)","args":[],"env":{}}}}' > $(PLUGINS_DIR)/$(DARWIN_AMD64_DIR)/.mcp.json
	@perl -pi -e 's/"version":\s*"[^"]*"/"version": "$(VERSION)"/' $(PLUGINS_DIR)/$(DARWIN_AMD64_DIR)/.claude-plugin/plugin.json
	@echo "✅ 完成: $(PLUGINS_DIR)/$(DARWIN_AMD64_DIR)/"

# Linux ARM64
.PHONY: package-linux-arm64
package-linux-arm64:
	@echo "📦 构建 Linux ARM64 插件包..."
	@mkdir -p $(PLUGINS_DIR)/$(LINUX_ARM64_DIR)/.claude-plugin
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(PLUGINS_DIR)/$(LINUX_ARM64_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo '{"mcpServers":{"shared-server":{"command":"./$(APP_NAME)","args":[],"env":{}}}}' > $(PLUGINS_DIR)/$(LINUX_ARM64_DIR)/.mcp.json
	@perl -pi -e 's/"version":\s*"[^"]*"/"version": "$(VERSION)"/' $(PLUGINS_DIR)/$(LINUX_ARM64_DIR)/.claude-plugin/plugin.json
	@echo "✅ 完成: $(PLUGINS_DIR)/$(LINUX_ARM64_DIR)/"

# Linux AMD64
.PHONY: package-linux-amd64
package-linux-amd64:
	@echo "📦 构建 Linux AMD64 插件包..."
	@mkdir -p $(PLUGINS_DIR)/$(LINUX_AMD64_DIR)/.claude-plugin
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(PLUGINS_DIR)/$(LINUX_AMD64_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo '{"mcpServers":{"shared-server":{"command":"./$(APP_NAME)","args":[],"env":{}}}}' > $(PLUGINS_DIR)/$(LINUX_AMD64_DIR)/.mcp.json
	@perl -pi -e 's/"version":\s*"[^"]*"/"version": "$(VERSION)"/' $(PLUGINS_DIR)/$(LINUX_AMD64_DIR)/.claude-plugin/plugin.json
	@echo "✅ 完成: $(PLUGINS_DIR)/$(LINUX_AMD64_DIR)/"

# Windows AMD64
.PHONY: package-windows-amd64
package-windows-amd64:
	@echo "📦 构建 Windows AMD64 插件包..."
	@mkdir -p $(PLUGINS_DIR)/$(WINDOWS_AMD64_DIR)/.claude-plugin
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(PLUGINS_DIR)/$(WINDOWS_AMD64_DIR)/$(APP_NAME).exe $(MAIN_FILE)
	@echo '{"mcpServers":{"shared-server":{"command":"./$(APP_NAME).exe","args":[],"env":{}}}}' > $(PLUGINS_DIR)/$(WINDOWS_AMD64_DIR)/.mcp.json
	@perl -pi -e 's/"version":\s*"[^"]*"/"version": "$(VERSION)"/' $(PLUGINS_DIR)/$(WINDOWS_AMD64_DIR)/.claude-plugin/plugin.json
	@echo "✅ 完成: $(PLUGINS_DIR)/$(WINDOWS_AMD64_DIR)/"

# 构建所有平台插件包
.PHONY: package-all
package-all: package-darwin-arm64 package-darwin-amd64 package-linux-arm64 package-linux-amd64 package-windows-amd64
	@echo ""
	@echo "🎉 所有平台插件包构建完成！"
	@echo ""
	@echo "📦 插件包列表:"
	@echo "  ├── $(DARWIN_ARM64_DIR)/  (macOS Apple Silicon)"
	@echo "  ├── $(DARWIN_AMD64_DIR)/  (macOS Intel)"
	@echo "  ├── $(LINUX_ARM64_DIR)/   (Linux ARM64)"
	@echo "  ├── $(LINUX_AMD64_DIR)/   (Linux AMD64)"
	@echo "  └── $(WINDOWS_AMD64_DIR)/ (Windows AMD64)"

# macOS 别名
.PHONY: package-macos-arm64
package-macos-arm64: package-darwin-arm64

.PHONY: package-macos-amd64
package-macos-amd64: package-darwin-amd64

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

# 清理构建产物
.PHONY: clean
clean:
	@echo "🧹 清理构建产物..."
	@rm -f $(PLUGINS_DIR)/document-*/$(APP_NAME) $(PLUGINS_DIR)/document-*/$(APP_NAME).exe
	@rm -f $(PLUGINS_DIR)/document-*/.mcp.json
	@rm -f $(APP_NAME) $(APP_NAME).exe
	@echo "✅ 清理完成"

# 帮助信息
.PHONY: help
help:
	@echo "Makefile 使用说明:"
	@echo ""
	@echo "  构建命令:"
	@echo "    make build              - 本地构建（当前平台）"
	@echo "    make package-all        - 构建所有平台插件包"
	@echo ""
	@echo "  平台插件包:"
	@echo "    make package-darwin-arm64  - macOS ARM64 (Apple Silicon)"
	@echo "    make package-darwin-amd64  - macOS AMD64 (Intel)"
	@echo "    make package-linux-arm64   - Linux ARM64"
	@echo "    make package-linux-amd64   - Linux AMD64"
	@echo "    make package-windows-amd64 - Windows AMD64"
	@echo ""
	@echo "  别名:"
	@echo "    make package-macos-arm64 = package-darwin-arm64"
	@echo "    make package-macos-amd64 = package-darwin-amd64"
	@echo ""
	@echo "  开发工具:"
	@echo "    make deps    - 安装依赖"
	@echo "    make fmt     - 格式化代码"
	@echo "    make lint    - 代码检查"
	@echo "    make test    - 运行测试"
	@echo "    make clean   - 清理构建产物"
