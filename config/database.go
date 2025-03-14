package config

import (
	"context"
	"fmt"
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

	sslCertPath := "certs/server.crt"

	return pgxpool.New(context.Background(), fmt.Sprintf(dbConnect, sslCertPath))
}
