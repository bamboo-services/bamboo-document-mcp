// Package route 提供 MCP 工具的路由注册功能。
package route

import (
	"context"

	"github.com/bamboo-services/bamboo-document-mcp/tool"
	"github.com/go-resty/resty/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Route MCP 路由管理器，负责管理上下文、MCP 服务器实例和工具集合。
type Route struct {
	ctx  context.Context // ctx 上下文
	mcp  *mcp.Server     // mcp MCP 服务器实例
	tool *tool.Tool      // tool 工具集合
}

// NewRoute 创建新的路由管理器实例，初始化 HTTP 客户端和工具集合。
func NewRoute(ctx context.Context, mcp *mcp.Server) *Route {
	client := resty.New()
	return &Route{
		ctx:  ctx,
		mcp:  mcp,
		tool: tool.NewTool(client),
	}
}

// RouteBuild 注册所有 MCP 工具到服务器，包括文档列表、文档详情和文档搜索工具。
func (r *Route) RouteBuild() {
	mcp.AddTool(r.mcp, &mcp.Tool{
		Name:        "list",
		Description: "查看文档目录列表，可根据板块和关键词筛选",
	}, r.tool.DocumentList)
	mcp.AddTool(r.mcp, &mcp.Tool{
		Name:        "detail",
		Description: "获取指定文档的完整 Markdown 内容",
	}, r.tool.DocumentDetail)
	mcp.AddTool(r.mcp, &mcp.Tool{
		Name:        "search",
		Description: "在文档内容中搜索关键词，返回匹配行及上下文段落",
	}, r.tool.DocumentSearch)
}
