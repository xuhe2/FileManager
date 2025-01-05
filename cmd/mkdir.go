package cmd

import (
	"StarFileManager/internal/call"
	"errors"

	"github.com/spf13/cobra"
)

// mkdirCmd 创建目录
var mkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "make dir",
	Long:  `make dir`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, _ := cmd.Flags().GetBool("parents")
		if len(args) < 1 {
			return errors.New("missing dir name")
		}
		_, err := call.MakeDir(cmd.Context(), args[0], p)
		return err
	},
}

func init() {
	mkdirCmd.Flags().BoolP("parents", "p", false, "Allows automatic creation of parent directories when necessary")

	rootCmd.AddCommand(mkdirCmd)
}
