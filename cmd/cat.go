package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
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
		fmt.Println(content)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catCmd)
}
