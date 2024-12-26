package cmd

import (
	"StarFileManager/internal/call"
	"github.com/spf13/cobra"
)

var rmdirCmd = &cobra.Command{
	Use:   "rmdir",
	Short: "删除空目录",
	Long:  `删除指定的目录,必须是空目录`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return call.DeleteFile(cmd.Context(), args[0], false, true, true)
	},
}

func init() {
	rootCmd.AddCommand(rmdirCmd)
}
