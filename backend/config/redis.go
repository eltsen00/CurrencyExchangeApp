package config

import (
	"eltsen00/CurrencyExchangeApp/backend/global"
	"log"
	"os"

	"github.com/go-redis/redis"
)

func InitRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v, system don't start", err)
		os.Exit(1)
	}
	global.RedisDB = RedisClient
}
