package call

import (
	"StarFileManager/internal/factory"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateFileBlocks 创建文件时完成块的创建
func CreateFileBlocks(ctx context.Context) ([]primitive.ObjectID, error) {
	mg := ctx.Value("mongo").(*mongo.Client)
	blocks := mg.Database("starfile").Collection("blocks")

	var res []primitive.ObjectID

	// 创建1级块
	create1b := func() ([]primitive.ObjectID, error) {
		var res []primitive.ObjectID
		// 10个一级块
		for i := 0; i < 10; i++ {
			block := factory.CreatBlock(1, "")
			r, err := blocks.InsertOne(context.Background(), block)
			if err != nil {
				return nil, err
			}
			res = append(res, r.InsertedID.(primitive.ObjectID))
		}
		return res, nil
	}
	b1, err := create1b()
	if err != nil {
		return nil, err
	}
	res = append(res, b1...)

	// 创建2级块(包含10个一级块)
	create2b := func() ([]primitive.ObjectID, error) {
		// 创建1级块
		b1, err := create1b()
		if err != nil {
			return nil, err
		}

		// 创建自身
		b2 := factory.CreatBlock(2, b1)
		r, err := blocks.InsertOne(context.Background(), b2)
		if err != nil {
			return nil, err
		}
		return []primitive.ObjectID{r.InsertedID.(primitive.ObjectID)}, nil
	}
	b2, err := create2b()
	if err != nil {
		return nil, err
	}
	res = append(res, b2...)

	// 创建3级块(包含2个二级块)
	var b2s []primitive.ObjectID
	for i := 0; i < 2; i++ {
		b2, err := create2b()
		if err != nil {
			return nil, err
		}
		b2s = append(b2s, b2...)
	}
	b3 := factory.CreatBlock(3, b2s)
	r, err := blocks.InsertOne(context.Background(), b3)
	if err != nil {
		return nil, err
	}
	res = append(res, r.InsertedID.(primitive.ObjectID))
	return res, nil
}

// ReadFileBlocks 读取块内容
func ReadFileBlocks(ctx context.Context, block bson.A) (string, error) {
	mg := ctx.Value("mongo").(*mongo.Client)
	blocks := mg.Database("starfile").Collection("blocks")

	res := ""
	for _, b := range block {
		// 读取当前块
		filter := bson.M{"_id": b}
		var bl bson.M
		blocks.FindOne(ctx, filter).Decode(&bl)

		if bl["level"].(int32) == 1 {
			if bl["content"].(string) == "" {
				// 当前块已空,后续块无内容读取结束
				return res, nil
			}
			// 1级块,直接访问
			res += bl["content"].(string)
		} else {
			// 多级块,递归继续读取
			re, err := ReadFileBlocks(ctx, bl["content"].(bson.A))
			if err != nil {
				return "", err
			}
			res += re
		}
	}
	return res, nil
}

// WriteFileBlocks 写入块内容
func WriteFileBlocks(ctx context.Context, block bson.A, content *string) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	blocks := mg.Database("starfile").Collection("blocks")

	for _, b := range block {
		// 读取当前块
		filter := bson.M{"_id": b}
		var bl bson.M
		blocks.FindOne(ctx, filter).Decode(&bl)

		if bl["level"].(int32) == 1 {
			if bl["content"].(string) == "" && *content == "" {
				// 当前块已空,后续块无内容写入结束
				return nil
			}
			// 1级块,直接写入
			bl["content"] = (*content)[0:min(1024, len(*content))]
			update := bson.M{"$set": bl}
			blocks.UpdateOne(ctx, filter, update)
			// 去掉已写入部分
			if len(*content) <= 1024 {
				*content = (*content)[0:0]
			} else {
				*content = (*content)[1024:]
			}
		} else {
			// 多级块,递归继续写入
			err := WriteFileBlocks(ctx, bl["content"].(bson.A), content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteFileBlocks 删除文件所有块
func DeleteFileBlocks(ctx context.Context, block bson.A) error {
	mg := ctx.Value("mongo").(*mongo.Client)
	blocks := mg.Database("starfile").Collection("blocks")

	for _, b := range block {
		// 读取当前块
		filter := bson.M{"_id": b}
		var bl bson.M
		blocks.FindOne(ctx, filter).Decode(&bl)

		if bl["level"].(int32) == 1 {
			// 1级块,直接删除
			blocks.DeleteOne(ctx, bl)
		} else {
			// 多级块,递归继续删除
			err := DeleteFileBlocks(ctx, bl["content"].(bson.A))
			if err != nil {
				return err
			}
			blocks.DeleteOne(ctx, bl)
		}
	}
	return nil
}
