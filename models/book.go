package models

import (
	"fmt"
	"strings"
	"time"
)

type BookStatus string

const (
	StatusTBR BookStatus = "TBR"
	StatusReading BookStatus = "Reading"
	StatusFinished BookStatus = "Finished"
	StatusDNF BookStatus = "DNF"
)

type Book struct {
	Id int
	Title string
	Authors []Author
	Isbn13 string
	Status BookStatus
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
		Isbn13: isbn13,
		Status: StatusTBR,
		StartedAt: nil,
		FinishedAt: nil,
	}
}
