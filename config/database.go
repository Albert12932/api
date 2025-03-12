package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func ConnectDB() (*pgxpool.Pool, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("Error while loading env")
		os.Exit(1)
	}

	dbConnect := os.Getenv("DATABASE_CONNECT")

	return pgxpool.New(context.Background(), dbConnect)
}
