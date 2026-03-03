// Package command 提供命令行接口的实现。
package command

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// 版本信息，编译时通过 ldflags 注入
var cmdVersion = "0.0.3"

// versionCmd 版本命令，显示版本信息。
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 bamboo-document 的版本信息，包括版本号和运行环境。`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("bamboo-document v%s\n", cmdVersion)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

// GetVersionCmd 返回版本命令实例。
func GetVersionCmd() *cobra.Command {
	return versionCmd
}

// SetCmdVersion 设置命令版本号。
func SetCmdVersion(v string) {
	cmdVersion = v
}
