package main

import (
	"fmt"

	"github.com/jmhoie/booktracker/models"
)

func main() {
	book := models.NewBook(
		"The Go Programming Language",
		[]models.Author{{Name: "Alan Donovan"}, {Name: "Brian Kernighan"}}, 	
		"9780134190440",
	)
	fmt.Println(book)
}
