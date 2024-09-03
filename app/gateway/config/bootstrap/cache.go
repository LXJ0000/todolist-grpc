package bootstrap

import (
	"log"
	"time"

	"context"

	"github.com/LXJ0000/todolist-grpc-gateway/pkg/cache"
	"github.com/redis/go-redis/v9"
)

func NewRedisCache(env *Env) cache.RedisCache {
	cmd := redis.NewClient(&redis.Options{
		Addr:     env.RedisAddr,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})
	if _, err := cmd.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}
	return cache.NewRedisCache(cmd, time.Duration(env.RedisExpiration)*time.Minute)
}
