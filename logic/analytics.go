package logic

import (
	"github.com/go-redis/redis"
)

func AddVisit() int64 {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // 指定
		Password: "",
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	num, _ := redisdb.Incr("visit").Result()
	return num
}

func QueryVisit() string {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // 指定
		Password: "",
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	num, _ := redisdb.Get("visit").Result()
	return num
}
