package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func SetUserSession(redis *redis.Client, ctx context.Context, userId uint, token string, role string) error {
	key := fmt.Sprintf("user:%d:session", userId)
	value := map[string]interface{}{
		"token": token,
		"role":  role,
	}
	return redis.HMSet(ctx, key, value).Err()
}

func GetUserSession(redis *redis.Client, ctx context.Context, userId uint) (map[string]string, error) {
	key := fmt.Sprintf("user:%d:session", userId)
	return redis.HGetAll(ctx, key).Result()
}

func DeleteUserSession(redis *redis.Client, ctx context.Context, userId uint) error {
	key := fmt.Sprintf("user:%d:session", userId)
	return redis.Del(ctx, key).Err()
}

func SetUserPermissions(redis *redis.Client, ctx context.Context, userId uint, permissions []string) error {
	key := fmt.Sprintf("user:%d:permissions", userId)
	return redis.SAdd(ctx, key, permissions).Err()
}

func GetUserPermissions(redis *redis.Client, ctx context.Context, userId uint) ([]string, error) {
	key := fmt.Sprintf("user:%d:permissions", userId)
	return redis.SMembers(ctx, key).Result()
}

func DeleteUserPermissions(redis *redis.Client, ctx context.Context, userId uint) error {
	key := fmt.Sprintf("user:%d:permissions", userId)
	return redis.Del(ctx, key).Err()
}

func SetExpiration(redis *redis.Client, ctx context.Context, key string, expiration time.Duration) error {
	return redis.Expire(ctx, key, expiration).Err()
}
