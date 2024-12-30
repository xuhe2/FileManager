package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"github.com/spf13/cobra"
)

var mvCmd = &cobra.Command{
	Use:   "mv",
	Short: "移动文件",
	Long:  `移动文件到指定的位置`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("缺少操作数")
		}

		err := call.MoveFile(cmd.Context(), args[0], args[1])
		return err
	},
}

func init() {
	rootCmd.AddCommand(mvCmd)
}
