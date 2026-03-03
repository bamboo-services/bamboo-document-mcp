package tool

import (
	"context"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestDocumentSearch(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:  "bamboo-base-go",
		Path:    "/architecture",
		Keyword: "架构",
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	// 检查有多个 Content（汇总 + 基础信息 + 匹配结果）
	if len(result.Content) < 5 {
		t.Errorf("Expected at least 5 content items, got %d", len(result.Content))
	}

	// 第一个应该是汇总信息
	summary := result.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(summary, "共找到") {
		t.Error("Expected first content to contain '共找到'")
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

func TestDocumentSearchWithContextLines(t *testing.T) {
	tool := NewTool(resty.New())

	contextLines := 5
	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:       "bamboo-base-go",
		Path:         "/architecture",
		Keyword:      "架构",
		ContextLines: &contextLines,
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	summary := result.Content[0].(*mcp.TextContent).Text
	t.Logf("Summary: %s", summary)
}

func TestDocumentSearchNilContextLines(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:       "bamboo-base-go",
		Path:         "/architecture",
		Keyword:      "架构",
		ContextLines: nil, // 测试默认值
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	summary := result.Content[0].(*mcp.TextContent).Text
	t.Logf("Summary with nil contextLines: %s", summary)
}

func TestDocumentSearchCaseInsensitive(t *testing.T) {
	tool := NewTool(resty.New())

	// 测试大小写不敏感
	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:  "bamboo-base-go",
		Path:    "/architecture",
		Keyword: "ARCHITECTURE", // 大写
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	summary := result.Content[0].(*mcp.TextContent).Text
	t.Logf("Case insensitive search: %s", summary)
}

func TestDocumentSearchInvalidPath(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:  "bamboo-base-go",
		Path:    "/non-existent-path-12345",
		Keyword: "test",
	})
	if err == nil {
		t.Error("Expected error for non-existent path")
	}

	if result != nil && !result.IsError {
		t.Error("Expected error result")
	}
}

func TestDocumentSearchPathNormalization(t *testing.T) {
	tool := NewTool(resty.New())

	// 测试不带斜杠的路径
	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:  "bamboo-base-go",
		Path:    "architecture",
		Keyword: "架构",
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	summary := result.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(summary, "共找到") {
		t.Error("Expected summary to contain '共找到'")
	}
}

func TestDocumentSearchNoMatch(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:  "bamboo-base-go",
		Path:    "/architecture",
		Keyword: "不存在的关键词xyz123",
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	// 最后一个应该是"未找到匹配内容"
	lastContent := result.Content[len(result.Content)-1].(*mcp.TextContent).Text
	if !strings.Contains(lastContent, "未找到匹配内容") {
		t.Errorf("Expected last content to contain '未找到匹配内容', got: %s", lastContent)
	}
}
