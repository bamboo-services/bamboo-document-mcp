package tool

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestDocumentDetail(t *testing.T) {
	tool := NewTool(resty.New())

	_, output, err := tool.DocumentDetail(context.Background(), nil, DocumentDetailInput{
		Sector: "bamboo-base-go",
		Path:   "/architecture",
	})
	if err != nil {
		t.Errorf("DocumentDetail failed: %v", err)
		return
	}

	if output.Title == "" {
		t.Error("Expected non-empty title")
	}
	if output.Content == "" {
		t.Error("Expected non-empty content")
	}
	if output.Sector != "bamboo-base-go" {
		t.Errorf("Expected sector 'bamboo-base-go', got '%s'", output.Sector)
	}

	t.Logf("Title: %s", output.Title)
	t.Logf("DocumentURI: %s", output.DocumentURI)
	t.Logf("Content: %s", output.Content)
}

func TestDocumentDetailWithInvalidPath(t *testing.T) {
	tool := NewTool(resty.New())

	_, _, err := tool.DocumentDetail(context.Background(), nil, DocumentDetailInput{
		Sector: "bamboo-base-go",
		Path:   "/non-existent-path-12345",
	})
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestDocumentDetailPathNormalization(t *testing.T) {
	tool := NewTool(resty.New())

	// 测试不带斜杠的路径
	_, output, err := tool.DocumentDetail(context.Background(), nil, DocumentDetailInput{
		Sector: "bamboo-base-go",
		Path:   "architecture",
	})
	if err != nil {
		t.Errorf("DocumentDetail failed: %v", err)
		return
	}

	if output.Title == "" {
		t.Error("Expected non-empty title")
	}
	t.Logf("Title: %s", output.Title)
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
