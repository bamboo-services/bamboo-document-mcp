// Package command 提供命令行接口的实现。
package command

import (
	"context"
	"log"

	"github.com/bamboo-services/bamboo-document-mcp/route"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

// 版本信息，编译时通过 ldflags 注入
var version = "0.0.3"

// mcpCmd MCP 服务命令，启动 MCP 服务器。
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "启动 MCP 服务器",
	Long: `启动 Bamboo 竹简库文档 MCP 服务器。

这是一个 MCP (Model Context Protocol) 服务器，用于从 doc.x-lf.com
获取 Bamboo 竹简库的文档信息，支持文档目录列表、文档详情获取、文档内容搜索等功能。

安装后，Claude Code 会自动检测并加载此 MCP 服务器提供的工具。`,
	Run: runMCP,
}

// runMCP 启动 MCP 服务器。
func runMCP(_ *cobra.Command, _ []string) {
	ctx := context.Background()
	server := mcp.NewServer(&mcp.Implementation{Name: "BambooDocument", Version: version}, nil)

	// 构建路由
	rTer := route.NewRoute(ctx, server)
	rTer.RouteBuild()

	// 启动 MCP 服务器
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

// GetMCPCmd 返回 MCP 命令实例。
func GetMCPCmd() *cobra.Command {
	return mcpCmd
}

// SetVersion 设置版本号。
func SetVersion(v string) {
	version = v
}
