package tool

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SectorListInput 板块列表工具的输入参数（无参数）。
type SectorListInput struct{}

// SectorList 获取所有可用的文档板块列表。
func (t *Tool) SectorList(
	_ context.Context,
	_ *mcp.CallToolRequest,
	_ SectorListInput,
) (*mcp.CallToolResult, any, error) {
	resp, err := t.client.R().Get("https://doc.x-lf.com/llms.txt")
	if err != nil {
		return nil, nil, err
	}

	// 解析文档内容提取板块
	content := resp.String()
	lines := strings.Split(content, "\n")
	re := regexp.MustCompile(`- \[.*?\]\(/docs/([a-zA-Z0-9_-]+)`)

	sectorMap := make(map[string]bool)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches != nil && len(matches) >= 2 {
			sectorMap[matches[1]] = true
		}
	}

	// 转换为列表
	sectorStr := "板块列表\n\n"
	sectorNum := 1
	for sector := range sectorMap {
		sectorStr += fmt.Sprintf("%d. %s\n", sectorNum, sector)
		sectorNum++
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: sectorStr},
		},
	}, nil, nil
}
