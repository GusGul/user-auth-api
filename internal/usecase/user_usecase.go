package usecase

import (
	"context"
	"errors"
	"time"

	"user-auth-api/internal/domain"
	"user-auth-api/internal/infra/security"
	"user-auth-api/internal/repository"
)

type UserUseCase struct {
	userRepo     repository.UserRepository
	cacheRepo    repository.CacheRepository
	rsaDecrypter *security.RSADecrypter
	jwtProvider  *security.JWTProvider
}

func NewUserUseCase(ur repository.UserRepository, cr repository.CacheRepository, rsa *security.RSADecrypter, jwt *security.JWTProvider) *UserUseCase {
	return &UserUseCase{
		userRepo:     ur,
		cacheRepo:    cr,
		rsaDecrypter: rsa,
		jwtProvider:  jwt,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, email, encryptedPasswordBase64 string) error {
	plainPassword, err := uc.rsaDecrypter.Decrypt([]byte(encryptedPasswordBase64))
	if err != nil {
		return errors.New("failed to decrypt password: " + err.Error())
	}

	hashedPassword, err := security.HashPassword(plainPassword)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return uc.userRepo.Create(ctx, user)
}

func (uc *UserUseCase) Login(ctx context.Context, email, encryptedPasswordBase64 string) (string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	plainPassword, err := uc.rsaDecrypter.Decrypt([]byte(encryptedPasswordBase64))
	if err != nil {
		return "", errors.New("failed to decrypt password")
	}

	if !security.CheckPasswordHash(plainPassword, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := uc.jwtProvider.GenerateToken(email, 24*time.Hour)
	if err != nil {
		return "", err
	}

	_ = uc.cacheRepo.SetUserSession(ctx, email, token, 3600*24)

	return token, nil
}
