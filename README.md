# Bamboo Document MCP

Bamboo 文档 - 竹简库是一套跨语言的后端服务组件库，旨在提供统一风格的 API 设计和开发体验。

本项目是一个 MCP (Model Context Protocol) 服务器，用于从 [doc.x-lf.com](https://doc.x-lf.com) 获取 Bamboo 竹简库的文档信息。

## 功能特性

- 📚 **文档目录列表** - 查看指定板块的文档目录，支持关键词筛选
- 📄 **文档详情获取** - 获取指定文档的完整 Markdown 内容
- 🔍 **文档内容搜索** - 在文档中搜索关键词，返回匹配行及上下文

## 安装

### macOS / Linux (Homebrew)

通过 Homebrew 安装是最简单的方式：

```bash
# 添加 Tap
brew tap xiaolfeng/tap

# 安装
brew install bamboo-document
```

更新到最新版本：

```bash
brew update
brew upgrade bamboo-document
```

卸载：

```bash
brew uninstall bamboo-document
```

### Windows (PowerShell)

一键安装：

```powershell
iwr -useb https://raw.githubusercontent.com/bamboo-services/bamboo-document-mcp/master/scripts/install.ps1 | iex
```

指定版本安装：

```powershell
# 下载脚本
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/bamboo-services/bamboo-document-mcp/master/scripts/install.ps1" -OutFile "install.ps1"
# 执行安装
.\install.ps1 -Version "v1.0.0"
```

卸载：

```powershell
iwr -useb https://raw.githubusercontent.com/bamboo-services/bamboo-document-mcp/master/scripts/uninstall.ps1 | iex
```

### 手动构建

```bash
# 克隆仓库
git clone https://github.com/bamboo-services/bamboo-document-mcp.git
cd bamboo-document-mcp

# 安装依赖
make deps

# 构建当前平台
make build

# 构建所有平台插件包
make package-all
```

## MCP 工具

### list - 文档目录列表

查看指定板块的文档目录列表。

**参数：**
- `sector` (必填): 板块标识，如 `bamboo-base-go`
- `search` (可选): 搜索关键词

**示例：**
```json
{
  "sector": "bamboo-base-go",
  "search": "architecture"
}
```

### detail - 文档详情获取

获取指定文档的完整 Markdown 内容。

**参数：**
- `sector` (必填): 板块标识
- `path` (必填): 文档路径，如 `/architecture`

**示例：**
```json
{
  "sector": "bamboo-base-go",
  "path": "/architecture"
}
```

### search - 文档内容搜索

在文档内容中搜索关键词，返回匹配行及上下文段落。

**参数：**
- `sector` (必填): 板块标识
- `path` (必填): 文档路径
- `keyword` (必填): 搜索关键词
- `context_lines` (可选): 上下文行数，默认 3

**示例：**
```json
{
  "sector": "bamboo-base-go",
  "path": "/architecture",
  "keyword": "API",
  "context_lines": 5
}
```

## 开发

### 项目结构

```
bamboo-document-mcp/
├── main.go                 # 入口文件
├── Makefile               # 构建配置
├── go.mod                 # Go 模块定义
├── go.sum                 # 依赖校验
├── models/                # 数据模型
│   ├── list.go
│   └── search.go
├── route/                 # 路由注册
│   └── route.go
├── tool/                  # MCP 工具实现
│   ├── tool.go
│   ├── list.go
│   ├── detail.go
│   └── search.go
├── plugins/               # 插件输出目录
│   ├── document-macos-aarch/
│   ├── document-macos-amd/
│   ├── document-linux-aarch/
│   ├── document-linux-amd/
│   └── document-windows-amd/
└── .claude-plugin/        # 插件市场配置
    └── marketplace.json
```

### 常用命令

```bash
make deps      # 安装依赖
make build     # 本地构建
make test      # 运行测试
make lint      # 代码检查
make fmt       # 格式化代码
make clean     # 清理构建产物
```

## 技术栈

- **语言**: Go 1.21+
- **MCP SDK**: [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)
- **HTTP 客户端**: [github.com/go-resty/resty](https://github.com/go-resty/resty)

## 许可证

本项目采用 MIT 许可证，详见 [MIT-LICENSE](LICENSE) 文件。

## 作者

**XiaoLFeng**
