package main

import (
//	"github.com/jmhoie/booktracker/models"
	"github.com/jmhoie/booktracker/cmd"
)

func main() {
	// book := models.NewBook(
	// 	"The Go Programming Language",
	// 	[]models.Author{{Name: "Alan Donovan"}, {Name: "Brian Kernighan"}}, 	
	// 	"9780134190440",
	// )

	cmd.Run()
}
