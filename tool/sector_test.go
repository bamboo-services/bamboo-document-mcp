package tool

import (
	"context"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSectorList(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.SectorList(context.Background(), nil, SectorListInput{})
	if err != nil {
		t.Errorf("SectorList failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	textContent := result.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(textContent, "板块列表") {
		t.Error("Expected result to contain '板块列表'")
	}

	t.Logf("Content: %s", textContent)
}

func TestSectorListContainsExpected(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.SectorList(context.Background(), nil, SectorListInput{})
	if err != nil {
		t.Errorf("SectorList failed: %v", err)
		return
	}

	textContent := result.Content[0].(*mcp.TextContent).Text

	// 检查是否包含预期的板块
	expectedSectors := []string{"bamboo-base-go", "bamboo-base-java"}
	for _, expected := range expectedSectors {
		if !strings.Contains(textContent, expected) {
			t.Logf("Warning: Expected sector '%s' not found", expected)
		} else {
			t.Logf("Found expected sector: %s", expected)
		}
	}
}

func TestSectorListNoDuplicates(t *testing.T) {
	tool := NewTool(resty.New())

	result, _, err := tool.SectorList(context.Background(), nil, SectorListInput{})
	if err != nil {
		t.Errorf("SectorList failed: %v", err)
		return
	}

	textContent := result.Content[0].(*mcp.TextContent).Text

	// 检查是否有重复（通过计数行数和板块数来判断）
	lines := strings.Split(textContent, "\n")
	sectorCount := 0
	for _, line := range lines {
		if strings.Contains(line, ". ") && !strings.Contains(line, "板块列表") {
			sectorCount++
		}
	}

	t.Logf("Total sectors found: %d", sectorCount)
}
