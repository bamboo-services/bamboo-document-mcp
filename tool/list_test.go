package tool

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestDocumentList(t *testing.T) {
	tool := NewTool(resty.New())

	_, output, err := tool.DocumentList(context.Background(), nil, DocumentListInput{Sector: "bamboo-base-java", Search: new("ID")})
	if err != nil {
		t.Errorf("DocumentList failed: %v", err)
		return
	}

	for _, list := range output.List {
		t.Log(list)
	}
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
