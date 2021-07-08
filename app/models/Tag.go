package models

import "database/sql"

type Tag_holder struct {
	Profile_picture string
	Id              int
	Line            string
	User            string
	Fname           string
	Lname           string
	RequestDate     string
	CloseDate       sql.NullString
	Type            string
	Color           string
	Priority        string
	Equipment       string
	Anomaly         string
	Qa              string
	Cost            string
	Product         string
	Mortal          string
	Deliver         string
	Safety          string
	Affect          string
	Before          string
	AUser           string
	AFname          string
	ALname          string
	Improvement     string
	State           string
	SColor          string
}

type QaTag_holder struct {
	Profile_picture string
	Id              int
	Line            string
	User            string
	Fname           string
	Lname           string
	RequestDate     string
	CloseDate       sql.NullString
	Type            string
	Color           string
	Priority        string
	Equipment       string
	Anomaly         string
	Improvement     string
	State           string
	SColor          string
}

type TagCounter struct {
	Blue   int
	Red    int
	Green  int
	Orange int
}
