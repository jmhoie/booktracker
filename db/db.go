package db

import (
	"time"
	"database/sql"
	"github.com/jmhoie/booktracker/models"

	_ "modernc.org/sqlite"
)

// DB wraps the database connection
type DB struct {
	conn *sql.DB
}

// createTables creates the database tables: books, authors, book_authors and
// as well as a trigger to clean up orphaned authors
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

// getOrCreateAuthor tries to fetch an author based on their name, if no 
// author is found, it creates a new one and adds it to the authors table
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

// Open creates a new database connection and initializes tables.
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

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// GetAuthors fetches the authors for a book based on the book id
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

// AddBook creates and adds a book to the database
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

// GetBook fetches a book by id, with its authors, from the database
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

// GetAllBooks fetches all books with their authors from the database
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

// UpdateBookStatus sets the status of the book in the database
func (db *DB) UpdateBookStatus(id int, status models.BookStatus) error {
	_, err := db.conn.Exec(
		`UPDATE books SET status = ? WHERE id = ?`,
		status,
		id,
	)
	return err
}

// UpdateStartedAt sets the started_at timestamp in the database
func (db *DB) UpdateStartedAt(id int, startedAt time.Time) error {
	_, err := db.conn.Exec(
		`UPDATE books SET started_at = ? WHERE id = ?`,
		startedAt,
		id,
	)
	return err
}

// UpdateFinishedAt sets the finished_at timestamp in the database
func (db *DB) UpdateFinishedAt(id int, finishedAt time.Time) error {
	_, err := db.conn.Exec(
		`UPDATE books SET started_at = ? WHERE id = ?`,
		finishedAt,
		id,
	)
	return err
}

// DeleteBook deletes book by id from the database
func (db *DB) DeleteBook(id int) error {
	_, err := db.conn.Exec(
		`DELETE FROM books WHERE id = ?`,
		id,
	)
	return err
}

// DeleteAuthor deletes author by id from the database
func (db *DB) DeleteAuthor(id int) error {
	_, err := db.conn.Exec(
		`DELETE FROM authors WHERE id = ?`,
		id,
	)
	return err
}
