package call

import (
	"StarFileManager/internal/factory"
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
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
	file, err := GetFile(ctx, target, true)
	if err != nil {
		return err
	}

	// 验证权限
	if !CheckMod(ctx, file, "x") {
		return errors.New("权限不够")
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
				// 验证权限
				if !CheckMod(ctx, current, "w") {
					return primitive.NilObjectID, errors.New("权限不够")
				}
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

	// 验证权限
	if !CheckMod(ctx, current, "w") {
		return primitive.NilObjectID, errors.New("权限不够")
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

// CopyDir 拷贝目录
func CopyDir(ctx context.Context, src, tar string) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")

	// 是否复制到自己
	if isSub := strings.HasPrefix(tar, src+"/"); isSub {
		return errors.New("无法将目录复制到自己")
	}

	// 找到源文件
	srcFile, err := GetFile(ctx, src, true)
	if err != nil {
		return err
	}
	src, err = GetRealPath(ctx, src)
	if err != nil {
		return err
	}
	// 验证源文件权限
	if !CheckMod(ctx, srcFile, "r") {
		return errors.New("权限不够")
	}

	// 找到目标位置父目录
	tarFile, err := GetFile(ctx, tar, false)
	if err != nil {
		return err
	}
	tar, err = GetRealPath(ctx, tar)
	if err != nil {
		return err
	}
	// 验证目标父目录权限
	if !CheckMod(ctx, tarFile, "w") {
		return errors.New("权限不够")
	}
	dirname := filepath.Base(tar)

	// 判断是否是目录
	if tarFile["type"] != "dir" {
		return errors.New("目标地址所在位置不是目录")
	}

	// 如果目标已存在报错
	if _, err := GetChildFile(ctx, tarFile, dirname); err == nil {
		return errors.New("目标已存在")
	}

	// 拷贝当前
	delete(srcFile, "_id")
	srcFile["time"] = time.Now()
	res, err := files.InsertOne(ctx, srcFile)
	if err != nil {
		return err
	}
	newId := res.InsertedID
	tarFile["content"].(bson.M)[dirname] = newId

	// 目标父目录保存id
	filter := bson.M{"_id": tarFile["_id"]}
	update := bson.M{"$set": bson.M{"content": tarFile["content"]}}
	files.UpdateOne(context.Background(), filter, update)

	// 递归拷贝子目录
	for filename, _ := range srcFile["content"].(bson.M) {
		current, err := GetChildFile(ctx, srcFile, filename)
		if err != nil {
			return err
		}
		if current["type"] == "dir" {
			// 拷贝目录
			err := CopyDir(ctx, fmt.Sprintf("%s/%s", src, filename), fmt.Sprintf("%s/%s", tar, filename))
			if err != nil {
				return err
			}
		} else {
			// 拷贝文件
			err := CopyFile(ctx, fmt.Sprintf("%s/%s", src, filename), fmt.Sprintf("%s/%s", tar, filename))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ListFiles 列出指定目录下的所有文件
func ListFiles(ctx context.Context, path string) ([]string, error) {
	path, err := GetRealPath(ctx, path)
	if err != nil {
		return nil, err
	}

	// 获取指定目录
	filename := filepath.Base(path)
	current, err := GetFile(ctx, path, true)
	if err != nil {
		return nil, err
	}
	// 验证权限
	if !CheckMod(ctx, current, "r") {
		return nil, errors.New("权限不够")
	}

	// 当前不是目录,输出当前文件
	if current["type"] != "dir" {
		return []string{filename}, nil
	} else {
		// 遍历当前目录下的文件
		var res []string
		for name, _ := range current["content"].(bson.M) {
			res = append(res, name)
		}
		return res, nil
	}
}

// ListFilesDetail 列出指定目录下的所有文件详细信息
func ListFilesDetail(ctx context.Context, path string) ([]table.Row, error) {
	path, err := GetRealPath(ctx, path)
	if err != nil {
		return nil, err
	}

	// 获取指定目录
	filename := filepath.Base(path)
	current, err := GetFile(ctx, path, true)
	if err != nil {
		return nil, err
	}
	// 验证权限
	if !CheckMod(ctx, current, "r") {
		return nil, errors.New("权限不够")
	}

	// 当前不是目录,输出当前文件
	if current["type"] != "dir" {
		return []table.Row{{
			GetModString(ctx, int(current["chmod"].(int32))),
			strconv.FormatInt(int64(GetHLinkCount(ctx, current["_id"].(primitive.ObjectID))), 10),
			current["owner"].(string),
			current["time"].(primitive.DateTime).Time().Format("2006-01-02 15:04:05"),
			current["type"].(string),
			filename,
		}}, nil
	} else {
		// 遍历当前目录下的文件
		var res []table.Row
		for name, _ := range current["content"].(bson.M) {
			file, err := GetChildFile(ctx, current, name)
			if err != nil {
				return nil, err
			}
			res = append(res, table.Row{
				GetModString(ctx, int(file["chmod"].(int32))),
				strconv.FormatInt(int64(GetHLinkCount(ctx, file["_id"].(primitive.ObjectID))), 10),
				file["owner"].(string),
				file["time"].(primitive.DateTime).Time().Format("2006-01-02 15:04:05"),
				file["type"].(string),
				name,
			})
		}
		return res, nil
	}
}
