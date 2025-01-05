package cmd

import (
	"StarFileManager/internal/call"
	"errors"

	"github.com/spf13/cobra"
)

// chownCmd 修改文件所有者
var chownCmd = &cobra.Command{
	Use:   "chown",
	Short: "change file owner",
	Long:  `First enter the modified owner username, then enter the modification target. You can use the '-R' flag to recursively modify the permissions of all files in the directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("缺少操作数")
		}
		r, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		return call.SetChown(cmd.Context(), args[1], args[0], r)
	},
}

func init() {
	chownCmd.Flags().BoolP("recursive", "R", false, "recursively change file owner")
	rootCmd.AddCommand(chownCmd)
}
