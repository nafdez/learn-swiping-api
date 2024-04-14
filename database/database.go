package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrConnFailed = errors.New("connection failed")
)

// Attempts to connect to the database.
//
// Returns a connection pool and an error.
func Connect() (*sql.DB, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	url := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	_ = mysql.Config{}

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true", user, pass, url, dbName))
	if err != nil {
		return nil, ErrConnFailed
	}

	if err = db.Ping(); err != nil {
		return nil, ErrConnFailed
	}

	return db, nil
}
