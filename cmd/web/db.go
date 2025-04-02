package main

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func getRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv(*redisAddr),
		Password: "",
		DB:       0,
	})

	return rdb
}
