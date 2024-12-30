package cmd

import "github.com/spf13/cobra"

var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "拷贝文件",
	Long:  `拷贝第一个参数指定的文件到第二个参数指定的目标位置(终止于文件),使用-r标志递归拷贝文件夹`,
}
