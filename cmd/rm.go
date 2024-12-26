package cmd

import (
	"StarFileManager/internal/call"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "删除文件",
	Long:  `删除文件,如果传入--recursive标志则可递归删除文件夹`,
	RunE: func(cmd *cobra.Command, args []string) error {
		deleteDir, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}
		return call.DeleteFile(cmd.Context(), args[0], true, deleteDir, false)
	},
}

func init() {
	rmCmd.Flags().BoolP("recursive", "r", false, "用户名,如果为空则为root")
	rootCmd.AddCommand(rmCmd)
}
