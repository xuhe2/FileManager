package factory

import "go.mongodb.org/mongo-driver/bson"

// CreatBlock 创建磁盘块
func CreatBlock(level int, content interface{}) bson.M {
	return bson.M{
		"level":   level,
		"content": content,
	}
}
