package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func ConnectDB() (*pgxpool.Pool, error) {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	cert := os.Getenv("DB_CERT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=verify-ca&sslrootcert=%s", user, password, host, port, dbname)
postgres: //albertt1001:albertt1001@51.250.48.59:5432/mydatabase?sslmode=verify-ca&sslrootcert=%s
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	return pool, nil

}
