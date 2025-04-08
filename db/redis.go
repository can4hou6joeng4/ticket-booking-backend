package db

import (
	"github.com/can4hou6joeng4/ticket-booking-project-v1/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(config *config.EnvConfig) *redis.Client {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.RedisHost + ":" + config.RedisConfig.RedisPort,
		Password: config.RedisConfig.RedisPassword,
		DB:       config.RedisConfig.RedisDB,
	})
	return RedisClient
}
