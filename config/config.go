package config

import (
	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ServerPort  string `env:"SERVER_PORT,required"`
	DBConfig    DBConfig
	RedisConfig RedisConfig
}

type DBConfig struct {
	DBHost     string `env:"DB_HOST,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSSLMode  string `env:"DB_SSLMODE,required"`
}

type RedisConfig struct {
	RedisHost     string `env:"REDIS_HOST,required"`
	RedisPort     string `env:"REDIS_PORT,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	RedisDB       int    `env:"REDIS_DB,required"`
}

func NewEnvConfig() *EnvConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load the .env: %v", err)
	}

	config := &EnvConfig{}
	if err = env.Parse(config); err != nil {
		log.Fatal("Unable to parse config: %v", err)
	}

	dbConfig := &DBConfig{}
	redisConfig := &RedisConfig{}
	if err = env.Parse(dbConfig); err != nil {
		log.Fatal("Unable to parse DB config: %v", err)
	}
	if err = env.Parse(redisConfig); err != nil {
		log.Fatal("Unable to parse Redis config: %v", err)
	}

	config.DBConfig = *dbConfig
	config.RedisConfig = *redisConfig

	return config
}
