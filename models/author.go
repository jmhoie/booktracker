package models

import "fmt"

type Author struct {
	ID int
	Name string
}

func (a *Author) String() string {
	return fmt.Sprintf("%s", a.Name)
}
