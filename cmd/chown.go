package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"github.com/spf13/cobra"
)

// chownCmd 修改文件所有者
var chownCmd = &cobra.Command{
	Use:   "chown",
	Short: "修改文件的所有者",
	Long:  `先输入修改后的所有者用户名,然后输入修改目标.可以使用-R标志递归修改目录下的所有文件的权限`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("缺少操作数")
		}
		r, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		return call.SetChown(cmd.Context(), args[1], args[0], r)
	},
}

func init() {
	chownCmd.Flags().BoolP("recursive", "R", false, "递归修改目录下的所有文件的所有者")
	rootCmd.AddCommand(chownCmd)
}
