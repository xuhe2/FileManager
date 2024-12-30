package cmd

import (
	"StarFileManager/internal/call"
	"StarFileManager/internal/view"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// 编辑文件
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "编辑文件内容",
	Long:  `编辑指定的文件的内容,不可编辑目录`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("缺少操作数")
		}

		// 获取文件
		path := args[0]
		path, err := call.GetRealPath(cmd.Context(), path)
		if err != nil {
			return err
		}
		file, err := call.GetFile(cmd.Context(), path, true)
		if err != nil {
			return err
		}

		// 判断类型
		if file["type"] != "file" {
			return errors.New("目标必须是文件")
		}

		// 启动编辑文件界面
		p := tea.NewProgram(view.NewEditArea(cmd.Context(), path, file["content"].(string)))

		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
