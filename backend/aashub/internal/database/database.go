package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() (*sql.DB, error) {
	// Construct the Data Source Name (DSN)
	dsn := "user:password@tcp(mariadb)/aashub?parseTime=true"

	// Open a DB connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the DB connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
