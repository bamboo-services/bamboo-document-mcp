// Package tool 提供所有 MCP 工具的实现，包括文档列表、详情和搜索功能。
package tool

import (
	"github.com/88250/lute"
	"github.com/go-resty/resty/v2"
)

// Tool 工具集合，封装 HTTP 客户端用于与文档服务通信。
type Tool struct {
	client *resty.Client // client HTTP 客户端
	lute   *lute.Lute
}

// NewTool 创建新的工具实例，使用指定的 HTTP 客户端。
func NewTool(client *resty.Client) *Tool {
	return &Tool{
		client: client,
		lute:   lute.New(),
	}
}
