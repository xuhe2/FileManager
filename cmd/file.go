package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// fileCmd 查看文件类型
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Check the file type",
	Long:  `Check the type of file specified by the parameter`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing file path")
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
