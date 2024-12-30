package cmd

import (
	"StarFileManager/internal/call"
	"StarFileManager/internal/model"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"path/filepath"
)

// catCmd 查看文件内容
var catCmd = &cobra.Command{
	Use:   "cat",
	Short: "输出文件内容",
	Long:  `输出文件内容`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("缺少操作数")
		}

		content, err := call.GetFileContent(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		//fmt.Println(content)

		// 显示内容
		p := tea.NewProgram(
			model.CatView{
				Title:   filepath.Base(args[0]),
				Content: content,
			},
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)

		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catCmd)
}
