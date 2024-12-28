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

// MakeFile 创建目录
func MakeFile(ctx context.Context, path string) (primitive.ObjectID, error) {
	mg := ctx.Value("mongo").(*mongo.Client)
	re := ctx.Value("redis").(*redis.Client)
	files := mg.Database("starfile").Collection("files")
	path, err := GetRealPath(ctx, path)

	// 逐个进入目录
	var current bson.M
	files.FindOne(context.Background(), bson.M{"_id": os.Getenv("rootInode")}).Decode(&current)
	currentDir := ""
	segments := strings.Split(path[1:], "/")
	log.Debugln(segments)
	for _, segment := range segments[:len(segments)-1] {
		log.Debugln("当前目录:", current)
		current, err = GetChildFile(ctx, current, segment) // 获取下一层目录
		currentDir = currentDir + "/" + segment
		if err != nil {
			return primitive.NilObjectID, err
		}
		if current["type"] != "dir" {
			return primitive.NilObjectID, errors.New("不是目录")
		}
	}

	// 验证是否重复创建
	log.Debugln(current["_id"])
	filename := segments[len(segments)-1]
	if _, ok := current["content"].(bson.M)[filename]; ok {
		return primitive.NilObjectID, errors.New("文件名重复")
	}

	// 创建当前目录
	umask, err := re.Get(context.Background(), fmt.Sprintf("%s:umask", GetUser(ctx))).Int()
	log.Debugln("创建文件:", filename, umask)
	if err != nil {
		return primitive.NilObjectID, err
	}
	var inodeId primitive.ObjectID
	if res, err := files.InsertOne(ctx, factory.CreateFile("", GetUser(ctx), 0666 & ^umask)); err != nil {
		return primitive.NilObjectID, err
	} else {
		inodeId = res.InsertedID.(primitive.ObjectID)
	}

	//上级目录补充当前文件信息
	filter := bson.M{"_id": current["_id"]}
	update := bson.M{"$set": bson.M{
		"content." + filename: inodeId,
	}}
	files.UpdateOne(context.Background(), filter, update)

	return inodeId, nil
}

// GetChildFile 获取指定目录下指定文件名的文件
func GetChildFile(ctx context.Context, fatherDir bson.M, filename string) (bson.M, error) {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")

	inodeId, ok := fatherDir["content"].(bson.M)[filename]
	if !ok {
		return nil, errors.New("没有那个文件或目录")
	}
	res := bson.M{}
	filter := bson.M{"_id": inodeId}
	err := files.FindOne(ctx, filter).Decode(&res)
	return res, err
}

// GetFileType 获取文件类型
func GetFileType(ctx context.Context, path string) (string, error) {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")

	path, err := GetRealPath(ctx, path)
	if err != nil {
		return "", err
	}

	var current bson.M
	files.FindOne(context.Background(), bson.M{"_id": os.Getenv("rootInode")}).Decode(&current)
	segments := strings.Split(path[1:], "/")
	log.Debugln(segments)
	for i, segment := range segments {
		log.Debugln("当前目录:", current)
		current, err = GetChildFile(ctx, current, segment) // 获取下一层目录
		if err != nil {
			return "", err
		}
		if i != len(segments)-1 && current["type"] != "dir" {
			return "", errors.New("不是目录")
		}
	}

	return current["type"].(string), nil
}

// DeleteFile 删除文件(目录)
func DeleteFile(ctx context.Context, path string, deleteFile bool, deleteDir bool, deleteOnlyEmptyDir bool) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")

	// 找到目标文件
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
	for i, segment := range segments {
		log.Debugln("当前目录:", current)
		father = current
		current, err = GetChildFile(ctx, current, segment) // 获取下一层目录
		if err != nil {
			return err
		}
		if i != len(segments)-1 && current["type"] != "dir" {
			return errors.New("不是目录")
		}
	}

	// 判断目标类型
	if current["type"] == "dir" {
		if !deleteDir {
			return errors.New("目标是文件夹,无法删除")
		} else if len(current["content"].(bson.M)) != 0 {
			// 当前为文件夹,判断是否可以删非空
			if deleteOnlyEmptyDir {
				return errors.New("目录非空,无法删除")
			} else {
				// 删除所有子目录
				for name, _ := range current["content"].(bson.M) {
					err := DeleteFile(ctx, filepath.Join(path, name), true, true, false)
					if err != nil {
						return err
					}
				}
			}
		}
	} else if current["type"] == "file" {
		if !deleteFile {
			return errors.New("目标是文件,无法删除")
		}
	}

	// 删除目标文件
	log.Debugln("删除", segments[len(segments)-1])
	delete(father["content"].(bson.M), segments[len(segments)-1])
	files.UpdateOne(context.Background(), bson.M{"_id": father["_id"]}, bson.M{"$set": father})
	files.DeleteOne(context.Background(), bson.M{"_id": current["_id"]})
	nowPath, _ := GetPwd(ctx)
	if nowPath == path {
		ChangePath(ctx, fmt.Sprintf("/home/%s", GetUser(ctx)))
	}

	return nil
}
