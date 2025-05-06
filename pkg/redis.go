package pkg

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func RedisConnect() *redis.Client {
	redisHost := os.Getenv("RDSHOST")
	redisPort := os.Getenv("RDSPORT")
	return redis.NewClient(&redis.Options{Addr: redisHost + ":" + redisPort})
}
