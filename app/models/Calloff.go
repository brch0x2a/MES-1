package models

import (
	"database/sql"
)

type Calloff_holder struct {
	Profile_picture		 string
	Id           int
	Linea        string
	NickName     string
	Fname        string
	Lname        string
	RequestDate  string
	CloseDate    sql.NullString //para poder hacer el query de un null sin que se caiga por null
	CodMaterial  int
	MaterialName string
	Amount       int
	Comment      string
	State        string
	StateColor   string
}
