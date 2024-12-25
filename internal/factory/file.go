package factory

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func CreateFile(content, owner string, chmod int) bson.M {
	return bson.M{
		"content":  content,
		"type":     "file",
		"size":     len([]byte(content)) / 8,
		"owner":    owner,
		"chmod":    chmod,
		"created":  time.Now(),
		"modified": time.Now(),
		"updated":  time.Now(),
	}
}

func CreateDir(owner string, chmod int) bson.M {
	return bson.M{
		"content": bson.M{},
		"type":    "dir",
		"size":    0,
		"owner":   owner,
		"chmod":   chmod,
		"created": time.Now(),
		"updated": time.Now(),
	}
}
