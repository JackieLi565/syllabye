package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() (*DB, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv(config.PostgresUser), os.Getenv(config.PostgresPassword), os.Getenv(config.PostgresHost), os.Getenv(config.PostgresPort), os.Getenv(config.PostgresDatabase))

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.TODO())
	if err != nil {
		log.Fatal("failed to create a database connection pool")
	} else {
		log.Println("database connection pool created")
	}
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
