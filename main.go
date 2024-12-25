package main

import (
	"StarFileManager/cmd"
	"StarFileManager/internal/factory"
	"context"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	godotenv.Load()

	// 配置日志
	log.SetFormatter(&nested.Formatter{
		HideKeys: true,
	})
	if os.Getenv("debug") == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.FatalLevel)
	}

	// mongoDB数据库连接
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", os.Getenv("mgHost")))
	// 连接到MongoDB
	mg, err := mongo.Connect(nil, clientOptions)
	if err != nil {
		log.Fatalln("mongoDB数据库连接失败", err)
	} else {
		defer mg.Disconnect(nil)
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()
	if err := mg.Ping(timeoutCtx, nil); err != nil {
		log.Fatalln("mongoDB数据库连接失败", err)
	}
	ctx = context.WithValue(ctx, "mongo", mg)

	// redis数据库连接
	var re = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", os.Getenv("reHost")), // Redis 服务器地址
	})
	if err := re.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("redis数据库连接失败", err)
	}
	ctx = context.WithValue(ctx, "redis", re)

	// 创建根目录
	files := mg.Database("starfile").Collection("files")
	if err := files.FindOne(context.Background(), bson.M{"_id": os.Getenv("rootInode")}).Err(); err != nil {
		root := factory.CreateDir("root", 0755)
		root["_id"] = os.Getenv("rootInode")
		// 创建home目录
		home := factory.CreateDir("root", 0755)
		res, _ := files.InsertOne(context.Background(), home)
		homeId := res.InsertedID
		root["content"].(bson.M)["home"] = homeId.(primitive.ObjectID)
		files.InsertOne(context.Background(), root)
	}

	// 当前登录用户
	user, err := re.Get(context.Background(), string(os.Getppid())).Result()
	if err == nil && user != "" {
		ctx = context.WithValue(ctx, "user", user)
		log.Debugln("当前用户", ctx.Value("user"))
	}

	// 通过上下文注入依赖
	cmd.ExecuteContext(ctx)
}
