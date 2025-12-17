package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"net/http"
	"os"

	"user-auth-api/config"
	"user-auth-api/internal/infra/database"
	internalHttp "user-auth-api/internal/infra/http"
	"user-auth-api/internal/infra/http/handlers"
	"user-auth-api/internal/infra/security"
	"user-auth-api/internal/usecase"
)

func main() {
	cfg := config.LoadConfig()

	if err := ensureKeysExist(); err != nil {
		log.Fatalf("Failed to setup keys: %v", err)
	}

	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	userRepo := database.NewMySQLUserRepository(db)
	cacheRepo := database.NewRedisCacheRepository(redisClient)

	rsaDecrypter, err := security.NewRSADecrypter("private.pem")
	if err != nil {
		log.Fatalf("Failed to load RSA private key: %v", err)
	}

	jwtProvider := security.NewJWTProvider(cfg.JWTSecret)

	userUC := usecase.NewUserUseCase(userRepo, cacheRepo, rsaDecrypter, jwtProvider)
	userHandler := handlers.NewUserHandler(userUC)

	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.HandleFunc("GET /public-key", func(w http.ResponseWriter, r *http.Request) {
		key, _ := os.ReadFile("public.pem")
		w.Write(key)
	})

	// Protected Routes
	mux.HandleFunc("GET /profile", internalHttp.AuthMiddleware(jwtProvider, func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id")
		w.Write([]byte("Hello user " + userID.(string)))
	}))

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := http.ListenAndServe(cfg.AppPort, internalHttp.LoggerMiddleware(mux)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func ensureKeysExist() error {
	if _, err := os.Stat("private.pem"); err == nil {
		return nil
	}

	log.Println("Generating RSA keys for demo...")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	if err := os.WriteFile("private.pem", privPEM, 0600); err != nil {
		return err
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	if err := os.WriteFile("public.pem", pubPEM, 0644); err != nil {
		return err
	}

	return nil
}
