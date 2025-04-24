package pkg

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PG *pgxpool.Pool

func Connect() (*pgxpool.Pool, error) {
	// CONNECT TO DB
	dbEnv := []any{}
	dbEnv = append(dbEnv, os.Getenv("DBUSER"))
	dbEnv = append(dbEnv, os.Getenv("DBPASS"))
	dbEnv = append(dbEnv, os.Getenv("DBHOST"))
	dbEnv = append(dbEnv, os.Getenv("DBPORT"))
	dbEnv = append(dbEnv, os.Getenv("DBNAME"))

	dbString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbEnv...)
	var err error
	PG, err = pgxpool.New(context.Background(), dbString)
	if err != nil {
		return nil, err
	}
	err = PG.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	log.Println("DB connected")
	return PG, nil
}
