package cmd

import (
	"StarFileManager/internal/call"
	"github.com/spf13/cobra"
)

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "进入目录",
	Long:  `进入参数指定的目录,支持绝对和相对路径`,
	Run: func(cmd *cobra.Command, args []string) {
		call.ChangePath(cmd.Context(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(cdCmd)
}
