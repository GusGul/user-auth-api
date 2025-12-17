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

// Register receives encrypted password, decrypts it, hashes it, and stores the user
func (uc *UserUseCase) Register(ctx context.Context, email, encryptedPasswordBase64 string) error {
	// 1. Decrypt Password
	// Note: We expect the input password to be encrypted with the public key
	// If the user sends raw password (for testing without frontend RSA), we might need a flag or separate check.
	// But strictly following requirements: "receber a senha do usuário criptografado".

	// However, usually register receives the raw password over TLS (HTTPS).
	// If the requirement is strict "encrypted with public key", we decrypt here.
	// Assuming the client sends base64 encoded RSA encrypted string.
	plainPassword, err := uc.rsaDecrypter.Decrypt([]byte(encryptedPasswordBase64))
	if err != nil {
		return errors.New("failed to decrypt password: " + err.Error())
	}

	// 2. Hash Password
	hashedPassword, err := security.HashPassword(plainPassword)
	if err != nil {
		return err
	}

	// 3. Create User
	user := &domain.User{
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return uc.userRepo.Create(ctx, user)
}

func (uc *UserUseCase) Login(ctx context.Context, email, encryptedPasswordBase64 string) (string, error) {
	// 1. Get User
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 2. Decrypt Incoming Password
	plainPassword, err := uc.rsaDecrypter.Decrypt([]byte(encryptedPasswordBase64))
	if err != nil {
		return "", errors.New("failed to decrypt password")
	}

	// 3. Verify Password
	if !security.CheckPasswordHash(plainPassword, user.Password) {
		return "", errors.New("invalid credentials")
	}

	// 4. Generate Token
	token, err := uc.jwtProvider.GenerateToken(email, 24*time.Hour) // Using email or ID as subject
	if err != nil {
		return "", err
	}

	// 5. Cache Session (Optional but fulfilling Redis requirement)
	// We can store the token or just a flag that user is logged in
	_ = uc.cacheRepo.SetUserSession(ctx, email, token, 3600*24) // 24h

	return token, nil
}
