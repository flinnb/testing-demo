package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var pool *sql.DB

func Connect() (err error) {
	connStr := fmt.Sprintf(
		"user=%s password='%s' host=%s port=%s dbname=%s sslmode=%s connect_timeout=%d application_name='%s'",
		"connectrn",
		"canary",
		os.Getenv("DB_HOST"),
		"5432",
		os.Getenv("DB_NAME"),
		"disable",
		5,
		"connectrn-eval",
	)
	pool, err = sql.Open("postgres", connStr)
	return
}

func Pool() *sql.DB {
	return pool
}
