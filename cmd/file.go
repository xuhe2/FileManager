package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "查看文件类型",
	Long:  `查看参数所指定的文件的类型`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("缺少操作数")
		}
		typename, err := call.GetFileType(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		fmt.Println(typename)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
}
