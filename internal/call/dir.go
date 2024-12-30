package call

import (
	"StarFileManager/internal/factory"
	"StarFileManager/internal/model"
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GetPwd 获取当前所在的目录
func GetPwd(ctx context.Context) (string, error) {
	re := ctx.Value("redis").(*redis.Client)
	user := GetUser(ctx)
	res, err := re.Get(context.Background(), fmt.Sprintf("%s:%s", user, "path")).Result()
	return res, err
}

// GetRealPath 获取目录所对应的绝对路径
func GetRealPath(ctx context.Context, path string) (string, error) {
	// 绝对路径,直接返回
	if filepath.IsAbs(path) {
		return path, nil
	}

	// 相对路径,拼接
	pwd, err := GetPwd(ctx)
	if err != nil {
		return "", err
	}
	return filepath.Join(pwd, path), nil
}

// ChangePath 修改当前路径
func ChangePath(ctx context.Context, target string) error {
	re := ctx.Value("redis").(*redis.Client)

	user := GetUser(ctx)
	path, err := GetRealPath(ctx, target)
	if err != nil {
		return err
	}

	// 获取目标目录
	_, err = GetFile(ctx, target, true)
	if err != nil {
		return err
	}

	// 修改缓存
	_, err = re.Set(context.Background(), fmt.Sprintf("%s:%s", user, "path"), path, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// MakeDir 创建目录
func MakeDir(ctx context.Context, path string, isCreateP bool) (primitive.ObjectID, error) {
	mg := ctx.Value("mongo").(*mongo.Client)
	re := ctx.Value("redis").(*redis.Client)
	files := mg.Database("starfile").Collection("files")
	path, err := GetRealPath(ctx, path)

	// 逐个确认目录
	var current bson.M
	// 获取根目录
	files.FindOne(context.Background(), bson.M{"_id": os.Getenv("rootInode")}).Decode(&current)
	currentDir := ""
	segments := strings.Split(path[1:], "/")
	log.Debugln(segments)
	for _, segment := range segments[:len(segments)-1] {
		if current["type"] != "dir" {
			return primitive.NilObjectID, errors.New("不是目录")
		}
		log.Debugln("当前目录:", current)

		// 进入下一层目录
		var next bson.M
		next, err = GetChildFile(ctx, current, segment)
		currentDir = currentDir + "/" + segment
		// 目录不存在
		if err != nil {
			if isCreateP {
				// 创建目录
				id, err := MakeDir(ctx, currentDir, isCreateP)
				if err != nil {
					return primitive.NilObjectID, err
				}
				files.FindOne(context.Background(), bson.M{"_id": id}).Decode(&next)
			} else {
				return primitive.NilObjectID, err
			}
		}
		current = next
	}

	// 验证是否重复创建
	dirname := segments[len(segments)-1]
	log.Debugln(current["_id"])
	if _, ok := current["content"].(bson.M)[dirname]; ok {
		return primitive.NilObjectID, errors.New("名称重复")
	}

	// 创建当前目录
	umask, err := re.Get(context.Background(), fmt.Sprintf("%s:umask", GetUser(ctx))).Int()
	log.Debugln("创建目录:", dirname, umask)
	if err != nil {
		return primitive.NilObjectID, err
	}
	var inodeId primitive.ObjectID
	if res, err := files.InsertOne(ctx, factory.CreateDir(GetUser(ctx), time.Now(), 0777 & ^umask)); err != nil {
		return primitive.NilObjectID, err
	} else {
		inodeId = res.InsertedID.(primitive.ObjectID)
	}

	//上级目录补充当前目录信息
	filter := bson.M{"_id": current["_id"]}
	update := bson.M{"$set": bson.M{
		"content." + dirname: inodeId,
	}}
	files.UpdateOne(context.Background(), filter, update)

	return inodeId, nil
}

// ListFiles 列出指定目录下的所有文件信息
func ListFiles(ctx context.Context, path string, showDetail bool) error {
	path, err := GetRealPath(ctx, path)
	if err != nil {
		return err
	}

	// 获取指定目录
	filename := filepath.Base(path)
	current, err := GetFile(ctx, path, true)
	if err != nil {
		return err
	}

	// 表头
	cols := []table.Column{
		{Title: "权限", Width: 10},
		{Title: "硬连接数", Width: 10},
		{Title: "所有者", Width: 10},
		{Title: "编辑时间", Width: 30},
		{Title: "类型", Width: 10},
		{Title: "文件名", Width: 10},
	}

	// 当前不是目录,输出当前文件
	if current["type"] != "dir" {
		if showDetail {
			rows := []table.Row{
				{
					GetModString(ctx, int(current["chmod"].(int32))),
					strconv.FormatInt(int64(GetHLinkCount(ctx, current["_id"].(primitive.ObjectID))), 10),
					current["owner"].(string),
					current["time"].(primitive.DateTime).Time().Format("2006-01-02 15:04:05"),
					current["type"].(string),
					filename,
				},
			}
			t := table.New(
				table.WithColumns(cols),
				table.WithRows(rows),
				table.WithFocused(false),
				table.WithHeight(5),
			)

			// 设置样式
			borderStyle := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, true, false)
			defaultStyle := lipgloss.NewStyle()
			t.SetStyles(
				table.Styles{
					Header:   borderStyle,
					Cell:     defaultStyle,
					Selected: defaultStyle,
				},
			)

			m := model.LsTable{Table: t}
			if _, err := tea.NewProgram(m).Run(); err != nil {
				return err
			}
			return nil
		} else {
			fmt.Println(filename)
		}
	} else {
		// 遍历当前目录下的文件
		var rows []table.Row
		for name, _ := range current["content"].(bson.M) {
			if showDetail {
				file, err := GetChildFile(ctx, current, name)
				if err != nil {
					return err
				}
				rows = append(rows, table.Row{
					GetModString(ctx, int(file["chmod"].(int32))),
					strconv.FormatInt(int64(GetHLinkCount(ctx, file["_id"].(primitive.ObjectID))), 10),
					file["owner"].(string),
					file["time"].(primitive.DateTime).Time().Format("2006-01-02 15:04:05"),
					file["type"].(string),
					name,
				})
			} else {
				fmt.Printf("%s\t", name)
			}
		}

		// 显示表格
		if showDetail {
			t := table.New(
				table.WithColumns(cols),
				table.WithRows(rows),
				table.WithFocused(false),
				table.WithHeight(5),
			)

			// 设置样式
			borderStyle := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, true, false)
			defaultStyle := lipgloss.NewStyle()
			t.SetStyles(
				table.Styles{
					Header:   borderStyle,
					Cell:     defaultStyle,
					Selected: defaultStyle,
				},
			)

			m := model.LsTable{Table: t}
			if _, err := tea.NewProgram(m).Run(); err != nil {
				return err
			}
		} else {
			fmt.Println()
		}
	}

	return nil
}
