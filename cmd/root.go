package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"os"
)

// tuiFlag 是否启用gui标志
var tuiFlag bool

// rootCmd 根命令,输出提示信息和启动tui界面
var rootCmd = &cobra.Command{
	Use:   "StarFileManager",
	Short: "一个多用户仿文件管理系统",
	Long: `一个多用户仿文件管理系统
仿照Linux文件管理系统实现,支持文件的增删改查和多用户权限管理
项目依赖MongoDB,需要保证27017端口正常运行MongoDB服务`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flags().BoolVarP(&tuiFlag, "tui", "t", false, "启用TUI")

		if tuiFlag {
			// TODO 实现tui
		} else {
			cmd.Help()
		}
	}}

func ExecuteContext(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
