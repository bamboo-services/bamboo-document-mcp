package models

// SearchResponse API 搜索响应，包含搜索结果列表。
type SearchResponse []SearchResult

// SearchResult 单个搜索结果。
type SearchResult struct {
	ID                    string             `json:"id"`                    // ID 结果唯一标识
	Type                  string             `json:"type"`                  // Type 结果类型 (page/heading/text)
	Content               string             `json:"content"`               // Content 内容文本
	Breadcrumbs           []string           `json:"breadcrumbs,omitempty"` // Breadcrumbs 面包屑导航
	ContentWithHighlights []ContentHighlight `json:"contentWithHighlights"` // ContentWithHighlights 高亮内容片段
	URL                   string             `json:"url"`                   // URL 文档地址
}

// ContentHighlight 高亮内容片段。
type ContentHighlight struct {
	Type    string          `json:"type"`             // Type 片段类型
	Content string          `json:"content"`          // Content 片段内容
	Styles  map[string]bool `json:"styles,omitempty"` // Styles 样式标记 (如 highlight)
}
