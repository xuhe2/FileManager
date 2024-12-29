package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"github.com/spf13/cobra"
)

// rmCmd 删除文件
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "删除文件",
	Long:  `删除文件,如果传入-r标志则可递归删除文件夹`,
	RunE: func(cmd *cobra.Command, args []string) error {
		deleteDir, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}
		if len(args) < 1 {
			return errors.New("缺少操作数")
		}
		return call.DeleteFile(cmd.Context(), args[0], true, deleteDir, false)
	},
}

func init() {
	rmCmd.Flags().BoolP("recursive", "r", false, "用户名,如果为空则为root")
	rootCmd.AddCommand(rmCmd)
}
