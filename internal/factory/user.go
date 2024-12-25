package factory

import "go.mongodb.org/mongo-driver/bson"

// CreateUser 创建用户对象工厂
func CreateUser(username string, password string) bson.M {
	return bson.M{
		"username": username,
		"password": password,
		"umask":    0022,
	}
}
