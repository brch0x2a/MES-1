package models

type SalsitasOutputControl struct {
	Id int

	Line         string
	Presentation string
	Output_date  string
	Batch_init   string
	Batch_end    string

	MiProfilePicture string
	MiFname          string
	MiLname          string

	MaProfilePicture string
	MaFname          string
	MaLname          string

	Observations string
}
