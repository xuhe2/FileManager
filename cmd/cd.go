package cmd

import (
	"StarFileManager/internal/call"
	"github.com/spf13/cobra"
)

// cdCmd 进入指定目录
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "进入目录",
	Long:  `进入参数指定的目录,支持绝对和相对路径`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return nil
		}
		return call.ChangePath(cmd.Context(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(cdCmd)
}
