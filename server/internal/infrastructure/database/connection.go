package database

import (
	"context"
	"fmt"
	"log"

	"github.com/JackieLi565/syllabye/server/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() (*DB, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.PostgresConfig.User, config.PostgresConfig.Password, config.PostgresConfig.Host, config.PostgresConfig.Port, config.PostgresConfig.Database)

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	log.Println("database connection pool created")
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
