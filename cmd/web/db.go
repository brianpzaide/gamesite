package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func getRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: "",
		DB:       0,
	})

	err := rdb.Ping(context.Background()).Err()

	numOftriesLeft := 5

	for err != nil && numOftriesLeft > 0 {
		time.Sleep(100 * time.Millisecond)
		err = rdb.Ping(context.Background()).Err()
		numOftriesLeft--
	}

	return rdb, err
}
