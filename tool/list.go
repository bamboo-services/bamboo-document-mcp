package tool

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bamboo-services/bamboo-document-mcp/models"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DocumentListInput 文档列表工具的输入参数。
type DocumentListInput struct {
	Sector string  `json:"sector" jsonschema:"required,板块内容"`  // Sector 板块标识，如 bamboo-base-go
	Search *string `json:"search" jsonschema:"optional,搜索关键词"` // Search 可选的搜索关键词
}

// DocumentListOutput 文档列表工具的输出结果。
type DocumentListOutput struct {
	BaseURI string                `json:"base_uri" jsonschema:"基础地址"` // BaseURI 文档基础地址
	List    []models.DocumentList `json:"list" jsonschema:"文档目录"`     // List 文档目录列表
}

// DocumentList 获取指定板块的文档目录列表，支持按关键词筛选。
func (t *Tool) DocumentList(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input DocumentListInput,
) (*mcp.CallToolResult, DocumentListOutput, error) {
	resp, err := t.client.R().Get("https://doc.x-lf.com/llms.txt")
	if err != nil {
		return nil, DocumentListOutput{}, err
	}

	// 格式化文本
	newStr := strings.Replace(resp.String(), "# Documentation\n\n", "", 1)
	lines := strings.Split(newStr, "\n")
	re := regexp.MustCompile(`- \[(.*?)]\((.*?)\):\s*(.*)`)

	list := make([]models.DocumentList, 0)
	for _, line := range lines {
		// 解析
		matches := re.FindStringSubmatch(line)
		if matches == nil || len(matches) != 4 {
			return nil, DocumentListOutput{}, fmt.Errorf("格式无法解析")
		}

		// 筛选
		if !strings.HasPrefix(matches[2], "/docs/"+input.Sector) {
			continue
		}
		if input.Search != nil && !strings.Contains(strings.ToLower(matches[1]), strings.ToLower(*input.Search)) && !strings.Contains(strings.ToLower(matches[3]), strings.ToLower(*input.Search)) {
			continue
		}

		// 组合器
		list = append(list, models.DocumentList{
			Name:        matches[1],
			Sector:      input.Sector,
			Path:        strings.Replace(matches[2], "/docs/"+input.Sector, "", 1),
			Description: matches[3],
		})
	}

	return nil, DocumentListOutput{BaseURI: "https://doc.x-lf.com/", List: list}, nil
}
