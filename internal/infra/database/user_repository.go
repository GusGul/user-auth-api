package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"user-auth-api/internal/domain"
	"user-auth-api/internal/repository"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) repository.UserRepository {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
	}

	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (email, password, created_at, updated_at) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *MySQLUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := "SELECT id, email, password, created_at, updated_at FROM users WHERE email = ?"
	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
