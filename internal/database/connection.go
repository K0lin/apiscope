package database

import (
	"APIScope/internal/config"
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitDatabase(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.DatabasePath,
		Password: cfg.DatabasePassword,
		DB:       0,
	})

	_, err := RedisClient.Ping(ctx).Result()
	return err
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func GetContext() context.Context {
	return ctx
}
