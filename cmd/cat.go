package cmd

import (
	"StarFileManager/internal/call"
	"StarFileManager/internal/view"
	"errors"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// catCmd 查看文件内容
var catCmd = &cobra.Command{
	Use:   "cat",
	Short: "show file content",
	Long:  `show file content`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing file name")
		}

		content, err := call.GetFileContent(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		// 显示内容
		p := tea.NewProgram(
			view.CatView{
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
