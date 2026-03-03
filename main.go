// Bamboo Document MCP - Bamboo 竹简库文档服务 MCP 服务器
//
// 本服务提供从 doc.x-lf.com 获取 Bamboo 竹简库文档信息的能力，
// 支持文档目录列表、文档详情获取、文档内容搜索等功能。
package main

import (
	"github.com/bamboo-services/bamboo-document-mcp/command"
)

// 版本信息，编译时通过 ldflags 注入
var version = "0.0.3"

// main 程序入口，初始化命令行并执行。
func main() {
	// 设置版本号
	command.SetRootVersion(version)
	command.SetCmdVersion(version)
	command.SetVersion(version)

	// 执行命令
	if err := command.Execute(); err != nil {
		panic(err)
	}
}
