package tool

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestDocumentSearch(t *testing.T) {
	tool := NewTool(resty.New())

	_, output, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:  "bamboo-base-go",
		Path:    "/architecture",
		Keyword: "架构",
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if output.TotalMatches == 0 {
		t.Error("Expected at least one match")
	}
	if len(output.Matches) != output.TotalMatches {
		t.Errorf("Matches count mismatch: %d vs %d", len(output.Matches), output.TotalMatches)
	}

	t.Logf("TotalMatches: %d", output.TotalMatches)
	for i, match := range output.Matches {
		t.Logf("Match %d - Line %d: %s", i+1, match.LineNumber, match.Line)
	}
}

func TestDocumentSearchWithContextLines(t *testing.T) {
	tool := NewTool(resty.New())

	contextLines := 5
	_, output, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:       "bamboo-base-go",
		Path:         "/architecture",
		Keyword:      "架构",
		ContextLines: &contextLines,
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	if output.TotalMatches == 0 {
		t.Error("Expected at least one match")
	}

	t.Logf("TotalMatches with contextLines=5: %d", output.TotalMatches)
}

func TestDocumentSearchNilContextLines(t *testing.T) {
	tool := NewTool(resty.New())

	_, output, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:       "bamboo-base-go",
		Path:         "/architecture",
		Keyword:      "架构",
		ContextLines: nil, // 测试默认值
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	t.Logf("TotalMatches with nil contextLines: %d", output.TotalMatches)
}

func TestDocumentSearchCaseInsensitive(t *testing.T) {
	tool := NewTool(resty.New())

	// 测试大小写不敏感
	_, output, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:  "bamboo-base-go",
		Path:    "/architecture",
		Keyword: "ARCHITECTURE", // 大写
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	// 如果文档中有 architecture 相关内容，应该能匹配到
	t.Logf("Case insensitive search - TotalMatches: %d", output.TotalMatches)
}

func TestDocumentSearchInvalidPath(t *testing.T) {
	tool := NewTool(resty.New())

	_, _, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:  "bamboo-base-go",
		Path:    "/non-existent-path-12345",
		Keyword: "test",
	})
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestDocumentSearchPathNormalization(t *testing.T) {
	tool := NewTool(resty.New())

	// 测试不带斜杠的路径
	_, output, err := tool.DocumentSearch(context.Background(), nil, DocumentSearchInput{
		Sector:  "bamboo-base-go",
		Path:    "architecture",
		Keyword: "架构",
	})
	if err != nil {
		t.Errorf("DocumentSearch failed: %v", err)
		return
	}

	t.Logf("Path normalization - TotalMatches: %d", output.TotalMatches)
}
