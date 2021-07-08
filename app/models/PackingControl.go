package models

type PackingControl struct {
	Id           int
	Line         string
	PPhoto       string
	Presentation string
	Date         string
	Profile      string
	Fname        string
	Lname        string
	Lote         int
	Autoclave    int
	Pallet       int
	Box          int
	Observation  string
}
