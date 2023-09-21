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

type Storage interface {
	RunMigration(*MigrationBody) error
}
type PostgresStore struct {
	db *sql.DB
}

type MigrationBody struct {
	query string
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

	fmt.Println("Connected to DB!")
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) RunMigration(m *MigrationBody) error {

	_, err := s.db.Exec(m.query)
	if err != nil {
		// errors.New(FmtRed("Error trying to run query => ") + err.Error())
		log.Fatal(FmtRed("Error trying to run query => ") + err.Error())
		return err
	}
	fmt.Println(FmtGreen("Migration done!"))
	return nil
}
