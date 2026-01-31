package models

import "fmt"

type Author struct {
	Id int
	Name string
}

func (a *Author) String() string {
	return fmt.Sprintf("%s", a.Name)
}
