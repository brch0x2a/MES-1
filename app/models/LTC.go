package models

type LTC struct {
	Id          int
	Description string
}

//

type LTCLine struct {
	Id     int
	Line   string
	LTC    string //Line time classification
	Sub    string //Sub classification
	Branch string
	Event  string
	Code   int
	Color  string
}
