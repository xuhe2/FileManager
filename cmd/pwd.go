package cmd

import (
	"StarFileManager/internal/call"
	"fmt"
	"github.com/spf13/cobra"
)

// pwdCmd 输出当前路径
var pwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "获取当前所在的文件路径",
	Long:  `获取当前所在的文件路径地址,需要先登录`,
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
