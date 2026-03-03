package tool

import (
	"context"
	"regexp"
	"strings"

	"github.com/bamboo-services/bamboo-document-mcp/models"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SectorListInput 板块列表工具的输入参数（无参数）。
type SectorListInput struct{}

// SectorListOutput 板块列表工具的输出结果。
type SectorListOutput struct {
	Sectors []models.SectorInfo `json:"sectors" jsonschema:"可用板块列表"` // Sectors 可用板块列表
}

// SectorList 获取所有可用的文档板块列表。
func (t *Tool) SectorList(
	_ context.Context,
	_ *mcp.CallToolRequest,
	_ SectorListInput,
) (*mcp.CallToolResult, SectorListOutput, error) {
	resp, err := t.client.R().Get("https://doc.x-lf.com/llms.txt")
	if err != nil {
		return nil, SectorListOutput{}, err
	}

	// 解析文档内容提取板块
	content := resp.String()
	lines := strings.Split(content, "\n")
	re := regexp.MustCompile(`- \[.*?\]\(/docs/([a-zA-Z0-9_-]+)`)

	sectorMap := make(map[string]bool)
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches != nil && len(matches) >= 2 {
			sectorMap[matches[1]] = true
		}
	}

	// 转换为列表
	sectors := make([]models.SectorInfo, 0, len(sectorMap))
	for sector := range sectorMap {
		sectors = append(sectors, models.SectorInfo{
			Sector: sector,
		})
	}

	return nil, SectorListOutput{Sectors: sectors}, nil
}
