package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"github.com/spf13/cobra"
)

// rmdirCmd 删除空目录
var rmdirCmd = &cobra.Command{
	Use:   "rmdir",
	Short: "删除空目录",
	Long:  `删除指定的目录,必须是空目录`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("缺少操作数")
		}
		return call.DeleteFile(cmd.Context(), args[0], false, true, true)
	},
}

func init() {
	rootCmd.AddCommand(rmdirCmd)
}
