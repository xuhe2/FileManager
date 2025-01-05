package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"time"

	"github.com/spf13/cobra"
)

// touchCmd 访问或创建文件
var touchCmd = &cobra.Command{
	Use:   "touch",
	Short: "create a file",
	Long:  `create or update the access and modification times of each FILE to the current time`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stamp := cmd.Flag("stamp").Value.String()
		var t time.Time
		var err error
		if stamp != "" {
			t, err = call.GetTimeFromStamp(stamp)
			if err != nil {
				return err
			}
		} else {
			t = time.Now()
		}

		if len(args) < 1 {
			return errors.New("missing file operand")
		}
		_, err = call.MakeFile(cmd.Context(), args[0], t)
		return err
	},
}

func init() {
	touchCmd.Flags().StringP("stamp", "t", "", "指定修改时间戳(采用[[CC]YY]MMDDhhmm[.ss]格式")
	rootCmd.AddCommand(touchCmd)
}
