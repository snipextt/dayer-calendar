package storage

import "github.com/redis/go-redis/v9"

var rclient *redis.Client

func connectToRedis() {
	opt, _ := redis.ParseURL("rediss://default:c8e2823af9604f0f82d0e3959a4a5024@apn1-enabled-tapir-34942.upstash.io:34942")
	rclient = redis.NewClient(opt)
}

func GetRedisInstance() *redis.Client {
	return rclient
}
