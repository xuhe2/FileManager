package cmd

import (
	"StarFileManager/internal/call"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"strconv"
)

// umaskCmd 修改用户mask
var umaskCmd = &cobra.Command{
	Use:   "umask",
	Short: "查看或修改当前用户的user mask",
	Long:  `不传参数时查看当前用户的user mask,传参数时修改当前用户的user mask.采用8进制数形式`,
	RunE: func(cmd *cobra.Command, args []string) error {
		re := cmd.Context().Value("redis").(*redis.Client)
		username := call.GetUser(cmd.Context())
		umask, err := re.Get(context.Background(), fmt.Sprintf("%s:umask", username)).Int64()
		if len(args) < 1 {
			fmt.Printf("%04o\n", umask)
			return nil
		}
		umask, err = strconv.ParseInt(args[0], 8, 32)
		if err != nil {
			return err
		}
		return call.SetUmask(cmd.Context(), int(umask))
	},
}

func init() {
	rootCmd.AddCommand(umaskCmd)
}
