package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "拷贝文件",
	Long:  `拷贝第一个参数指定的文件到第二个参数指定的目标位置(终止于文件),使用-r标志递归拷贝文件夹`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cpDir, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}
		if len(args) < 2 {
			return errors.New("缺少操作数")
		}

		if cpDir {
			err := call.CopyDir(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
		} else {
			err := call.CopyFile(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	cpCmd.Flags().BoolP("recursive", "r", false, "递归拷贝文件夹")

	rootCmd.AddCommand(cpCmd)
}
