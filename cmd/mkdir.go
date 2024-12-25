package cmd

import (
	"StarFileManager/internal/call"
	"github.com/spf13/cobra"
)

// mkdirCommand 创建目录
var mkdirCommand = &cobra.Command{
	Use:   "mkdir",
	Short: "创建目录",
	Long:  `在指定目录下创建目录,`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, _ := cmd.Flags().GetBool("parents")
		return call.MakeDir(cmd.Context(), args[0], p)
	},
}

func init() {
	rootCmd.Flags().BoolP("parents", "p", false, "必要时创建父目录")

	rootCmd.AddCommand(mkdirCommand)
}
