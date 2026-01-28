package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	conn, err := sql.Open("sqlite", "books.db")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
