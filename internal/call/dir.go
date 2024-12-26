package call

import (
	"StarFileManager/internal/factory"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"path/filepath"
	"strings"
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

func ChangePath(ctx context.Context, target string) error {
	re := ctx.Value("redis").(*redis.Client)
	user := GetUser(ctx)
	now, err := GetRealPath(ctx, target)
	if err != nil {
		return err
	}
	_, err = re.Set(context.Background(), fmt.Sprintf("%s:%s", user, "path"), now, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetFile 获取指定目录下指定文件名的文件
func GetFile(ctx context.Context, fatherDir bson.M, filename string) (bson.M, error) {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")

	inodeId, ok := fatherDir["content"].(bson.M)[filename]
	if !ok {
		return nil, errors.New("文件不存在")
	}
	res := bson.M{}
	filter := bson.M{"_id": inodeId}
	err := files.FindOne(ctx, filter).Decode(&res)
	return res, err
}

// MakeDir 创建目录
func MakeDir(ctx context.Context, path string, isCreateP bool) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	re := ctx.Value("redis").(*redis.Client)
	files := mg.Database("starfile").Collection("files")
	path, err := GetRealPath(ctx, path)

	// 逐个确认目录
	var current bson.M
	files.FindOne(context.Background(), bson.M{"_id": os.Getenv("rootInode")}).Decode(&current)
	currentDir := ""
	segments := strings.Split(path[1:], "/")
	log.Debugln(segments)
	for _, segment := range segments[:len(segments)-1] {
		log.Debugln("当前目录:", current)
		current, err = GetFile(ctx, current, segment) // 获取下一层目录
		currentDir = currentDir + "/" + segment
		if err != nil {
			if isCreateP {
				// 创建目录
				err := MakeDir(ctx, currentDir, isCreateP)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// 创建当前目录
	dirname := segments[len(segments)-1]
	umask, err := re.Get(context.Background(), fmt.Sprintf("%s:umask", GetUser(ctx))).Int()
	log.Debugln("创建目录:", dirname, umask)
	if err != nil {
		return err
	}
	var inodeId primitive.ObjectID
	if res, err := files.InsertOne(ctx, factory.CreateDir(GetUser(ctx), 0777 & ^umask)); err != nil {
		return err
	} else {
		inodeId = res.InsertedID.(primitive.ObjectID)
	}

	//上级目录补充当前目录信息
	log.Debugln(current["_id"])
	filter := bson.M{"_id": current["_id"]}
	update := bson.M{"$set": bson.M{
		"content." + dirname: inodeId,
	}}
	files.UpdateOne(context.Background(), filter, update)

	return nil
}

func DeleteDir(ctx context.Context, path string, deleteOnlyEmpty bool) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")

	// 找到目标目录
	path, err := GetRealPath(ctx, path)
	log.Debugln(path)
	if err != nil {
		return err
	}
	var current bson.M
	var father bson.M
	files.FindOne(context.Background(), bson.M{"_id": os.Getenv("rootInode")}).Decode(&current)
	segments := strings.Split(path[1:], "/")
	log.Debugln(segments)
	for _, segment := range segments {
		log.Debugln("当前目录:", current)
		father = current
		current, err = GetFile(ctx, current, segment) // 获取下一层目录
		if err != nil {
			return err
		}
	}

	if len(current["content"].(bson.M)) != 0 {
		if deleteOnlyEmpty {
			return errors.New("目录非空,无法删除")
		} else {
			// TODO 递归删除
		}
	}

	// 删除目录
	delete(father["content"].(bson.M), segments[len(segments)-1])
	files.UpdateOne(context.Background(), bson.M{"_id": father["_id"]}, bson.M{"$set": father})
	files.DeleteOne(context.Background(), bson.M{"_id": current["_id"]})

	return nil
}
