package internal

import (
	"fmt"

	"github.com/alimy/freecar/app/profile/conf"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/go-redis/redis/v8"
)

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.GlobalServerConfig.RedisInfo.Host, conf.GlobalServerConfig.RedisInfo.Port),
		Password: conf.GlobalServerConfig.RedisInfo.Password,
		DB:       consts.RedisProfileClientDB,
	})
}
