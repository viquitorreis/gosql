package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("An error occurred. Err: %s", err)
	}
	connStr := os.Getenv("CONN_STR")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.New(FmtRed("Error") + err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, errors.New(FmtRed("Error") + err.Error())
	}

	fmt.Println("connected")
	return &PostgresStore{
		db: db,
	}, nil
}
