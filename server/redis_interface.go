package main

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func ConnectRedis() {
	password := os.Getenv("REDIS_PASSWORD")

	if password == "" {
		panic("missing REDIS_PASSWORD environment variable")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-11110.c257.us-east-1-3.ec2.redns.redis-cloud.com:11110",
		Username: "default",
		Password: password,
		DB:       0,
	})
}

func LIndex(key string, index int) string {
	ctx := context.Background()

	value, err := rdb.LIndex(ctx, key, int64(index)).Result()

	if err != nil {
		panic(err) //TODO replace
	}

	return value
}
