package database

import (
	"database/sql"
	"fmt"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://sandeep:admin123@localhost/whatsapp?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}
