package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bamboo-services/bamboo-document-mcp/models"
)

// DocumentSearchInput 文档搜索工具的输入参数。
type DocumentSearchInput struct {
	Query  string  `json:"query" jsonschema:"required,搜索关键词（只支持英文）"` // Query 搜索关键词（必需）
	Sector *string `json:"sector" jsonschema:"optional,板块标识筛选"`      // Sector 板块筛选（可选）
	Path   *string `json:"path" jsonschema:"optional,路径筛选"`          // Path 路径筛选（可选）
}

// DocumentSearch 调用 API 在文档中搜索关键词，返回匹配结果。
func (t *Tool) DocumentSearch(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input DocumentSearchInput,
) (*mcp.CallToolResult, any, error) {
	// 1. 调用搜索 API
	resp, err := t.client.R().
		SetQueryParam("query", input.Query).
		Get("https://doc.x-lf.com/api/search")
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "搜索 API 调用失败"},
			},
		}, nil, err
	}
	if resp.StatusCode() != 200 {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("搜索请求失败: HTTP %d", resp.StatusCode())},
			},
		}, nil, fmt.Errorf("搜索请求失败: HTTP %d", resp.StatusCode())
	}

	// 2. 解析 JSON 响应
	var searchResults models.SearchResponse
	if err := json.Unmarshal(resp.Body(), &searchResults); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "搜索结果解析失败"},
			},
		}, nil, err
	}

	// 3. 根据 Sector 和 Path 筛选结果
	filteredResults := t.filterSearchResults(searchResults, input.Sector, input.Path)

	// 4. 构建输出内容
	contents := make([]mcp.Content, 0)
	contents = append(contents,
		&mcp.TextContent{Text: fmt.Sprintf("关键词: %s", input.Query)},
	)

	if input.Sector != nil {
		contents = append(contents, &mcp.TextContent{Text: fmt.Sprintf("板块筛选: %s", *input.Sector)})
	}
	if input.Path != nil {
		contents = append(contents, &mcp.TextContent{Text: fmt.Sprintf("路径筛选: %s", *input.Path)})
	}

	// 5. 添加汇总信息
	if len(filteredResults) == 0 {
		contents = append(contents, &mcp.TextContent{Text: "未找到匹配内容"})
	} else {
		contents = append(contents, &mcp.TextContent{Text: fmt.Sprintf("共找到 %d 个匹配", len(filteredResults))})
	}

	// 6. 格式化每个搜索结果
	for i, result := range filteredResults {
		// 提取 Sector 和 Path
		sector := t.extractSector(result.URL)
		detailPath := t.extractDetailPath(result.URL)

		// 获取概览内容
		overview := t.buildHighlightContent(result.ContentWithHighlights)

		// 格式化输出: 1. [Sector](Path)
		contents = append(contents, &mcp.TextContent{Text: fmt.Sprintf("%d. [%s](%s)", i+1, sector, detailPath)})
		contents = append(contents, &mcp.TextContent{Text: fmt.Sprintf(" ⎿ %s", overview)})
	}

	return &mcp.CallToolResult{
		Content: contents,
	}, nil, nil
}

// filterSearchResults 根据 Sector 和 Path 筛选搜索结果。
func (t *Tool) filterSearchResults(results models.SearchResponse, sector, path *string) models.SearchResponse {
	if sector == nil && path == nil {
		return results
	}

	filtered := make(models.SearchResponse, 0)
	for _, result := range results {
		// 检查 Sector 筛选
		if sector != nil {
			if !strings.Contains(result.ID, *sector) {
				continue
			}
		}

		// 检查 Path 筛选
		if path != nil {
			pathVal := *path
			if !strings.HasPrefix(pathVal, "/") {
				pathVal = "/" + pathVal
			}
			if !strings.Contains(result.URL, pathVal) {
				continue
			}
		}

		filtered = append(filtered, result)
	}

	return filtered
}

// buildHighlightContent 从高亮片段构建内容摘要。
func (t *Tool) buildHighlightContent(highlights []models.ContentHighlight) string {
	var builder strings.Builder
	for _, h := range highlights {
		builder.WriteString(h.Content)
	}
	return builder.String()
}

// extractDetailPath 从 URL 中提取可用于 detail 接口的 path。
// URL 格式: /docs/{sector}{path} -> 返回 {path}
func (t *Tool) extractDetailPath(url string) string {
	// 移除 /docs/ 前缀
	if !strings.HasPrefix(url, "/docs/") {
		return url
	}
	// 移除 "/docs/" (6 个字符)
	remaining := strings.TrimPrefix(url, "/docs/")
	// 找到第一个 "/" 来确定 sector 边界
	idx := strings.Index(remaining, "/")
	if idx == -1 {
		return "/" + remaining // 没有 path 部分，返回根路径
	}
	return remaining[idx:] // 返回 sector 之后的部分
}

// extractSector 从 URL 中提取板块标识。
// URL 格式: /docs/{sector}{path} -> 返回 {sector}
func (t *Tool) extractSector(url string) string {
	// 移除 /docs/ 前缀
	if !strings.HasPrefix(url, "/docs/") {
		return ""
	}
	// 移除 "/docs/" (6 个字符)
	remaining := strings.TrimPrefix(url, "/docs/")
	// 找到第一个 "/" 来确定 sector 边界
	idx := strings.Index(remaining, "/")
	if idx == -1 {
		// 没有 path 部分，整个就是 sector，但需要移除锚点
		return t.removeAnchor(remaining)
	}
	return t.removeAnchor(remaining[:idx]) // 返回 sector 部分
}

// removeAnchor 移除 URL 中的锚点部分。
func (t *Tool) removeAnchor(s string) string {
	if idx := strings.Index(s, "#"); idx != -1 {
		return s[:idx]
	}
	return s
}
