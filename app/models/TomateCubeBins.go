package models

type TomateCubeBins struct {
	Id          int
	Line        string
	Date        string
	Aperture    string
	Consumption string

	Exposure string

	Profile string
	Fname   string
	Lname   string
	Lote    int
	Weight  float32

	Sticker bool

	Observation string
}
