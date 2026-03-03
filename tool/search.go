package tool

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DocumentSearchInput 文档搜索工具的输入参数。
type DocumentSearchInput struct {
	Sector       string `json:"sector" jsonschema:"required,板块标识"`                   // Sector 板块标识
	Path         string `json:"path" jsonschema:"required,文档路径"`                     // Path 文档路径
	Keyword      string `json:"keyword" jsonschema:"required,搜索关键词"`                 // Keyword 搜索关键词
	ContextLines *int   `json:"context_lines" jsonschema:"optional,上下文行数，默认为 5,nil"` // ContextLines 上下文行数
}

// DocumentSearch 在指定文档内容中搜索关键词，返回匹配行及上下文段落。
func (t *Tool) DocumentSearch(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input DocumentSearchInput,
) (*mcp.CallToolResult, any, error) {
	// 1. 设置默认上下文行数
	contextLines := 5
	if input.ContextLines != nil && *input.ContextLines > 0 {
		contextLines = *input.ContextLines
	}

	// 2. 确保 path 以 / 开头
	path := input.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 3. 获取文档内容
	uri := fmt.Sprintf("https://doc.x-lf.com/llms.mdx/docs/%s%s", input.Sector, path)
	resp, err := t.client.R().Get(uri)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "文档获取失败"},
			},
		}, nil, err
	}
	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("文档不存在: HTTP %d", resp.StatusCode())},
			},
		}, nil, fmt.Errorf("文档不存在: HTTP %d", resp.StatusCode())
	}

	// 4. 按行分割并搜索
	content := resp.String()
	lines := strings.Split(content, "\n")
	keywordLower := strings.ToLower(input.Keyword)

	// 5. 构建搜索结果（动态构造 Content）
	contents := make([]mcp.Content, 0)
	contents = append(contents,
		&mcp.TextContent{Text: fmt.Sprintf("板块: %s", input.Sector)},
		&mcp.TextContent{Text: fmt.Sprintf("路径: %s", input.Path)},
		&mcp.TextContent{Text: fmt.Sprintf("关键词: %s", input.Keyword)},
		&mcp.TextContent{Text: fmt.Sprintf("文档地址: %s", uri)},
	)

	matchNum := 0
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), keywordLower) {
			matchNum++
			// 构建上下文
			contextStart := max(0, i-contextLines)
			contextEnd := min(len(lines)-1, i+contextLines)

			var contextBuilder strings.Builder
			for j := contextStart; j <= contextEnd; j++ {
				prefix := "  "
				if j == i {
					prefix = ">>"
				}
				contextBuilder.WriteString(fmt.Sprintf("%s %d: %s\n", prefix, j+1, lines[j]))
			}

			// 每个匹配作为独立的 TextContent
			contents = append(contents,
				&mcp.TextContent{Text: fmt.Sprintf("--- 匹配 %d (第 %d 行) ---", matchNum, i+1)},
				&mcp.TextContent{Text: t.lute.FormatStr("search", contextBuilder.String())},
			)
		}
	}

	// 6. 添加汇总信息
	if matchNum == 0 {
		contents = append(contents, &mcp.TextContent{Text: "未找到匹配内容"})
	} else {
		// 在开头插入汇总信息
		summary := &mcp.TextContent{Text: fmt.Sprintf("共找到 %d 个匹配", matchNum)}
		contents = append([]mcp.Content{summary}, contents...)
	}

	return &mcp.CallToolResult{
		Content: contents,
	}, nil, nil
}
