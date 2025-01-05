package cmd

import (
	"StarFileManager/internal/call"
	"StarFileManager/internal/view"
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 编辑文件
var editCmd = &cobra.Command{
	Use:   "vim",
	Short: "Editing file contents",
	Long:  `Edit the contents of the specified file, not the directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		re := cmd.Context().Value("redis").(*redis.Client)
		if len(args) < 1 {
			return errors.New("missing arguments")
		}

		// 获取文件
		path := args[0]
		path, err := call.GetRealPath(cmd.Context(), path)
		if err != nil {
			return err
		}
		file, err := call.GetFile(cmd.Context(), path, true)
		if err != nil {
			return err
		}
		// 验证权限
		if !call.CheckMod(cmd.Context(), file, "w") {
			return errors.New("Insufficient permissions")
		}

		// 判断类型
		if file["type"] != "file" {
			return errors.New("not a file")
		}

		// 读取文件内容
		content, err := call.GetFileContent(cmd.Context(), path)
		if err != nil {
			return err
		}

		// 上写锁
		if s, _ := re.SetNX(context.Background(), file["_id"].(primitive.ObjectID).String(), 1, 0).Result(); !s {
			return errors.New("The file is currently occupied")
		}
		defer re.Del(context.Background(), file["_id"].(primitive.ObjectID).String())

		// 启动编辑文件界面
		p := tea.NewProgram(view.NewEditArea(cmd.Context(), path, content))

		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
