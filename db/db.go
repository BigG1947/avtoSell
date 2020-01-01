package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "avtosell.db")
	if err != nil {
		log.Printf("Error in connection to database: %s\n", err)
	}
	return db, err
}