package cmd

import (
	"StarFileManager/internal/call"
	"errors"

	"github.com/spf13/cobra"
)

var mvCmd = &cobra.Command{
	Use:   "mv",
	Short: "Move files",
	Long:  `Move the file of the first parameter to the location specified by the second parameter`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("missing arguments")
		}

		err := call.MoveFile(cmd.Context(), args[0], args[1])
		return err
	},
}

func init() {
	rootCmd.AddCommand(mvCmd)
}
