package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewDBConnect(user, password, dbname string, port int) (*PostgresDB, error) {
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable", user, password, dbname, port)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		dbErr := fmt.Errorf("db connection error: %w", err)
		return nil, dbErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db error: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresDB{DB: db}, nil
}

func (pq *PostgresDB) Close() error {
	err := pq.DB.Close()
	if err != nil {
		return fmt.Errorf("close db error: %w", err)
	}
	return nil
}
