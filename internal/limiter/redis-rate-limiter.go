package limiter

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/felipezschornack/golang-rate-limit-middleware/config"
	"github.com/go-redis/redis/v8"
)

type RedisRateLimiter struct {
	Config      *config.EnvironmentVariables
	RedisClient *redis.Client
}

func NewRedisRateLimiter() *RedisRateLimiter {
	conf := getConfig()
	redisClient := getRedisClient(conf)
	return &RedisRateLimiter{
		RedisClient: redisClient,
		Config:      conf,
	}
}

func getConfig() *config.EnvironmentVariables {
	conf, err := config.LoadConfig("../")
	if err != nil {
		panic(err)
	}
	return conf
}

func getRedisClient(conf *config.EnvironmentVariables) *redis.Client {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", conf.RedisHostname, conf.RedisPort),
	})
	_ = redisClient.FlushDB(ctx).Err()
	return redisClient
}

func (rrl *RedisRateLimiter) IsBlocked(r *http.Request) bool {
	token := r.Header.Get("API_KEY")

	if len(token) > 0 {
		return rrl.dealWithAccessToken(token) != nil
	} else {
		return rrl.dealWithIpAddress(r) != nil
	}
}

func (rrl *RedisRateLimiter) dealWithAccessToken(token string) error {
	rateLimit := strToInt64(rrl.Config.AccessTokenRateLimitInSeconds)
	blockingWindow := strToInt64(rrl.Config.AccessTokenBlockingWindowInSeconds)
	return rrl.verifyBlock(token, rateLimit, time.Duration(blockingWindow)*time.Second)
}

func (rrl *RedisRateLimiter) dealWithIpAddress(r *http.Request) error {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	ip, _, _ = net.SplitHostPort(ip)

	rateLimit := strToInt64(rrl.Config.IpRateLimitInSeconds)
	blockingWindow := strToInt64(rrl.Config.IpBlockingWindowInSeconds)
	return rrl.verifyBlock(ip, rateLimit, time.Duration(blockingWindow)*time.Second)
}

func strToInt64(in string) int64 {
	out, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		panic(err)
	}
	return out
}

func (rrl *RedisRateLimiter) verifyBlock(key string, limit int64, window time.Duration) error {
	return rrl.RedisClient.Watch(context.Background(), func(tx *redis.Tx) error {
		currentTime := time.Now()
		keyWindow := fmt.Sprintf("%s_%d", key, currentTime.Unix()/int64(window.Seconds()))

		count, err := rrl.RedisClient.Get(context.Background(), keyWindow).Int64()
		if err != nil && err != redis.Nil {
			return err
		}

		if count > limit {
			return errors.New("Rate limit exceeded")
		}

		pipe := rrl.RedisClient.TxPipeline()
		pipe.Incr(context.Background(), keyWindow)
		pipe.Expire(context.Background(), keyWindow, window)
		_, err = pipe.Exec(context.Background())
		if err != nil {
			return err
		}
		return nil
	}, key)
}
