package cmd

import (
	"StarFileManager/internal/call"
	"fmt"

	"github.com/spf13/cobra"
)

// pwdCmd 输出当前路径
var pwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "Get the current file path",
	Long:  `Get the current file path address, you need to log in first`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := call.GetPwd(cmd.Context())
		if err != nil {
			return err
		} else {
			fmt.Println(res)
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(pwdCmd)
}
