package models

// SectorInfo 板块信息，表示一个可用的文档板块。
type SectorInfo struct {
	Sector string `json:"sector" jsonschema:"板块标识"` // Sector 板块标识
}
