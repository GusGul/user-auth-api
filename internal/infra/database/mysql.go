package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"user-auth-api/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Retry loop for Docker startup race conditions
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to MySQL")
			return db, nil
		}
		log.Println("Waiting for MySQL to be ready...")
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to MySQL after retries: %v", err)
}
