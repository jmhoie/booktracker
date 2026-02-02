package db

import (
	"time"
	"github.com/jmhoie/booktracker/models"

	_ "modernc.org/sqlite"
)

// ResetAndSeed deletes all data and populates with mock data (for testing)
func (db *DB) ResetAndSeed() error {
	_, err := db.conn.Exec(`DELETE FROM books`)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(`DELETE from authors`)
	if err != nil {
		return err
	}

	// reset autoincrement
	_, err = db.conn.Exec(`DELETE FROM sqlite_sequence`)
	if err != nil {
		return nil
	}

	now := time.Now()
    lastWeek := now.AddDate(0, 0, -7)
    twoWeeksAgo := now.AddDate(0, 0, -14)
    lastMonth := now.AddDate(0, -1, 0)
    twoMonthsAgo := now.AddDate(0, -2, 0)
    threeMonthsAgo := now.AddDate(0, -3, 0)
    
    mockBooks := []*models.Book{
        {
            Title:      "The Go Programming Language",
            Isbn13:     "9780134190440",
            Status:     models.StatusFinished,
            StartedAt:  &twoMonthsAgo,
            FinishedAt: &lastMonth,
            Authors: []models.Author{
                {Name: "Alan Donovan"},
                {Name: "Brian Kernighan"},
            },
        },
        {
            Title:     "1984",
            Isbn13:    "9780451524935",
            Status:    models.StatusReading,
            StartedAt: &lastWeek,
            Authors: []models.Author{
                {Name: "George Orwell"},
            },
        },
        {
            Title:  "Dune",
            Isbn13: "9780441172719",
            Status: models.StatusTBR,
            Authors: []models.Author{
                {Name: "Frank Herbert"},
            },
        },
        {
            Title:  "The Pragmatic Programmer",
            Isbn13: "9780135957059",
            Status: models.StatusTBR,
            Authors: []models.Author{
                {Name: "David Thomas"},
                {Name: "Andrew Hunt"},
            },
        },
        {
            Title:      "Clean Code",
            Isbn13:     "9780132350884",
            Status:     models.StatusFinished,
            StartedAt:  &threeMonthsAgo,
            FinishedAt: &twoMonthsAgo,
            Authors: []models.Author{
                {Name: "Robert C. Martin"},
            },
        },
        {
            Title:     "The Hobbit",
            Isbn13:    "9780547928227",
            Status:    models.StatusReading,
            StartedAt: &twoWeeksAgo,
            Authors: []models.Author{
                {Name: "J.R.R. Tolkien"},
            },
        },
        {
            Title:  "Sapiens",
            Isbn13: "9780062316097",
            Status: models.StatusTBR,
            Authors: []models.Author{
                {Name: "Yuval Noah Harari"},
            },
        },
        {
            Title:      "The Mythical Man-Month",
            Isbn13:     "9780201835953",
            Status:     models.StatusFinished,
            StartedAt:  &threeMonthsAgo,
            FinishedAt: &twoMonthsAgo,
            Authors: []models.Author{
                {Name: "Frederick P. Brooks Jr."},
            },
        },
        {
            Title:  "Project Hail Mary",
            Isbn13: "9780593135204",
            Status: models.StatusTBR,
            Authors: []models.Author{
                {Name: "Andy Weir"},
            },
        },
        {
            Title:      "Designing Data-Intensive Applications",
            Isbn13:     "9781449373320",
            Status:     models.StatusFinished,
            StartedAt:  &threeMonthsAgo,
            FinishedAt: &lastMonth,
            Authors: []models.Author{
                {Name: "Martin Kleppmann"},
            },
        },
    }

	for _, book := range mockBooks {
		if err := db.AddBook(book); err != nil {
			return err
		}
	}

	return nil
}
