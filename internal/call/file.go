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
	"time"
)

// MakeFile 创建文件
func MakeFile(ctx context.Context, path string, stamp time.Time) (primitive.ObjectID, error) {
	mg := ctx.Value("mongo").(*mongo.Client)
	re := ctx.Value("redis").(*redis.Client)
	files := mg.Database("starfile").Collection("files")

	// 逐个进入目录
	current, err := GetFile(ctx, path, false)

	// 验证是否重复创建
	log.Debugln(current["_id"])
	filename := filepath.Base(path)
	if _, ok := current["content"].(bson.M)[filename]; ok {
		//重复创建,仅修改访问与创建时间
		file, err := GetChildFile(ctx, current, filename)
		if err != nil {
			return primitive.NilObjectID, err
		}
		inodeId := file["_id"].(primitive.ObjectID)
		filter := bson.M{"_id": inodeId}
		update := bson.M{"$set": bson.M{"time": stamp}}
		files.UpdateOne(context.Background(), filter, update)
		return inodeId, nil
	}

	// 创建当前目录
	umask, err := re.Get(context.Background(), fmt.Sprintf("%s:umask", GetUser(ctx))).Int()
	log.Debugln("创建文件:", filename, umask)
	if err != nil {
		return primitive.NilObjectID, err
	}
	var inodeId primitive.ObjectID
	if res, err := files.InsertOne(ctx, factory.CreateFile("", GetUser(ctx), stamp, 0666 & ^umask)); err != nil {
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

// GetFile 获取指定的文件
func GetFile(ctx context.Context, path string, includeLast bool) (bson.M, error) {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")

	// 获取绝对路径
	path, err := GetRealPath(ctx, path)
	log.Debugln(path)
	if err != nil {
		return nil, err
	}

	// 从根目录开始寻找
	var current bson.M
	files.FindOne(context.Background(), bson.M{"_id": os.Getenv("rootInode")}).Decode(&current)
	segments := strings.Split(path[1:], "/")
	log.Debugln(segments)
	// 是否需要包括最后一级目录
	if !includeLast {
		segments = segments[:len(segments)-1]
	}

	for _, segment := range segments {
		log.Debugln("当前目录:", current)
		// 当前不是目录
		if current["type"] != "dir" {
			return nil, errors.New("不是目录")
		}

		// 获取下一层目录
		current, err = GetChildFile(ctx, current, segment)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
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
	path, err := GetRealPath(ctx, path)
	if err != nil {
		return "", err
	}

	current, err := GetFile(ctx, path, true)
	return current["type"].(string), err
}

// GetFileContent 获取文件内容
func GetFileContent(ctx context.Context, path string) (string, error) {
	// 找到目标文件
	current, err := GetFile(ctx, path, true)
	if err != nil {
		return "", err
	}

	if current["type"] != "file" {
		return "", errors.New("不可获取文件夹内容")
	}
	return current["content"].(string), nil
}

// SetChmod 设置文件访问权限
func SetChmod(ctx context.Context, path string, chmod int, modifyInner bool) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")

	// 获取目标文件
	current, err := GetFile(ctx, path, true)
	if err != nil {
		return err
	}

	// 修改文件权限
	filter := bson.M{"_id": current["_id"]}
	update := bson.M{"$set": bson.M{"chmod": chmod}}
	files.UpdateOne(context.Background(), filter, update)
	if current["type"] == "dir" && modifyInner {
		// 递归修改所有子目录的权限
		for name, _ := range current["content"].(bson.M) {
			err := SetChmod(ctx, filepath.Join(path, name), chmod, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetChown 设置文件所有者
func SetChown(ctx context.Context, path string, owner string, modifyInner bool) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")
	users := mg.Database("starfile").Collection("users")

	// 验证用户是否存在
	filter := bson.M{"username": owner}
	if users.FindOne(ctx, filter).Err() != nil {
		return errors.New("用户不存在")
	}

	// 获取目标文件
	current, err := GetFile(ctx, path, true)
	if err != nil {
		return err
	}

	// 修改文件所有者
	filter = bson.M{"_id": current["_id"]}
	update := bson.M{"$set": bson.M{"owner": owner}}
	files.UpdateOne(context.Background(), filter, update)
	if current["type"] == "dir" && modifyInner {
		// 递归修改所有子目录的权限
		for name, _ := range current["content"].(bson.M) {
			err := SetChown(ctx, filepath.Join(path, name), owner, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteFile 删除文件(目录)
func DeleteFile(ctx context.Context, path string, deleteFile bool, deleteDir bool, deleteOnlyEmptyDir bool) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	files := mg.Database("starfile").Collection("files")

	// 找到父目录
	father, err := GetFile(ctx, path, false)
	if err != nil {
		return err
	}
	// 找到目标文件
	filename := filepath.Base(path)
	current, err := GetChildFile(ctx, father, filename)

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
	log.Debugln("删除", filename)
	delete(father["content"].(bson.M), filename)
	files.UpdateOne(context.Background(), bson.M{"_id": father["_id"]}, bson.M{"$set": father})
	files.DeleteOne(context.Background(), bson.M{"_id": current["_id"]})
	nowPath, _ := GetPwd(ctx)
	if nowPath == path {
		ChangePath(ctx, fmt.Sprintf("/home/%s", GetUser(ctx)))
	}

	return nil
}
