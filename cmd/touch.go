package cmd

import (
	"StarFileManager/internal/call"
	"errors"
	"github.com/spf13/cobra"
	"time"
)

var touchCmd = &cobra.Command{
	Use:   "touch",
	Short: "创建或访问文件",
	Long:  `访问指定的文件,更新修改时间戳,如果不存在则创建它`,
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
			return errors.New("缺少操作数")
		}
		_, err = call.MakeFile(cmd.Context(), args[0], t)
		return err
	},
}

func init() {
	touchCmd.Flags().StringP("stamp", "t", "", "指定修改时间戳(采用[[CC]YY]MMDDhhmm[.ss]格式")
	rootCmd.AddCommand(touchCmd)
}
