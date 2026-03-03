// Bamboo Document MCP - Bamboo 竹简库文档服务 MCP 服务器
//
// 本服务提供从 doc.x-lf.com 获取 Bamboo 竹简库文档信息的能力，
// 支持文档目录列表、文档详情获取、文档内容搜索等功能。
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bamboo-services/bamboo-document-mcp/route"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// 版本信息，编译时通过 ldflags 注入
var version = "0.0.3"

// main 初始化并启动 MCP 服务器，注册所有工具路由后通过标准 I/O 传输运行。
func main() {
	// 命令行参数解析
	showHelp := flag.Bool("help", false, "显示帮助信息")
	showVersion := flag.Bool("version", false, "显示版本信息")
	flag.Parse()

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("bamboo-document v%s\n", version)
		os.Exit(0)
	}

	ctx := context.Background()
	server := mcp.NewServer(&mcp.Implementation{Name: "BambooDocument", Version: version}, nil)

	// 构建路由
	rTer := route.NewRoute(ctx, server)
	rTer.RouteBuild()

	// 启动 MCP 服务器
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Printf(`bamboo-document v%s - Bamboo 竹简库文档 MCP 服务

用法:
    bamboo-document [选项]

选项:
    -help       显示帮助信息
    -version    显示版本信息

说明:
    这是一个 MCP (Model Context Protocol) 服务器，用于从 doc.x-lf.com
    获取 Bamboo 竹简库的文档信息。

    安装后，Claude Code 会自动检测并加载此 MCP 服务器提供的工具。

可用工具:
    list        获取文档目录列表
    detail      获取文档详情
    search      搜索文档内容

更多信息:
    文档站点: https://doc.x-lf.com
    项目主页: https://github.com/bamboo-services/bamboo-document-mcp

`, version)
}
