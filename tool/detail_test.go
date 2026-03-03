package tool

import (
	"context"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestDocumentDetail(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.DocumentDetail(context.Background(), nil, DocumentDetailInput{
		Sector: "bamboo-base-go",
		Path:   "/architecture",
	})
	if err != nil {
		t.Errorf("DocumentDetail failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	// 检查多个 Content
	if len(result.Content) < 5 {
		t.Errorf("Expected at least 5 content items, got %d", len(result.Content))
	}

	// 验证各个字段
	textContent := result.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(textContent, "板块:") {
		t.Error("Expected first content to contain '板块:'")
	}

	t.Logf("Contents:")
	for _, c := range result.Content {
		text := c.(*mcp.TextContent).Text
		if len(text) > 100 {
			t.Logf("  %s...", text[:100])
		} else {
			t.Logf("  %s", text)
		}
	}
}

func TestDocumentDetailWithInvalidPath(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.DocumentDetail(context.Background(), nil, DocumentDetailInput{
		Sector: "bamboo-base-go",
		Path:   "/non-existent-path-12345",
	})
	if err == nil {
		t.Error("Expected error for non-existent path")
	}

	if result != nil && !result.IsError {
		t.Error("Expected error result")
	}
}

func TestDocumentDetailPathNormalization(t *testing.T) {
	tool := NewTool(resty.New())

	// 测试不带斜杠的路径
	result, _, err := tool.DocumentDetail(context.Background(), nil, DocumentDetailInput{
		Sector: "bamboo-base-go",
		Path:   "architecture",
	})
	if err != nil {
		t.Errorf("DocumentDetail failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	// 验证文档内容包含预期关键词
	contentText := result.Content[4].(*mcp.TextContent).Text
	if !strings.Contains(contentText, "架构") {
		t.Error("Expected content to contain '架构'")
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "标准标题",
			content:  "# 这是标题\n\n内容",
			expected: "这是标题",
		},
		{
			name:     "带空格的标题",
			content:  "#   空格标题   \n内容",
			expected: "空格标题",
		},
		{
			name:     "无标题",
			content:  "没有标题的内容",
			expected: "",
		},
		{
			name:     "多行内容中的标题",
			content:  "一些内容\n# 中间标题\n更多内容",
			expected: "中间标题",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTitle(tt.content)
			if result != tt.expected {
				t.Errorf("extractTitle() = '%s', expected '%s'", result, tt.expected)
			}
		})
	}
}
