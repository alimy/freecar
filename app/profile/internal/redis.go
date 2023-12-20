package internal

import (
	"fmt"

	"github.com/alimy/freecar/app/profile/config"
	"github.com/alimy/freecar/library/core/consts"
	"github.com/go-redis/redis/v8"
)

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GlobalServerConfig.RedisInfo.Host, config.GlobalServerConfig.RedisInfo.Port),
		Password: config.GlobalServerConfig.RedisInfo.Password,
		DB:       consts.RedisProfileClientDB,
	})
}
