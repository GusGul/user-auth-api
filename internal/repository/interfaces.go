package repository

import (
	"context"
	"user-auth-api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type CacheRepository interface {
	SetUserSession(ctx context.Context, userID string, token string, ttl int) error
	GetSession(ctx context.Context, userID string) (string, error)
}
