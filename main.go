package main

import (
	"fmt"
	"github.com/jmhoie/booktracker/db"
)

func main() {
	db, _ := db.Open()
	defer db.Close()

	_ = db.ResetAndSeed()

	books, _ := db.GetAllBooks()
	fmt.Println(books)
}
