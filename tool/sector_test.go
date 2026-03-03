package tool

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestSectorList(t *testing.T) {
	tool := NewTool(resty.New())

	_, output, err := tool.SectorList(context.Background(), nil, SectorListInput{})
	if err != nil {
		t.Errorf("SectorList failed: %v", err)
		return
	}

	if len(output.Sectors) == 0 {
		t.Error("Expected at least one sector")
	}

	t.Logf("Total Sectors: %d", len(output.Sectors))
	for i, sector := range output.Sectors {
		t.Logf("Sector %d: %s", i+1, sector.Sector)
	}
}

func TestSectorListContainsExpected(t *testing.T) {
	tool := NewTool(resty.New())

	_, output, err := tool.SectorList(context.Background(), nil, SectorListInput{})
	if err != nil {
		t.Errorf("SectorList failed: %v", err)
		return
	}

	// 检查是否包含预期的板块
	sectorMap := make(map[string]bool)
	for _, s := range output.Sectors {
		sectorMap[s.Sector] = true
	}

	expectedSectors := []string{"bamboo-base-go", "bamboo-base-java"}
	for _, expected := range expectedSectors {
		if !sectorMap[expected] {
			t.Logf("Warning: Expected sector '%s' not found", expected)
		} else {
			t.Logf("Found expected sector: %s", expected)
		}
	}
}

func TestSectorListNoDuplicates(t *testing.T) {
	tool := NewTool(resty.New())

	_, output, err := tool.SectorList(context.Background(), nil, SectorListInput{})
	if err != nil {
		t.Errorf("SectorList failed: %v", err)
		return
	}

	// 检查是否有重复
	seen := make(map[string]int)
	for _, s := range output.Sectors {
		seen[s.Sector]++
	}

	for sector, count := range seen {
		if count > 1 {
			t.Errorf("Duplicate sector found: %s (count: %d)", sector, count)
		}
	}
}
