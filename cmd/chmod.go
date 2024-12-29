package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"github.com/spf13/cobra"
	"strconv"
)

// chmodCmd 修改文件权限
var chmodCmd = &cobra.Command{
	Use:   "chmod",
	Short: "修改指定文件的权限",
	Long:  `采用数字模式,用三位八进制表示权限,然后输入修改目标.可以使用-R标志递归修改目录下的所有文件的权限`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("缺少操作数")
		}
		chmod, err := strconv.ParseInt(args[0], 8, 32)
		if err != nil {
			return err
		}

		r, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		return call.SetChmod(cmd.Context(), args[1], int(chmod), r)
	},
}

func init() {
	chmodCmd.Flags().BoolP("recursive", "R", false, "递归修改目录下的所有文件的权限")

	rootCmd.AddCommand(chmodCmd)
}
