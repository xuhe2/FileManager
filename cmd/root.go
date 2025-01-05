package cmd

import (
	"StarFileManager/internal/call"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd 根命令,输出提示信息
var rootCmd = &cobra.Command{
	Use:   "FileManager",
	Short: "A Multi-User Simulated File Management System",
	Long: `一个多用户仿文件管理系统
仿照Linux文件管理系统实现,支持文件的增删改查和多用户权限管理
需要使用docker compose启动服务`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func ExecuteContext(ctx context.Context) {
	rootCmd.ExecuteContext(ctx)
}

func init() {
	// 身份验证中间件
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// 根命令和登录注册命令不检测
		if cmd != rootCmd && cmd != loginCmd && cmd != registerCmd {
			if call.GetUser(cmd.Context()) == "" {
				fmt.Fprintln(os.Stderr, "need login first")
				os.Exit(1)
			}
		}
	}
}
