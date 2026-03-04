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
		Query: "api",
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	// 检查有多个 Content（关键词 + 汇总 + 匹配结果）
	if len(result.Content) < 3 {
		t.Errorf("Expected at least 3 content items, got %d", len(result.Content))
	}

	// 第一个应该是关键词信息
	firstContent := result.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(firstContent, "关键词") {
		t.Errorf("Expected first content to contain '关键词', got: %s", firstContent)
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

func TestDocumentSearchWithSector(t *testing.T) {
	tool := NewTool(resty.New())

	sector := "bamboo-base-go"
	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Query:  "api",
		Sector: &sector,
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	// 检查是否包含板块筛选信息
	hasSectorFilter := false
	for _, c := range result.Content {
		text := c.(*mcp.TextContent).Text
		if strings.Contains(text, "板块筛选") {
			hasSectorFilter = true
			break
		}
	}

	if !hasSectorFilter {
		t.Error("Expected sector filter info in results")
	}

	t.Logf("Sector filter test passed, results count: %d", len(result.Content))
}

func TestDocumentSearchWithPath(t *testing.T) {
	tool := NewTool(resty.New())

	path := "result"
	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Query: "response",
		Path:  &path,
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	t.Logf("Path filter test passed, results count: %d", len(result.Content))
}

func TestDocumentSearchWithBothFilters(t *testing.T) {
	tool := NewTool(resty.New())

	sector := "bamboo-base-go"
	path := "result"
	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Query:  "response",
		Sector: &sector,
		Path:   &path,
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	t.Logf("Both filters test passed, results count: %d", len(result.Content))
}

func TestDocumentSearchNoMatch(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Query: "nonexistentkeywordxyz12345",
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result (even with no matches)")
	}

	// 应该包含"未找到匹配内容"
	hasNoMatch := false
	for _, c := range result.Content {
		text := c.(*mcp.TextContent).Text
		if strings.Contains(text, "未找到匹配内容") {
			hasNoMatch = true
			break
		}
	}

	if !hasNoMatch {
		t.Error("Expected '未找到匹配内容' when no matches found")
	}

	t.Logf("No match test passed")
}

func TestDocumentSearchEmptyQuery(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Query: "",
	})

	// 空查询可能返回错误或空结果
	if err != nil {
		t.Logf("Empty query returned error (expected): %v", err)
		return
	}

	if result.IsError {
		t.Logf("Empty query returned error result (expected)")
		return
	}

	t.Logf("Empty query returned %d results", len(result.Content))
}

func TestDocumentSearchCaseInsensitive(t *testing.T) {
	tool := NewTool(resty.New())

	// 测试大写关键词
	result, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Query: "API",
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	t.Logf("Case insensitive search (uppercase) passed, results: %d", len(result.Content))
}

func TestExtractDetailPath(t *testing.T) {
	tool := NewTool(resty.New())

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "标准 URL",
			input:    "/docs/bamboo-base-go/core/init/config",
			expected: "/core/init/config",
		},
		{
			name:     "带锚点的 URL",
			input:    "/docs/bamboo-base-go/quick-start#api-结构定义",
			expected: "/quick-start#api-结构定义",
		},
		{
			name:     "仅有 sector 无 path",
			input:    "/docs/bamboo-base-go",
			expected: "/bamboo-base-go",
		},
		{
			name:     "不含 /docs/ 前缀",
			input:    "/other/path",
			expected: "/other/path",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "Java 板块文档",
			input:    "/docs/bamboo-base-java/mvc/filter/context",
			expected: "/mvc/filter/context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tool.extractDetailPath(tt.input)
			if result != tt.expected {
				t.Errorf("extractDetailPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractSector(t *testing.T) {
	tool := NewTool(resty.New())

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "标准 URL",
			input:    "/docs/bamboo-base-go/core/init/config",
			expected: "bamboo-base-go",
		},
		{
			name:     "带锚点的 URL",
			input:    "/docs/bamboo-base-go/quick-start#api-结构定义",
			expected: "bamboo-base-go",
		},
		{
			name:     "仅有 sector 无 path",
			input:    "/docs/bamboo-base-go",
			expected: "bamboo-base-go",
		},
		{
			name:     "仅有 sector 带锚点",
			input:    "/docs/guide#关于竹简库",
			expected: "guide",
		},
		{
			name:     "不含 /docs/ 前缀",
			input:    "/other/path",
			expected: "",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "Java 板块文档",
			input:    "/docs/bamboo-base-java/mvc/filter/context",
			expected: "bamboo-base-java",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tool.extractSector(tt.input)
			if result != tt.expected {
				t.Errorf("extractSector(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
