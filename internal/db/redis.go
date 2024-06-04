package db

import (
	"context"
	"fmt"

	"github.com/felipezschornack/golang-rate-limit-middleware/config"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

func GetRedisClient(conf *config.EnvironmentVariables) *Redis {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", conf.RedisHostname, conf.RedisPort),
	})
	_ = redisClient.FlushDB(ctx).Err()
	return &Redis{Client: redisClient}
}

func (r *Redis) ClearDatabase() {
	r.Client.FlushAll(context.Background()).Err()
}
