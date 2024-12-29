package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// loginCmd 登录
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "用户登录",
	Long:  `进行用户登录,如果不指明用户名则默认登录root根用户,允许后置输入密码`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}
		log.Debugln("username:", username, ", password:", password)

		if username != "root" && password == "" {
			// TODO tui输入密码
		}
		res := call.Login(cmd.Context(), username, password)
		if res {
			return nil
		} else {
			return errors.New("用户名或密码错误")
		}
	},
}

func init() {
	loginCmd.Flags().StringP("username", "u", "root", "用户名,如果为空则为root")
	loginCmd.Flags().StringP("password", "p", "", "密码,如果为空则后续手动输入")

	rootCmd.AddCommand(loginCmd)
}
