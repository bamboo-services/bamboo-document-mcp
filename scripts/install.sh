#!/bin/bash
#
# Bamboo Document MCP 一键安装脚本
# 项目地址：https://github.com/bamboo-services/bamboo-document-mcp
#
# 用法：
#   curl -fsSL https://raw.githubusercontent.com/bamboo-services/bamboo-document-mcp/master/scripts/install.sh | bash
#   或指定版本：
#   curl -fsSL https://raw.githubusercontent.com/bamboo-services/bamboo-document-mcp/master/scripts/install.sh | bash -s -- v0.0.3

set -e

# 配置
REPO="bamboo-services/bamboo-document-mcp"
BINARY_NAME="bamboo-document"
INSTALL_DIR="$HOME/.local/bin"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印函数
info() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# 检测操作系统
detect_os() {
    case "$(uname -s)" in
        Darwin*) echo "darwin" ;;
        Linux*)  echo "linux" ;;
        CYGWIN*|MINGW*|MSYS*) echo "windows" ;;
        *) error "不支持的操作系统: $(uname -s)" ;;
    esac
}

# 检测架构
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "不支持的架构: $(uname -m)" ;;
    esac
}

# 获取最新版本号
get_latest_version() {
    local version
    version=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        error "无法获取最新版本号"
    fi
    echo "$version"
}

# 主安装流程
main() {
    echo ""
    echo "╔══════════════════════════════════════════════════════╗"
    echo "║       Bamboo Document MCP 安装脚本                   ║"
    echo "║       https://github.com/bamboo-services/bamboo-document-mcp ║"
    echo "╚══════════════════════════════════════════════════════╝"
    echo ""

    # 获取版本
    VERSION="${1:-}"
    if [ -z "$VERSION" ]; then
        info "正在获取最新版本..."
        VERSION=$(get_latest_version)
    fi
    info "安装版本: $VERSION"

    # 检测系统
    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "检测到系统: $OS/$ARCH"

    # 设置文件名
    if [ "$OS" = "windows" ]; then
        BINARY="${BINARY_NAME}-${OS}-${ARCH}.exe"
    else
        BINARY="${BINARY_NAME}-${OS}-${ARCH}"
    fi

    # 下载 URL
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"

    # 创建安装目录
    mkdir -p "$INSTALL_DIR"

    # 下载二进制文件
    info "正在下载: $DOWNLOAD_URL"
    TEMP_FILE=$(mktemp)
    if ! curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_FILE"; then
        rm -f "$TEMP_FILE"
        error "下载失败，请检查版本号是否正确"
    fi

    # 下载校验和并验证
    CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
    info "正在验证文件完整性..."
    if curl -fsSL "$CHECKSUM_URL" -o /tmp/checksums.txt 2>/dev/null; then
        cd /tmp
        if grep -q "$(sha256sum "$TEMP_FILE" | cut -d' ' -f1)" checksums.txt; then
            success "SHA256 校验通过"
        else
            warn "SHA256 校验失败，但继续安装"
        fi
        rm -f checksums.txt
        cd - > /dev/null
    else
        warn "无法下载校验和文件，跳过验证"
    fi

    # 安装
    info "正在安装到: $INSTALL_DIR/$BINARY_NAME"
    mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    # 检查 PATH
    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        warn ""
        warn "安装目录 $INSTALL_DIR 不在 PATH 中"
        warn "请将以下内容添加到你的 shell 配置文件 (~/.bashrc, ~/.zshrc 等):"
        warn ""
        warn "    export PATH=\"\$HOME/.local/bin:\$PATH\""
        warn ""
        warn "然后运行: source ~/.bashrc (或 ~/.zshrc)"
    fi

    # 成功
    echo ""
    success "✅ $BINARY_NAME $VERSION 安装成功！"
    echo ""
    echo "使用方法："
    echo "    $BINARY_NAME --help"
    echo ""
    echo "文档站点：https://doc.x-lf.com"
    echo ""
}

# 运行
main "$@"
