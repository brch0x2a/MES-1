package models

type CodingVerification struct {
	Id                int
	Date_reg          string
	OFname            string
	OAProfile_Picture string
	OLname            string

	Line         string
	Presentation string
	Photo_coding string

	Type string

	Comment string

	Date_check        string
	CFname            string
	CAProfile_Picture string
	CLname            string

	State string
}
