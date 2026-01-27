package models

import (
	"fmt"
	"strings"
	"time"
)

type Book struct {
	ID int
	Title string
	Authors []Author
	ISBN13 string
	Status string
	StartedAt *time.Time
	FinishedAt *time.Time
}

func (b *Book) String() string {
	var authorNames []string
	for _, author := range b.Authors {
		authorNames = append(authorNames, author.Name)
	}
	authors := strings.Join(authorNames, ", ")

	return fmt.Sprintf("%s by %s (%s)", b.Title, authors, b.Status)
}

func NewBook(title string, authors []Author, isbn13 string) *Book {
	return &Book{
		Title: title,
		Authors: authors,
		ISBN13: isbn13,
		Status: "TBR",
		StartedAt: nil,
		FinishedAt: nil,
	}
}
