package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDb struct {
	Pool *pgxpool.Pool
}

func NewPostgresDb() (*PostgresDb, error) {
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
	return &PostgresDb{Pool: pool}, nil
}

func (db *PostgresDb) Close() {
	db.Pool.Close()
}
