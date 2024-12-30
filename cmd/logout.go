package cmd

import (
	"StarFileManager/internal/call"
	"github.com/spf13/cobra"
)

// logoutCmd 退出登录
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "退出登录",
	Long:  `退出当前登录的用户`,
	Run: func(cmd *cobra.Command, args []string) {
		call.Logout(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
