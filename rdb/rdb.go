package rdb

import "github.com/go-redis/redis"

var Client *redis.Client

func CreateCon() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
