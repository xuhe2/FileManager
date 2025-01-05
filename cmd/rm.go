package cmd

import (
	"StarFileManager/internal/call"
	"errors"

	"github.com/spf13/cobra"
)

// rmCmd 删除文件
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove file",
	Long:  `Delete files. If you pass the '-r' flag, you can recursively delete folders.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		deleteDir, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}
		if len(args) < 1 {
			return errors.New("missing file name")
		}
		return call.DeleteFile(cmd.Context(), args[0], true, deleteDir, false)
	},
}

func init() {
	rmCmd.Flags().BoolP("recursive", "r", false, "delete folder recursively")
	rootCmd.AddCommand(rmCmd)
}
