package limiter

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/felipezschornack/golang-rate-limit-middleware/config"
	"github.com/felipezschornack/golang-rate-limit-middleware/internal/db"
	"github.com/go-redis/redis/v8"
)

type RedisRateLimiter struct {
	Config *config.EnvironmentVariables
	Redis  *db.Redis
}

func NewRedisRateLimiter(conf *config.EnvironmentVariables, redisClient *db.Redis) *RedisRateLimiter {
	return &RedisRateLimiter{
		Redis:  redisClient,
		Config: conf,
	}
}

func (rrl *RedisRateLimiter) IsBlocked(token string, ipAddress string) bool {
	if len(token) > 0 {
		return rrl.dealWithAccessToken(token) != nil
	} else {
		return rrl.dealWithIpAddress(ipAddress) != nil
	}
}

func (rrl *RedisRateLimiter) dealWithAccessToken(token string) error {
	rateLimit := strToInt64(rrl.Config.AccessTokenRateLimitInSeconds)
	blockingWindow := strToInt64(rrl.Config.AccessTokenBlockingWindowInSeconds)
	return rrl.verifyBlock(token, rateLimit, time.Duration(blockingWindow)*time.Second)
}

func (rrl *RedisRateLimiter) dealWithIpAddress(ipAddress string) error {
	rateLimit := strToInt64(rrl.Config.IpRateLimitInSeconds)
	blockingWindow := strToInt64(rrl.Config.IpBlockingWindowInSeconds)
	return rrl.verifyBlock(ipAddress, rateLimit, time.Duration(blockingWindow)*time.Second)
}

func strToInt64(in string) int64 {
	out, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		panic(err)
	}
	return out
}

func (rrl *RedisRateLimiter) verifyBlock(key string, limit int64, window time.Duration) error {
	return rrl.Redis.Client.Watch(context.Background(), func(tx *redis.Tx) error {
		currentTime := time.Now()
		keyWindow := fmt.Sprintf("%s_%d", key, currentTime.Unix()/int64(window.Seconds()))

		count, err := rrl.Redis.Client.Get(context.Background(), keyWindow).Int64()
		if err != nil && err != redis.Nil {
			return err
		}

		if count < limit {
			pipe := rrl.Redis.Client.TxPipeline()
			pipe.Incr(context.Background(), keyWindow)
			pipe.Expire(context.Background(), keyWindow, window)
			_, err = pipe.Exec(context.Background())
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("Rate limit exceeded")

	}, key)
}
