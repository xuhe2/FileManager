package call

import (
	"StarFileManager/internal/factory"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

// GetUser 是否登录
func GetUser(ctx context.Context) string {
	user, ok := ctx.Value("user").(string)
	if !ok {
		user = ""
	}
	return user
}

// Login 用户登录
func Login(ctx context.Context, username string, password string) bool {
	mg := ctx.Value("mongo").(*mongo.Client)
	re := ctx.Value("redis").(*redis.Client)
	users := mg.Database("starfile").Collection("users")

	// 检查用户名密码
	filter := bson.M{"username": username, "password": password}
	var user bson.M
	//users.InsertOne(ctx, filter)
	if err := users.FindOne(context.Background(), filter).Decode(&user); err != nil {
		return false
	} else {
		// 写入缓存
		//标记当前会话对应当前用户
		re.Set(context.Background(), string(os.Getppid()), username, 0)

		// user mask
		re.Set(context.Background(), fmt.Sprintf("%s:umask", username), user["umask"], 0)

		// 初始路径(主目录)
		homepath := fmt.Sprintf("/home/%s", username)
		re.Set(context.Background(), fmt.Sprintf("%s:path", username), homepath, 0)
		// TODO 写入其他缓存初始值
		return true
	}
}

// Register 用户注册
func Register(ctx context.Context, username string, password string) bool {
	mg := ctx.Value("mongo").(*mongo.Client)
	users := mg.Database("starfile").Collection("users")
	files := mg.Database("starfile").Collection("files")

	// 验证是否重复注册
	filter := bson.M{"username": username}
	if err := users.FindOne(context.Background(), filter).Err(); err == nil {
		return false
	}

	// 写入用户
	user := factory.CreateUser(username, password)
	users.InsertOne(context.Background(), user)

	// 创建主目录
	var root bson.M
	files.FindOne(context.Background(), bson.M{"_id": os.Getenv("rootInode")}).Decode(&root)
	home, _ := GetChildFile(ctx, root, "home")
	userpath := factory.CreateDir(username, time.Now(), 0750)
	res, _ := files.InsertOne(context.Background(), userpath)
	userpathId := res.InsertedID
	filter = bson.M{"_id": home["_id"]}
	update := bson.M{"$set": bson.M{
		"content." + username: userpathId,
	}}
	files.UpdateOne(context.Background(), filter, update)

	return true
}

// SetUmask 设置用户mask
func SetUmask(ctx context.Context, umask int) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	re := ctx.Value("redis").(*redis.Client)
	users := mg.Database("starfile").Collection("users")
	username := GetUser(ctx)

	// 验证umask
	if umask > 0777 {
		return errors.New("username mask应为3位8进制")
	}

	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"umask": umask}}
	users.UpdateOne(context.Background(), filter, update)
	re.Get(context.Background(), fmt.Sprintf("%s:umask", username))
	re.Set(context.Background(), fmt.Sprintf("%s:umask", username), umask, 0)
	return nil
}
