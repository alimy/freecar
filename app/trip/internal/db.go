package internal

import (
	"context"
	"fmt"

	"github.com/alimy/freecar/app/trip/conf"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB() *mongo.Database {
	c := conf.GlobalServerConfig.MongoDBInfo
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(
		fmt.Sprintf(consts.MongoURI, c.User, c.Password, c.Host, c.Port)))
	if err != nil {
		klog.Fatal("cannot connect mongodb", err)
	}
	return mongoClient.Database(c.Name)
}
