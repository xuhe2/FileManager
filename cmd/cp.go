package cmd

import (
	"StarFileManager/internal/call"
	"errors"

	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "copy file or dir",
	Long:  `Copies the file specified by the first parameter to the destination specified by the second parameter (ending with the file). Use the -r flag to recursively copy folders.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cpDir, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}
		if len(args) < 2 {
			return errors.New("missing arguments")
		}

		if cpDir {
			err := call.CopyDir(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
		} else {
			err := call.CopyFile(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	cpCmd.Flags().BoolP("recursive", "r", false, "Recursively copy folders")

	rootCmd.AddCommand(cpCmd)
}
