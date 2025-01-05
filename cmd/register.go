package cmd

import (
	"StarFileManager/internal/call"
	"StarFileManager/internal/view"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// registerCmd 注册
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register",
	Long:  `For user login, the user name must be included, duplicate registration is not allowed, and password entry is allowed later`,
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

		if password == "" {
			i := view.NewPasswordInput(cmd.Context(), username, call.Register)
			p := tea.NewProgram(i)
			if _, err := p.Run(); err != nil {
				return err
			}
			return nil
		}
		err = call.Register(cmd.Context(), username, password)
		return err
	},
}

func init() {
	registerCmd.Flags().StringP("username", "u", "", "Username, must contain")
	registerCmd.MarkFlagRequired("username")
	registerCmd.Flags().StringP("password", "p", "", "Password, if it is empty, enter it manually later")

	rootCmd.AddCommand(registerCmd)
}
