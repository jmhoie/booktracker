package db

import (
	"database/sql"
	"github.com/jmhoie/booktracker/models"

	_ "modernc.org/sqlite"
)

// DB wraps the database connection.
type DB struct {
	conn *sql.DB
}

// Creates a new database connection and initializes tables.
func Open() (*DB, error) {
	conn, err := sql.Open("sqlite", "books.db")
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	db := &DB{conn: conn}

	if err = db.createTables(); err != nil {
		return nil, err
	}

	return db, nil
}

// Closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		isbn13 TEXT,
		status TEXT NOT NULL DEFAULT 'TBR',
		started_at DATETIME,
		finished_at DATETIME,
		created_at DATETAME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS book_authors (
		book_id INTEGER NOT NULL,
		author_id INTEGER NOT NULL,
		PRIMARY KEY (book_id, author_id),
		FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
		FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE
	);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// gets book by ID with its authors
func (db *DB) GetBook(id int) (*models.Book, error) {
	book := &models.Book{}

	err := db.conn.QueryRow(
		`SELECT id, title, isbn13, status, started_at, finished_at
		FROM books WHERE id = ?`,
		id,
	).Scan(
		&book.Id,
		&book.Title,
		&book.Isbn13,
		&book.Status,
		&book.StartedAt,
		&book.FinishedAt,
	)
	if err != nil {
		return nil, err
	}

	// get authors
	rows, err := db.conn.Query(
		`SELECT a.id, a.name
		FROM authors a
		JOIN book_authors ba ON a.id = ba.author_id
		WHERE ba.book_id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var author models.Author
		if err := rows.Scan(&author.Id, &author.Name); err != nil {
			return nil, err
		}
		book.Authors = append(book.Authors, author)
	}

	return book, rows.Err()
}
