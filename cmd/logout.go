package cmd

import (
	"StarFileManager/internal/call"

	"github.com/spf13/cobra"
)

// logoutCmd 退出登录
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout",
	Long:  `Log out the currently logged in user`,
	Run: func(cmd *cobra.Command, args []string) {
		call.Logout(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
