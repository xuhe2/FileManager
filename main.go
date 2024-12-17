package main

import (
	"StarFileManager/cmd"
	"context"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
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
	context.WithValue(ctx, "mongo", mg)

	// redis数据库连接
	var re = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", os.Getenv("reHost")), // Redis 服务器地址
	})
	if err := re.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("redis数据库连接失败", err)
	}
	context.WithValue(ctx, "redis", re)

	// 文件保存位置
	err = os.MkdirAll("data/files", os.ModePerm)
	if err != nil {
		log.Fatalln("创建图片缓存目录失败:", err)
	}

	// 通过上下文注入依赖
	cmd.ExecuteContext(ctx)
}
