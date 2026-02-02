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

// create tables: books, authors, book_authors.
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

	-- Automatically delete authors who are no longer referenced by a book
	CREATE TRIGGER IF NOT EXISTS cleanup_orphaned_authors
	AFTER DELETE ON book_authors
	BEGIN
		DELETE FROM authors
		WHERE id = OLD.author_id
		AND id NOT IN (SELECT DISTINCT author_id FROM book_authors);
	END;
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) getOrCreateAuthor(tx *sql.Tx, name string) (int, error) {
	var id int
	err := tx.QueryRow(`SELECT id FROM authors WHERE name = ?`, name).Scan(&id)
	if err == nil {
		return id, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	// author doesn't exists, create it
	result, err := tx.Exec(`INSERT INTO authors (name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}

	authorId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(authorId), nil
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

	// enable foreign keys (needed for CASCADE to work)
	_, err = conn.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
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

// get the authors for a book based on book id
func (db *DB) GetAuthors(bookId int) ([]models.Author, error) {
	var authors []models.Author	

	rows, err := db.conn.Query(
		`SELECT a.id, a.name
		FROM authors a
		JOIN book_authors ba ON a.id = ba.author_id
		WHERE book_id = ?`,
		bookId,
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
		authors = append(authors, author)
	}

	return authors, nil
}

// adds book to the database
func (db *DB) AddBook(b *models.Book) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		`INSERT INTO books (title, isbn13, status, started_at, finished_at)
		VALUES (?, ?, ?, ?, ?)`,
		b.Title,
		b.Isbn13,
		b.Status,
		b.StartedAt,
		b.FinishedAt,
	)
	if err != nil {
		return err
	}

	bookId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	b.Id = int(bookId)

	for i := range b.Authors {
		authorId, err := db.getOrCreateAuthor(tx, b.Authors[i].Name)
		if err != nil {
			return err
		}
		b.Authors[i].Id = authorId

		_, err = tx.Exec(
			`INSERT INTO book_authors (book_id, author_id) VALUES (?, ?)`,
			bookId,
			authorId,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// get book by id with its authors
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

	authors, err := db.GetAuthors(book.Id)
	if err != nil {
		return nil, err
	}
	book.Authors = authors

	return book, nil
}

// get all books with their authors
func (db *DB) GetAllBooks() ([]models.Book, error) {
	rows, err := db.conn.Query(
		`SELECT id, title, isbn13, status, started_at, finished_at FROM books`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.Id,
			&book.Title,
			&book.Isbn13,
			&book.Status,
			&book.StartedAt,
			&book.FinishedAt,
		); err != nil {
			return nil, err
		}

		authors, err := db.GetAuthors(book.Id)
		if err != nil {
			return nil, err
		}

		book.Authors = authors
		books = append(books, book)
	}

	return books, rows.Err()
}

func (db *DB) UpdateBookStatus(id int, status models.BookStatus) error {
	_, err := db.conn.Exec(
		`UPDATE books SET status = ? WHERE id = ?`,
		status,
		id,
	)
	return err
}

