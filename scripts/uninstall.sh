#!/bin/bash
#
# Bamboo Document MCP 卸载脚本
# 项目地址：https://github.com/bamboo-services/bamboo-document-mcp

set -e

BINARY_NAME="bamboo-document"
INSTALL_DIR="$HOME/.local/bin"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }

echo ""
echo "╔══════════════════════════════════════════════════════╗"
echo "║       Bamboo Document MCP 卸载脚本                   ║"
echo "╚══════════════════════════════════════════════════════╝"
echo ""

BINARY_PATH="$INSTALL_DIR/$BINARY_NAME"

if [ ! -f "$BINARY_PATH" ]; then
    warn "$BINARY_NAME 未安装在 $INSTALL_DIR"
    # 尝试查找其他位置
    FOUND=$(which $BINARY_NAME 2>/dev/null || true)
    if [ -n "$FOUND" ]; then
        info "发现安装位置: $FOUND"
        read -p "是否删除? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -f "$FOUND"
            success "已删除: $FOUND"
        fi
    fi
    exit 0
fi

info "正在删除: $BINARY_PATH"
rm -f "$BINARY_PATH"

success "✅ $BINARY_NAME 已成功卸载"
echo ""
