package tool

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestDocumentList(t *testing.T) {
	tool := NewTool(resty.New())
	search := "ID"

	result, _, err := tool.DocumentList(context.Background(), nil, DocumentListInput{Sector: "bamboo-base-java", Search: &search})
	if err != nil {
		t.Errorf("DocumentList failed: %v", err)
		return
	}

	if result.IsError {
		t.Error("Expected successful result")
	}

	// 检查第二个 Content (文档列表)
	if len(result.Content) < 2 {
		t.Error("Expected at least 2 content items")
		return
	}

	textContent := result.Content[1].(*mcp.TextContent).Text
	if !strings.Contains(textContent, "文档列表") {
		t.Error("Expected result to contain '文档列表'")
	}

	t.Logf("Base URL: %s", result.Content[0].(*mcp.TextContent).Text)
	t.Logf("List: %s", textContent)
}

func TestDebug(t *testing.T) {
	input := "- [筱工具(Golang)](/docs/bamboo-base-go): 筱锋的 Go 语言基础组件库文档"

	// 定义正则
	re := regexp.MustCompile(`- \[(.*?)]\((.*?)\):\s*(.*)`)

	// 查找匹配项
	matches := re.FindStringSubmatch(input)

	if len(matches) > 3 {
		fmt.Printf("1. 标题: %s\n", matches[1]) // 核心说明
		fmt.Printf("2. 路径: %s\n", matches[2]) // /docs/bamboo-base-go/architecture
		fmt.Printf("3. 说明: %s\n", matches[3]) // Bamboo Base 整体架构设计与模块说明
	}
}
