package cmd

import (
	"StarFileManager/internal/call"
	"StarFileManager/internal/view"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出指定目录下的目录",
	Long:  `列出指定目录下的所有文件,可以使用-l标志列出详细信息.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showDetail, err := cmd.Flags().GetBool("long")
		if err != nil {
			return err
		}

		// 获取查询目标
		target := ""
		if len(args) < 1 {
			// 默认为当前目录
			target = "."
		} else {
			target = args[0]
		}

		if showDetail {
			res, err := call.ListFilesDetail(cmd.Context(), target)
			if err != nil {
				return err
			}
			// 显示界面
			t := view.NewLsTable(res)
			if _, err := tea.NewProgram(t).Run(); err != nil {
				return err
			}
		} else {
			res, err := call.ListFiles(cmd.Context(), target)
			if err != nil {
				return err
			}
			for _, item := range res {
				fmt.Printf("%s\t", item)
			}
			fmt.Println()
		}
		return nil
	},
}

func init() {
	lsCmd.Flags().BoolP("long", "l", false, "")

	rootCmd.AddCommand(lsCmd)
}
