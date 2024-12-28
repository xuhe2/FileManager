package factory

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func CreateFile(content, owner string, stamp time.Time, chmod int) bson.M {
	return bson.M{
		"content": content,
		"type":    "file",
		"size":    len([]byte(content)) / 8,
		"owner":   owner,
		"chmod":   chmod,
		"time":    stamp,
	}
}

func CreateDir(owner string, stamp time.Time, chmod int) bson.M {
	return bson.M{
		"content": bson.M{},
		"type":    "dir",
		"size":    0,
		"owner":   owner,
		"chmod":   chmod,
		"time":    stamp,
	}
}
