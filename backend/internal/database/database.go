package database

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Open reads configuration from the process environment; no secrets are committed.
func Open(ctx context.Context) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = env("DB_ADDR", "127.0.0.1:3306")
	cfg.User = env("DB_USER", "root")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.DBName = env("DB_NAME", "mirmilowo")
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	cfg.Collation = "utf8mb4_unicode_ci"
	cfg.Timeout = 5 * time.Second
	cfg.ReadTimeout = 5 * time.Second
	cfg.WriteTimeout = 5 * time.Second
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(3 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
