package config

import (
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
)

var PostgresConfig pgconn.Config

func init() {
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		port = 5432
	}

	PostgresConfig = pgconn.Config{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DATABASE"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     uint16(port),
	}
}
