package factory

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// CreateFile 创建文件
func CreateFile(blocks []primitive.ObjectID, owner string, stamp time.Time, chmod int) bson.M {
	return bson.M{
		"content": blocks,
		"type":    "file",
		"owner":   owner,
		"chmod":   chmod,
		"time":    stamp,
	}
}

// CreateDir 创建目录
func CreateDir(owner string, stamp time.Time, chmod int) bson.M {
	return bson.M{
		"content": bson.M{},
		"type":    "dir",
		"owner":   owner,
		"chmod":   chmod,
		"time":    stamp,
	}
}
