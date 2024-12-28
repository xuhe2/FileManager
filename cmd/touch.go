package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"github.com/spf13/cobra"
)

var touchCmd = &cobra.Command{
	Use:   "touch",
	Short: "创建或访问文件",
	Long:  `访问指定的文件,如果不存在则创建它`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("缺少操作数")
		}
		_, err := call.MakeFile(cmd.Context(), args[0])
		return err
	},
}

func init() {
	rootCmd.AddCommand(touchCmd)
}
