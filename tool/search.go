package tool

import (
	"context"
	"fmt"
	"strings"

	"github.com/bamboo-services/bamboo-document-mcp/models"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DocumentSearchInput 文档搜索工具的输入参数。
type DocumentSearchInput struct {
	Sector       string `json:"sector" jsonschema:"required,板块标识"`              // Sector 板块标识
	Path         string `json:"path" jsonschema:"required,文档路径"`                // Path 文档路径
	Keyword      string `json:"keyword" jsonschema:"required,搜索关键词"`            // Keyword 搜索关键词
	ContextLines int    `json:"context_lines" jsonschema:"optional,上下文行数，默认 3"` // ContextLines 上下文行数
}

// DocumentSearchOutput 文档搜索工具的输出结果。
type DocumentSearchOutput struct {
	Sector       string               `json:"sector" jsonschema:"板块标识"`           // Sector 板块标识
	Path         string               `json:"path" jsonschema:"文档路径"`             // Path 文档路径
	Keyword      string               `json:"keyword" jsonschema:"搜索关键词"`         // Keyword 搜索关键词
	TotalMatches int                  `json:"total_matches" jsonschema:"总匹配数"`    // TotalMatches 总匹配数
	Matches      []models.SearchMatch `json:"matches" jsonschema:"匹配结果列表"`        // Matches 匹配结果列表
	DocumentURI  string               `json:"document_uri" jsonschema:"完整文档 URI"` // DocumentURI 完整文档 URI
}

// DocumentSearch 在指定文档内容中搜索关键词，返回匹配行及上下文段落。
func (t *Tool) DocumentSearch(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input DocumentSearchInput,
) (*mcp.CallToolResult, DocumentSearchOutput, error) {
	// 1. 设置默认上下文行数
	contextLines := input.ContextLines
	if contextLines <= 0 {
		contextLines = 3
	}

	// 2. 确保 path 以 / 开头
	path := input.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 3. 获取文档内容
	uri := fmt.Sprintf("https://doc.x-lf.com/docs/%s%s", input.Sector, path)
	resp, err := t.client.R().Get(uri)
	if err != nil {
		return nil, DocumentSearchOutput{}, err
	}
	if resp.StatusCode() != 200 {
		return nil, DocumentSearchOutput{}, fmt.Errorf("文档不存在: HTTP %d", resp.StatusCode())
	}

	// 4. 按行分割并搜索
	content := resp.String()
	lines := strings.Split(content, "\n")
	keywordLower := strings.ToLower(input.Keyword)
	matches := make([]models.SearchMatch, 0)

	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), keywordLower) {
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

			matches = append(matches, models.SearchMatch{
				LineNumber: i + 1,
				Line:       line,
				Context:    contextBuilder.String(),
			})
		}
	}

	return nil, DocumentSearchOutput{
		Sector:       input.Sector,
		Path:         input.Path,
		Keyword:      input.Keyword,
		TotalMatches: len(matches),
		Matches:      matches,
		DocumentURI:  uri,
	}, nil
}
