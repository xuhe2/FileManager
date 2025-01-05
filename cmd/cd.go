package cmd

import (
	"StarFileManager/internal/call"

	"github.com/spf13/cobra"
)

// cdCmd 进入指定目录
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "change directory",
	Long:  `change directory, can use abs and relative path`,
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
