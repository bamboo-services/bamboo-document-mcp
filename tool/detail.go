package tool

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DocumentDetailInput 文档详情工具的输入参数。
type DocumentDetailInput struct {
	Sector string `json:"sector" jsonschema:"required,板块标识，如 bamboo-base-go"` // Sector 板块标识
	Path   string `json:"path" jsonschema:"required,文档路径，如 /architecture"`    // Path 文档路径
}

// DocumentDetail 获取指定文档的完整 Markdown 内容。
func (t *Tool) DocumentDetail(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input DocumentDetailInput,
) (*mcp.CallToolResult, any, error) {
	// 1. 规范化 path：移除 /docs/{sector} 前缀（如果存在），确保以 / 开头
	path := input.Path
	sectorPrefix := "/docs/" + input.Sector
	if strings.HasPrefix(path, sectorPrefix) {
		path = strings.TrimPrefix(path, sectorPrefix)
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 2. 构建请求 URL
	uri := fmt.Sprintf("https://doc.x-lf.com/llms.mdx/docs/%s%s", input.Sector, path)

	// 3. 发送 GET 请求
	resp, err := t.client.R().Get(uri)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "文档获取失败"},
			},
		}, nil, err
	}

	// 4. 检查响应状态码
	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("文档不存在或获取失败: HTTP %d", resp.StatusCode())},
			},
		}, nil, fmt.Errorf("文档不存在或获取失败: HTTP %d", resp.StatusCode())
	}

	// 5. 提取标题
	content := resp.String()
	title := extractTitle(content)

	// 6. 格式化内容
	formattedContent := t.lute.FormatStr("detail", content)

	// 7. 构建返回结果（使用多个 TextContent）
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("板块: %s", input.Sector)},
			&mcp.TextContent{Text: fmt.Sprintf("路径: %s", input.Path)},
			&mcp.TextContent{Text: fmt.Sprintf("标题: %s", title)},
			&mcp.TextContent{Text: fmt.Sprintf("文档地址: %s", uri)},
			&mcp.TextContent{Text: formattedContent},
		},
	}, nil, nil
}

// extractTitle 从 Markdown 内容中提取第一个标题
func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}
