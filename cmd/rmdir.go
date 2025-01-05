package cmd

import (
	"StarFileManager/internal/call"
	"errors"

	"github.com/spf13/cobra"
)

// rmdirCmd 删除空目录
var rmdirCmd = &cobra.Command{
	Use:   "rmdir",
	Short: "remove empty directory",
	Long:  `Delete the specified directory, which must be an empty directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing directory name")
		}
		return call.DeleteFile(cmd.Context(), args[0], false, true, true)
	},
}

func init() {
	rootCmd.AddCommand(rmdirCmd)
}
