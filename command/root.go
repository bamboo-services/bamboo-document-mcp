// Package command 提供命令行接口的实现。
package command

import (
	"github.com/spf13/cobra"
)

// 版本信息，编译时通过 ldflags 注入
var rootVersion = "0.0.3"

// rootCmd 根命令，作为所有子命令的入口。
var rootCmd = &cobra.Command{
	Use:   "bamboo-document",
	Short: "Bamboo 竹简库文档 MCP 服务",
	Long: `Bamboo 竹简库文档 MCP 服务

本服务提供从 doc.x-lf.com 获取 Bamboo 竹简库文档信息的能力，
支持文档目录列表、文档详情获取、文档内容搜索等功能。

可用工具:
    list        获取文档目录列表
    detail      获取文档详情
    search      搜索文档内容
    sector      获取可用板块列表

更多信息:
    文档站点: https://doc.x-lf.com
    项目主页: https://github.com/bamboo-services/bamboo-document-mcp`,
}

// Execute 执行根命令。
func Execute() error {
	return rootCmd.Execute()
}

// SetRootVersion 设置根命令版本号。
func SetRootVersion(v string) {
	rootVersion = v
}

// GetRootCmd 返回根命令实例。
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(GetMCPCmd())
	rootCmd.AddCommand(GetVersionCmd())
}
