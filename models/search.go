package models

// SearchMatch 搜索匹配结果，表示文档中匹配关键词的单个结果。
type SearchMatch struct {
	LineNumber int    `json:"line_number" jsonschema:"匹配行号（从 1 开始）"` // LineNumber 匹配行号（从 1 开始）
	Line       string `json:"line" jsonschema:"匹配的行内容"`              // Line 匹配的行内容
	Context    string `json:"context" jsonschema:"上下文内容（包含前后 N 行）"`  // Context 上下文内容（包含前后 N 行）
}
