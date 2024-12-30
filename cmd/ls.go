package cmd

import (
	"StarFileManager/internal/call"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出指定目录下的目录",
	Long:  `列出指定目录下的所有文件,可以使用-l标志列出详细信息.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showDetail, err := cmd.Flags().GetBool("long")
		if err != nil {
			return err
		}

		target := ""
		if len(args) < 1 {
			target = "."
		} else {
			target = args[0]
		}
		call.ListFiles(cmd.Context(), target, showDetail)
		return nil
	},
}

func init() {
	lsCmd.Flags().BoolP("long", "l", false, "")

	rootCmd.AddCommand(lsCmd)
}
