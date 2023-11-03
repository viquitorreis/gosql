package main

import (
	"database/sql"
	"errors"
	"fmt"
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

func NewPostgresStore() (*PostgresStore, error) { // -------------------- TRATAR MELHOR ERROS
	err := godotenv.Load(".env") // PUXAR A FUNÇÃO CHECKUSERDBCONFIG
	if err != nil {
		return nil, errors.New(FmtRed("Error trying to read .env file") + err.Error())
	}
	connStr := os.Getenv("CONN_STR")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.New(FmtRed("Error trying to stablish database connection") + err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, errors.New(FmtRed("Error, connection to database may no be online") + err.Error())
	}

	fmt.Println(FmtGreen("Connected to DB!"))
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) RunMigration(m *MigrationBody) error {

	_, err := s.db.Exec(m.query)
	if err != nil {
		return errors.New(FmtRed("Error trying to run query => ") + err.Error())
	}
	fmt.Println(FmtGreen("Migration done!"))
	return nil

}
