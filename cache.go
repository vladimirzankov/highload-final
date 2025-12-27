package main

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewCache() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	return redis.NewClient(&redis.Options{
		Addr:        addr,
		DialTimeout: 5 * time.Second,
	})
}
