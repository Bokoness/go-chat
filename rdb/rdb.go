package rdb

import "github.com/go-redis/redis"

func CreateCon() *redis.Client {
	red := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return red
}
