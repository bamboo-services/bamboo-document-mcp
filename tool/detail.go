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

// DocumentDetailOutput 文档详情工具的输出结果。
type DocumentDetailOutput struct {
	Sector      string `json:"sector" jsonschema:"板块标识"`              // Sector 板块标识
	Path        string `json:"path" jsonschema:"文档路径"`                // Path 文档路径
	Title       string `json:"title" jsonschema:"文档标题"`               // Title 文档标题（从 Markdown 提取）
	Content     string `json:"content" jsonschema:"文档 Markdown 原始内容"` // Content 文档 Markdown 原始内容
	DocumentURI string `json:"document_uri" jsonschema:"完整文档 URI"`    // DocumentURI 完整文档 URI
}

// DocumentDetail 获取指定文档的完整 Markdown 内容。
func (t *Tool) DocumentDetail(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input DocumentDetailInput,
) (*mcp.CallToolResult, DocumentDetailOutput, error) {
	// 1. 确保 path 以 / 开头
	path := input.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 2. 构建请求 URL
	uri := fmt.Sprintf("https://doc.x-lf.com/docs/%s%s", input.Sector, path)

	// 3. 发送 GET 请求
	resp, err := t.client.R().Get(uri)
	if err != nil {
		return nil, DocumentDetailOutput{}, err
	}

	// 4. 检查响应状态码
	if resp.StatusCode() != 200 {
		return nil, DocumentDetailOutput{}, fmt.Errorf("文档不存在或获取失败: HTTP %d", resp.StatusCode())
	}

	// 5. 提取标题
	content := resp.String()
	title := extractTitle(content)

	// 6. 返回结果
	return nil, DocumentDetailOutput{
		Sector:      input.Sector,
		Path:        input.Path,
		Title:       title,
		Content:     content,
		DocumentURI: uri,
	}, nil
}

// extractTitle 从 Markdown 内容中提取第一个标题
func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}
