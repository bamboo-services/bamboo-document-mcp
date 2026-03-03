// Bamboo Document MCP - Bamboo 竹简库文档服务 MCP 服务器
//
// 本服务提供从 doc.x-lf.com 获取 Bamboo 竹简库文档信息的能力，
// 支持文档目录列表、文档详情获取、文档内容搜索等功能。
package main

import (
	"context"
	"log"

	"github.com/bamboo-services/bamboo-document-mcp/route"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// main 初始化并启动 MCP 服务器，注册所有工具路由后通过标准 I/O 传输运行。
func main() {
	ctx := context.Background()
	server := mcp.NewServer(&mcp.Implementation{Name: "BambooDocument", Version: "v1.0.0"}, nil)

	// 构建路由
	rTer := route.NewRoute(ctx, server)
	rTer.RouteBuild()

	// 启动 MCP 服务器
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
