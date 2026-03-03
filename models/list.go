// Package models 定义了 MCP 工具使用的数据模型结构体。
package models

// DocumentList 文档目录项，表示单个文档的基本信息。
type DocumentList struct {
	Name        string `json:"name"`        // Name 文档名称
	Sector      string `json:"sector"`      // Sector 所属板块标识
	Path        string `json:"path"`        // Path 文档路径
	Description string `json:"description"` // Description 文档描述
}
