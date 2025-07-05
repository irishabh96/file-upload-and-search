package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func initQueue() {
	rdb = redis.NewClient(&redis.Options{
		Addr: RedisURI,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Could not connect to Redis:", err)
		rdb = nil
	} else {
		fmt.Println("Connected to Redis!")
	}
}
