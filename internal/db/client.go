package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	host     = "localhost"
	user     = "avito_test"
	password = "test"
	dbname   = "avito_test"
	sslmode  = "disable"
	port     = 5432
)

func CreateDB(ctx context.Context) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, GenerateConn())
	if err != nil {
		log.Fatalf("Failed to connect to DB with err: %v", err)
	}
	return newDB(pool), nil
}

func GenerateConn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", host, user, password, dbname, port, sslmode)
}
