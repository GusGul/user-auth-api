package database

import (
	"context"
	"time"

	"user-auth-api/internal/repository"

	"github.com/redis/go-redis/v9"
)

type RedisCacheRepository struct {
	client *redis.Client
}

func NewRedisCacheRepository(client *redis.Client) repository.CacheRepository {
	return &RedisCacheRepository{client: client}
}

func (r *RedisCacheRepository) SetUserSession(ctx context.Context, userID string, token string, ttl int) error {
	return r.client.Set(ctx, "session:"+userID, token, time.Duration(ttl)*time.Second).Err()
}

func (r *RedisCacheRepository) GetSession(ctx context.Context, userID string) (string, error) {
	return r.client.Get(ctx, "session:"+userID).Result()
}
