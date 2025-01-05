package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"strconv"

	"github.com/spf13/cobra"
)

// chmodCmd 修改文件权限
var chmodCmd = &cobra.Command{
	Use:   "chmod",
	Short: "change file mode",
	Long:  `use digit mode, use 3 digit to change mode.use '-R' flag to recursive change file mode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("缺少操作数")
		}
		chmod, err := strconv.ParseInt(args[0], 8, 32)
		if err != nil {
			return err
		}

		r, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		return call.SetChmod(cmd.Context(), args[1], int(chmod), r)
	},
}

func init() {
	chmodCmd.Flags().BoolP("recursive", "R", false, "recursive change file mode")

	rootCmd.AddCommand(chmodCmd)
}
