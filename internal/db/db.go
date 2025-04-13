package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() (*pgxpool.Pool, error) {
	url := os.Getenv("DB_URL")
	if url == "" {
		log.Fatal("❌ DB_URL not set in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	DB, err = pgxpool.New(ctx, url)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("✅ Connected to fern-reporter database")
	return DB, nil
}
