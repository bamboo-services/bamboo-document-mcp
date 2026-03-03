package tool

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DocumentListInput 文档列表工具的输入参数。
type DocumentListInput struct {
	Sector string  `json:"sector" jsonschema:"required,板块内容"`  // Sector 板块标识，如 bamboo-base-go
	Search *string `json:"search" jsonschema:"optional,搜索关键词"` // Search 可选的搜索关键词
}

// DocumentList 获取指定板块的文档目录列表，支持按关键词筛选。
func (t *Tool) DocumentList(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input DocumentListInput,
) (*mcp.CallToolResult, any, error) {
	resp, err := t.client.R().Get("https://doc.x-lf.com/llms.txt")
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "文档列表获取失败"},
			},
		}, nil, err
	}

	// 格式化文本
	newStr := strings.Replace(resp.String(), "# Documentation\n\n", "", 1)
	lines := strings.Split(newStr, "\n")
	re := regexp.MustCompile(`- \[(.*?)]\((.*?)\):\s*(.*)`)

	listStr := "文档列表\n\n"
	listNum := 1
	for _, line := range lines {
		// 解析
		matches := re.FindStringSubmatch(line)
		if matches == nil || len(matches) != 4 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "格式无法解析"},
				},
			}, nil, fmt.Errorf("格式无法解析")
		}

		// 筛选
		if !strings.HasPrefix(matches[2], "/docs/"+input.Sector) {
			continue
		}
		if input.Search != nil && !strings.Contains(strings.ToLower(matches[1]), strings.ToLower(*input.Search)) && !strings.Contains(strings.ToLower(matches[3]), strings.ToLower(*input.Search)) {
			continue
		}

		// 组合器
		listStr += fmt.Sprintf("%d. [%s](%s): %s\n", listNum, matches[1], matches[2], matches[3])
		listNum++
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "基础地址: https://doc.x-lf.com\n"},
			&mcp.TextContent{Text: listStr},
		},
	}, nil, nil
}
