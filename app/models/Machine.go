package models

type Machine struct {
	Id   int
	Name string
}

type Component struct {
	Id          int
	Name        string
	Description string
	Photo       string
	Id_machine  int
}

type Component_holder struct {
	Id          int
	Name        string
	Description string
	Photo       string
	Machine     string
}
